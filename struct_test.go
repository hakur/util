package util

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, "127.0.0.1", appConfig.NatsStreaming.Url)
	assert.Equal(t, 3306, appConfig.DB.Port)
	assert.Equal(t, true, appConfig.DB.AutoMigrate)
	assert.Equal(t, "myTablePrefix", appConfig.DB.TablePrefix)
}

type blackJackInfomation struct {
	Name        string            `default:"BlackJack"`
	Age         int               `default:"18"`
	Money       uint64            `default:"420000"`
	Salary      float64           `default:"12000.56"`
	Alive       bool              `default:"true"`
	Family      [2]string         `default:"father,mother,child"`
	TVChannel   [2]int            `default:"66,99,55"`
	FavorColor  [3]byte           `default:"100,101,253"`
	Friends     []string          `default:"bob,alice,mike"`                       // only simple types like string or number
	PhonesBook  map[string]string `default:"bob=010-15235789,alice=0825-54567893"` // only simple types like string or number
	BlackSmith  *neighbor
	WhiteSmith  neighbor
	unexported1 neighbor
	unexported2 *neighbor
	unexported  int
}

type neighbor struct {
	XName      string `default:"nnn"`
	XAge       int8   `default:"99"`
	unexported int
}

func TestDefaultValue(t *testing.T) {
	data := &blackJackInfomation{BlackSmith: &neighbor{}, Friends: []string{"aa"}} // if set some field manually, will not set default value to them
	assert.Equal(t, nil, DefaultValue(data))
	assert.Equal(t, []string{"aa"}, data.Friends)
	assert.Equal(t, [2]int{66, 99}, data.TVChannel)
	assert.Equal(t, map[string]string{"bob": "010-15235789", "alice": "0825-54567893"}, data.PhonesBook)
	assert.Equal(t, int8(99), data.BlackSmith.XAge)
	assert.Equal(t, int8(99), data.WhiteSmith.XAge)
}

func TestBasicTypeReflectSetValue(t *testing.T) {
	aa := 11
	if err := BasicTypeReflectSetValue(reflect.ValueOf(&aa).Elem(), "22"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(aa)
	}
}

func TestParseDockerImageNameInfo(t *testing.T) {
	var imageNames = []string{
		"https://docker.io/pizza/rumia/rds-operator:v0.0.3",
		"https://docker.io/rumia/rds-operator:v0.0.3",
		"https://docker.io/rumia/rds-operator@sha256:123456789",
		"https://docker.io/rumia/rds-operator",
		"http://quay.io/rumia/rds-operator",
		"rumia/rds-operator",
		"centos",
	}

	for _, imageName := range imageNames {
		info := ParseDockerImageNameInfo(imageName)
		LogJSON(info)
	}
}
