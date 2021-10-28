package util

import (
	"fmt"
	"os"
	"reflect"
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

type blackJackInfomation struct {
	Name       string            `default:"BlackJack"`
	Age        int               `default:"18"`
	Money      uint64            `default:"420000"`
	Salary     float64           `default:"12000.56"`
	Alive      bool              `default:"true"`
	Family     [2]string         `default:"father,mother,child"`
	TVChannel  [2]int            `default:"66,99,55"`
	FavorColor [3]byte           `default:"100,101,253"`
	Friends    []string          `default:"bob,alice,mike"`                       // only simple types like string or number
	PhonesBook map[string]string `default:"bob=010-15235789,alice=0825-54567893"` // only simple types like string or number
	BlackSmith *neighbor
	WhiteSmith neighbor
}

type neighbor struct {
	XName string `default:"nnn"`
	XAge  int8   `default:"99"`
}

func TestDefaultValue(t *testing.T) {
	data := &blackJackInfomation{BlackSmith: &neighbor{}, Friends: []string{"aa"}} // if set some field manually, will not set default value to them
	if err := DefaultValue(data); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(data)
	}
}

func TestBasicTypeReflectSetValue(t *testing.T) {
	aa := 11
	if err := BasicTypeReflectSetValue(reflect.ValueOf(&aa).Elem(), "22"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(aa)
	}
}
