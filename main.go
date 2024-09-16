package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/olbrichattila/creategofra/internal/appwizard"
	"github.com/olbrichattila/creategofra/internal/specio"
)

//go:embed files/blank.zip
var blankAppZipData []byte

//go:embed files/regapp.zip
var regAppZipData []byte

// Migration data

//go:embed files/firebird.zip
var firebirdZipData []byte

//go:embed files/sqlite.zip
var sqliteZipData []byte

//go:embed files/mysql.zip
var mysqlZipData []byte

//go:embed files/pgsql.zip
var pgsqlZipData []byte

var processChars = []string{"\\", "|", "/", "-"}

var skipMigration = []string{
	"2024-07-31_21_01_53-migrate--user.sql",
	"2024-07-31_21_01_53-rollback--user.sql",
	"2024-08-15_20_58_24-migrate--reg_confirmations.sql",
	"2024-08-15_20_58_24-rollback--reg_confirmations.sql",
	"2024-08-26_15_35_47-migrate--password-reminder.sql",
	"2024-08-26_15_35_47-rollback--password-reminder.sql",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage creategofra <project-name>")
		return
	}

	projectName := os.Args[1]
	if validated := validate(projectName); validated != "" {
		fmt.Println(validated)
		return
	}

	selection := extractRequestedVersion(projectName)

	initGoApp(projectName)
	responses := appwizard.Wizard(projectName + "/.env")

	copyMigrations(projectName, selection, responses)

	fmt.Print("\nDone\n")
}

func extractRequestedVersion(projectName string) string {
	fmt.Println(`Select project type
  1. Blank project
  2. Project with login/registration`)

	for {
		selection := specio.Input("Please select (1 or 2): ", "")
		if selection == "1" {
			fmt.Println()
			extract(projectName, "project source code", "", []string{}, &blankAppZipData)
			return selection
		}

		if selection == "2" {
			fmt.Println()
			extract(projectName, "project source code", "", []string{}, &regAppZipData)
			return selection
		}

		fmt.Println("\nInvalid selection")
	}
}

func validate(projectName string) string {
	_, err := os.Stat(projectName)
	if os.IsNotExist(err) {
		return ""
	}

	return fmt.Sprintf("Project '%s' already exists!", projectName)
}

func extract(projectName, taskName, subFolder string, skipData []string, data *[]byte) {
	projectName = projectName + "/"
	if subFolder != "" {
		projectName = projectName + "/" + subFolder + "/"
	}

	// Read the embedded zip file
	zipReader, err := zip.NewReader(bytes.NewReader(*data), int64(len(*data)))
	if err != nil {
		fmt.Println("Failed to read zip file:", err)
		return
	}

	i := 0
	// Iterate through the files in the zip archive
	for _, file := range zipReader.File {
		i++
		process(i, taskName)
		if contains(file.Name, skipData) {
			continue
		}

		targetFileName := projectName + file.Name

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(targetFileName, os.ModePerm)
			if err != nil {
				fmt.Println("\nFailed to create directory:", err)
				return
			}
			continue
		}

		// Open the file inside the zip
		zipFile, err := file.Open()
		if err != nil {
			fmt.Println("\nFailed to open file in zip:", err)
			return
		}
		defer zipFile.Close()

		// Create directories for the file if necessary
		if err := os.MkdirAll(filepath.Dir(targetFileName), os.ModePerm); err != nil {
			fmt.Println("\nFailed to create directory for file:", err)
			return
		}

		// If the file is a Go file, read, replace and write its content
		if filepath.Ext(file.Name) == ".go" {
			content, err := io.ReadAll(zipFile)
			if err != nil {
				fmt.Println("\nFailed to read file:", err)
				return
			}

			// Replace "gofraapp/" with projectName
			modifiedContent := strings.ReplaceAll(string(content), "\"gofraapp/", "\""+projectName)

			// Create the file in the local filesystem
			outFile, err := os.Create(targetFileName)
			if err != nil {
				fmt.Println("\nFailed to create file:", err)
				return
			}
			defer outFile.Close()

			// Write the modified content back to the file
			_, err = outFile.Write([]byte(modifiedContent))
			if err != nil {
				fmt.Println("\nFailed to write modified content to file:", err)
				return
			}
		} else {
			// For non-Go files, just copy the file data as is
			outFile, err := os.Create(targetFileName)
			if err != nil {
				fmt.Println("\nFailed to create file:", err)
				return
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, zipFile)
			if err != nil {
				fmt.Println("\nFailed to copy file:", err)
				return
			}
		}
	}
}

func initGoApp(projectName string) {
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = projectName

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	}

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = projectName

	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	}
}

func process(i int, taskName string) {
	pos := i % 4
	fmt.Printf("Generating %s: %s\r", taskName, processChars[pos])
	// extraction is fast, so the user can see it does something
	time.Sleep(30 * time.Millisecond)
}

func copyMigrations(projectName, selection string, responses []appwizard.EnvData) {
	dbConnectionName := getDbConnection(responses)
	skipData := []string{}
	if selection == "1" {
		skipData = skipMigration
	}

	switch dbConnectionName {
	case "sqlite":
		extract(projectName, "sqlite migrations", "migrations", skipData, &sqliteZipData)
	case "mysql":
		extract(projectName, "MySql migrations", "migrations", skipData, &mysqlZipData)
	case "pgsql":
		extract(projectName, "PostgresQl migrations", "migrations", skipData, &pgsqlZipData)
	case "firebird":
		extract(projectName, "Firebird migrations", "migrations", skipData, &firebirdZipData)
	default:
		fmt.Print("Skip generating, migrations not set")
	}
}

func getDbConnection(responses []appwizard.EnvData) string {
	for _, e := range responses {
		if e.Key == "DB_CONNECTION" {
			return e.Value
		}
	}
	return ""
}

func contains(item string, data []string) bool {
	for _, v := range data {
		if v == item {
			return true
		}
	}
	return false
}
