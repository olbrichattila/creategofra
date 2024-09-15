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

	"github.com/olbrichattila/creategofra/appwizard"
)

//go:embed files/blank.zip
var zipData []byte

var processChars = []string{"\\", "|", "/", "-"}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage creategofra <project-name>")
		return
	}

	projectName := os.Args[1]

	extract(projectName)
	initGoApp(projectName)
	appwizard.Wizard(projectName + "/.env")

	fmt.Print("\nDone\n")
}

func extract(projectName string) {
	projectName = projectName + "/"
	// Read the embedded zip file
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		fmt.Println("Failed to read zip file:", err)
		return
	}

	i := 0
	// Iterate through the files in the zip archive
	for _, file := range zipReader.File {
		i++
		process(i)
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

	// 	// Create the file in the local filesystem
	// 	outFile, err := os.Create(targetFileName)
	// 	if err != nil {
	// 		fmt.Println("\nFailed to create file:", err)
	// 		return
	// 	}
	// 	defer outFile.Close()

	// 	// Copy the file data from the zip to the local file
	// 	_, err = io.Copy(outFile, zipFile)
	// 	if err != nil {
	// 		fmt.Println("\nFailed to copy file:", err)
	// 		return
	// 	}
	// }
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

	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error: %v\nOutput: %s", err, output)
	}

	fmt.Println(string(output))
}

func process(i int) {
	pos := i % 4
	fmt.Print("Generating code: " + processChars[pos] + "\r")
}
