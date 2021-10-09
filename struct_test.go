package util

import (
	"fmt"
	"os"
	"testing"
)

type TestAppConfig struct {
	TimeZone string `yaml:"timeZone"`
	LogLevel string `yaml:"logLevel"`
	Web      struct {
		Port     int    `yaml:"port"`
		Address  string `yaml:"address"`
		SiteName string `yaml:"siteName"`
	} `yaml:"web"`
	NatsStreaming struct {
		Url      string `yaml:"url"`
		Token    string `yaml:"token"`
		ClientID string `yaml:"clientId"`
	} `yaml:"nats"`
	DB struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Name        string `yaml:"name"`
		MinConn     int    `yaml:"minConn"`
		MaxConn     int    `yaml:"maxConn"`
		AutoMigrate bool   `yaml:"autoMigrate"`
		TablePrefix string `yaml:"tablePrefix"`
	} `yaml:"db"`
}

func TestParseStructWithEnv(t *testing.T) {
	os.Setenv("NATS_STREAMING_URL", "127.0.0.1")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_AUTO_MIGRATE", "true")
	os.Setenv(StrToEnvName("db.tablePrefix"), "myTablePrefix")
	appConfig := new(TestAppConfig)
	ParseStructWithEnv(appConfig, "")
	fmt.Println(appConfig.NatsStreaming.Url, appConfig.DB.Port, appConfig.DB.AutoMigrate, appConfig.DB.TablePrefix)
}
