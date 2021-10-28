package util

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ParseStructWithEnv pasr struct tag(convert style see StrToEnvName function) and check if OS environment name matched tag
// if matched tag, set environment variable name as struct field value
// current support struct field value type are [Int Bool String Struct]
// usage example see struct_test.go#TestParseStructWithEnv()
// ParseStructWithEnv 使用结构体的tag映射环境变量值，tag和环境变量名的转换参照函数StrToEnvName,当前只支持 [Int Bool String Struct]
// 使用示范见struct_test.go#TestParseStructWithEnv()

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

// DefaultValue simply set default value to struct field by "default" tag
// current support struct field types [ float32 float64 uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 bool string slice array map ],current not support pointer types, usage example see struct_test.go#TestDefaultValue()
// DefaultValue 简单的设置默认值给结构体字段，通过default标签
// 当前支持的结构体字段类型 [ float32 float64 uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 bool string slice array map ],当前不支持指针类型变量
// 使用示范见struct_test.go#TestDefaultValue()
func DefaultValue(data interface{}) (err error) {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	if dataType.Kind() != reflect.Ptr {
		return fmt.Errorf("param data is not a pointer")
	}

	strcutType := dataType.Elem()
	structValue := dataValue.Elem()

	for i := 0; i < strcutType.NumField(); i++ {
		fieldType := strcutType.Field(i)
		fieldValue := structValue.Field(i)
		defaultValue := fieldType.Tag.Get("default")

		if defaultValue == "" && fieldType.Type.Kind() != reflect.Ptr && fieldType.Type.Kind() != reflect.Struct {
			// struct has no default tag ,thier fields has default tag
			continue
		} else if !fieldValue.IsZero() && fieldType.Type.Kind() != reflect.Ptr && fieldType.Type.Kind() != reflect.Struct {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Bool:
			err = BasicTypeReflectSetValue(fieldValue, defaultValue)

		case reflect.Struct:
			err = DefaultValue(structValue.Field(i).Addr().Interface())

		case reflect.Slice:
			sliceData := strings.Split(defaultValue, ",")
			slcieDataLength := len(sliceData)
			sliceValue := reflect.MakeSlice(fieldType.Type, len(sliceData), 0)
			for k, v := range sliceData {
				if k < slcieDataLength {
					BasicTypeReflectSetValue(fieldValue.Index(i), v)
				}
			}
			fieldValue.Set(sliceValue)

		case reflect.Array:
			arrayLen := fieldValue.Len()
			for k, v := range strings.Split(defaultValue, ",") {
				if k < arrayLen {
					BasicTypeReflectSetValue(fieldValue.Index(k), v)
				}
			}

		case reflect.Map:
			mapValue := reflect.MakeMap(fieldType.Type)

			for _, v := range strings.Split(defaultValue, ",") {
				mapElemSlice := strings.Split(v, "=")
				sliceLength := len(mapElemSlice)
				if sliceLength > 1 {
					mapValue.SetMapIndex(reflect.ValueOf(mapElemSlice[0]), reflect.ValueOf(strings.Join(mapElemSlice[1:], "=")))
				} else if sliceLength == 1 {
					mapValue.SetMapIndex(reflect.ValueOf(mapElemSlice[0]), reflect.Value{})
				}
			}

			fieldValue.Set(mapValue)

		case reflect.Ptr:
			if !fieldValue.IsNil() {
				err = DefaultValue(structValue.Field(i).Interface())
			}

		default:
			err = fmt.Errorf("unsupport struct field type %s", fieldType.Type.String())
		}

		if err != nil {
			return err
		}
	}

	return err
}

// BasicTypeReflectSetValue use reflect.Value set value, main used for function DefaultValue(), usage see struct_test.go#TestBasicTypeReflectSetValue()，current not support pointer types,
// BasicTypeReflectSetValue 使用reflect.Value来修改变量的值，主要用于函数DefaultValue()，使用示范见struct_test.go#TestDefaultValue(),当前不支持指针类型变量
func BasicTypeReflectSetValue(rv reflect.Value, value string) (err error) {
	switch rv.Kind() {
	case reflect.String:
		rv.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("convert tag default value to int64 failed -> %s", err.Error())
		}
		rv.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("convert tag default value to uint64 failed -> %s", err.Error())
		}
		rv.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("convert tag default value to float64 failed -> %s", err.Error())
		}
		rv.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("convert tag default value to bool failed -> %s", err.Error())
		}
		rv.SetBool(boolValue)

	default:
		err = fmt.Errorf("unsupport struct field type %s", rv.Type().String())
	}
	return err
}
