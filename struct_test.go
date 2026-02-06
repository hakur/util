package util

import (
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
	// 测试 String 类型
	var str string
	err := BasicTypeReflectSetValue(reflect.ValueOf(&str).Elem(), "hello")
	assert.Nil(t, err)
	assert.Equal(t, "hello", str)

	// 测试 Int 类型
	var numInt int
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numInt).Elem(), "42")
	assert.Nil(t, err)
	assert.Equal(t, 42, numInt)

	// 测试 Int8 类型
	var numInt8 int8
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numInt8).Elem(), "127")
	assert.Nil(t, err)
	assert.Equal(t, int8(127), numInt8)

	// 测试 Int16 类型
	var numInt16 int16
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numInt16).Elem(), "32767")
	assert.Nil(t, err)
	assert.Equal(t, int16(32767), numInt16)

	// 测试 Int32 类型
	var numInt32 int32
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numInt32).Elem(), "2147483647")
	assert.Nil(t, err)
	assert.Equal(t, int32(2147483647), numInt32)

	// 测试 Int64 类型
	var numInt64 int64
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numInt64).Elem(), "9223372036854775807")
	assert.Nil(t, err)
	assert.Equal(t, int64(9223372036854775807), numInt64)

	// 测试 Uint 类型
	var numUint uint
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numUint).Elem(), "42")
	assert.Nil(t, err)
	assert.Equal(t, uint(42), numUint)

	// 测试 Uint8 类型
	var numUint8 uint8
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numUint8).Elem(), "255")
	assert.Nil(t, err)
	assert.Equal(t, uint8(255), numUint8)

	// 测试 Uint16 类型
	var numUint16 uint16
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numUint16).Elem(), "65535")
	assert.Nil(t, err)
	assert.Equal(t, uint16(65535), numUint16)

	// 测试 Uint32 类型
	var numUint32 uint32
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numUint32).Elem(), "4294967295")
	assert.Nil(t, err)
	assert.Equal(t, uint32(4294967295), numUint32)

	// 测试 Uint64 类型
	var numUint64 uint64
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numUint64).Elem(), "18446744073709551615")
	assert.Nil(t, err)
	assert.Equal(t, uint64(18446744073709551615), numUint64)

	// 测试 Float32 类型
	var numFloat32 float32
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numFloat32).Elem(), "3.14159")
	assert.Nil(t, err)
	assert.InDelta(t, float32(3.14159), numFloat32, 0.0001)

	// 测试 Float64 类型
	var numFloat64 float64
	err = BasicTypeReflectSetValue(reflect.ValueOf(&numFloat64).Elem(), "3.141592653589793")
	assert.Nil(t, err)
	assert.InDelta(t, float64(3.141592653589793), numFloat64, 0.0001)

	// 测试 Bool 类型
	var flag bool
	err = BasicTypeReflectSetValue(reflect.ValueOf(&flag).Elem(), "true")
	assert.Nil(t, err)
	assert.Equal(t, true, flag)

	err = BasicTypeReflectSetValue(reflect.ValueOf(&flag).Elem(), "false")
	assert.Nil(t, err)
	assert.Equal(t, false, flag)

	err = BasicTypeReflectSetValue(reflect.ValueOf(&flag).Elem(), "1")
	assert.Nil(t, err)
	assert.Equal(t, true, flag)

	err = BasicTypeReflectSetValue(reflect.ValueOf(&flag).Elem(), "0")
	assert.Nil(t, err)
	assert.Equal(t, false, flag)

	// 测试无效的 Int 值
	var invalidInt int
	err = BasicTypeReflectSetValue(reflect.ValueOf(&invalidInt).Elem(), "invalid")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "convert tag default value to int64 failed")

	// 测试无效的 Bool 值
	var invalidBool bool
	err = BasicTypeReflectSetValue(reflect.ValueOf(&invalidBool).Elem(), "yes")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "convert tag default value to bool failed")

	// 测试 Int 溢出
	var overflowInt8 int8
	err = BasicTypeReflectSetValue(reflect.ValueOf(&overflowInt8).Elem(), "999")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "int overflow")

	// 测试 Uint 溢出
	var overflowUint8 uint8
	err = BasicTypeReflectSetValue(reflect.ValueOf(&overflowUint8).Elem(), "999")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "uint overflow")

	// 测试 Float 溢出
	var overflowFloat32 float32
	err = BasicTypeReflectSetValue(reflect.ValueOf(&overflowFloat32).Elem(), "1e100")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "float overflow")

	// 测试不支持的类型
	slice := []int{1, 2, 3}
	err = BasicTypeReflectSetValue(reflect.ValueOf(&slice).Elem(), "test")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported struct field type")
}
