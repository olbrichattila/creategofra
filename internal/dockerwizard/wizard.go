// Package dockerwizard creates a docker-compose.yml file wit the selected items
package dockerwizard

import (
	_ "embed"
	"strings"

	"github.com/olbrichattila/creategofra/internal/appwizard"
)

var storageMap = map[string]string{
	"file":      "",
	"db":        "",
	"redis":     redisFile,
	"memcached": memcachedFile,
}

//go:embed files/docker-redis
var redisFile string

//go:embed files/docker-mysql
var mySqlFile string

//go:embed files/docker-pgsql
var pgsqlFile string

//go:embed files/docker-memcached
var memcachedFile string

//go:embed files/docker-mailtrap
var mailtrapFile string

//go:embed files/docker-firebird
var firebirdFile string

type wizard struct {
	dbConnectionName string
	envData          []appwizard.EnvData
	storages         []string
	hasMailConfig    bool
}

func Wizard(dbConnectionName string, envData []appwizard.EnvData, storages []string, hasMailConfig bool) string {
	w := &wizard{
		dbConnectionName: dbConnectionName,
		envData:          envData,
		storages:         storages,
		hasMailConfig:    hasMailConfig,
	}

	return w.Run()

}

func (w *wizard) Run() string {
	composerContent := w.getComposeHead()

	composerContent += w.getComposeVolumes()

	return composerContent
}

func (w *wizard) getComposeHead() string {
	head := `version: '3.8'

services:
`
	switch w.dbConnectionName {
	case "mysql":
		head += w.fillTemplate(mySqlFile)
	case "pgsql":
		head += w.fillTemplate(pgsqlFile)
	case "firebird":
		head += w.fillTemplate(firebirdFile)
	}

	for _, selectedStorageName := range w.storages {
		if storageContent, ok := storageMap[selectedStorageName]; ok {
			head += w.fillTemplate(storageContent)
		}
	}

	if w.hasMailConfig {
		head += w.fillTemplate(mailtrapFile)
	}

	return head
}

func (w *wizard) getComposeVolumes() string {
	dbVolume := ""
	redisData := ""

	switch w.dbConnectionName {
	case "mysql":
		dbVolume = "  mysql_data:\n"
	case "pgsql":
		dbVolume = "  postgres_data:\n"
	case "firebird":
		dbVolume = "  firebird_data:\n"
	}

	for _, storageName := range w.storages {
		if storageName == "redis" {
			redisData = "  redis_data:\n"
		}
	}

	if dbVolume == "" && redisData == "" {
		return ""
	}

	return "volumes:\n" + dbVolume + redisData
}

func (w *wizard) fillTemplate(content string) string {
	for _, e := range w.envData {
		content = strings.ReplaceAll(content, e.Key, e.Value)
	}

	return content
}
