package util

import (
	"os"
	"reflect"
	"strconv"
)

// ParseStructWithEnv pasr struct tag(convert style see StrToEnvName function) and check if OS environment name matched tag
// if matched tag, set environment variable name as struct field value
// current support struct field value type are [Int Bool String Struct]
// ParseStructWithEnv 使用结构体的tag映射环境变量值，tag和环境变量名的转换参照函数StrToEnvName,当前只支持 [Int Bool String Struct]
// usage example see struct_test.go
func ParseStructWithEnv(structNode interface{}, rootNodeName string) {
	tp := reflect.TypeOf(structNode)
	var val reflect.Value

	if tp.Kind() == reflect.Ptr {
		ov := reflect.ValueOf(structNode)
		val = reflect.Indirect(ov)
	} else {
		val = reflect.ValueOf(structNode)
	}

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Kind() == reflect.Struct {
			ParseStructWithEnv(val.Field(i).Addr().Interface(), val.Type().Field(i).Name)
		} else {
			envName := StrToEnvName(rootNodeName + "_" + val.Type().Field(i).Name)
			env := os.Getenv(envName)
			if env == "" {
				continue
			}
			switch val.Field(i).Type().Kind() {
			case reflect.Bool:
				v, _ := strconv.ParseBool(env)
				val.Field(i).SetBool(v)
			case reflect.Int:
				v, _ := strconv.ParseInt(env, 10, 64)
				val.Field(i).SetInt(v)
			case reflect.String:
				val.Field(i).SetString(env)
			}
		}
	}
}
