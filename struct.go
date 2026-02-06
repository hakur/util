package util

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ParseStructWithEnv 解析结构体字段并从环境变量中读取值填充
// 通过 StrToEnvName 函数将结构体字段名转换为环境变量名
// 支持的类型: [Int8 Int16 Int32 Int64 Uint8 Uint16 Uint32 Uint64 Uint Int Float32 Float64 Bool String Struct]
// 如果环境变量不存在或解析失败，跳过该字段
// 返回错误以提示解析过程中的问题
func ParseStructWithEnv(structNode interface{}, rootNodeName string) error {
	// 检查 nil
	if structNode == nil {
		return fmt.Errorf("structNode is nil")
	}

	// 获取反射值
	tp := reflect.TypeOf(structNode)
	val := reflect.ValueOf(structNode)

	// 处理指针类型
	if tp.Kind() == reflect.Ptr {
		if val.IsNil() {
			return fmt.Errorf("structNode is nil pointer")
		}
		val = reflect.Indirect(val)
		tp = val.Type()
	}

	// 验证必须是结构体
	if tp.Kind() != reflect.Struct {
		return fmt.Errorf("structNode must be struct or pointer to struct, got %s", tp.Kind())
	}

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := tp.Field(i)

		// 检查字段是否可设置且可导出
		if !field.CanSet() {
			continue
		}
		if !fieldType.IsExported() {
			continue
		}

		// 递归处理嵌套结构体
		if field.Kind() == reflect.Struct {
			if err := ParseStructWithEnv(field.Addr().Interface(), rootNodeName+"_"+fieldType.Name); err != nil {
				return err
			}
			continue
		}

		// 处理嵌套指针结构体 *Struct
		if field.Kind() == reflect.Ptr && !field.IsNil() && field.Type().Elem().Kind() == reflect.Struct {
			if err := ParseStructWithEnv(field.Interface(), rootNodeName+"_"+fieldType.Name); err != nil {
				return err
			}
			continue
		}

		// 构造环境变量名
		envName := StrToEnvName(rootNodeName + "_" + fieldType.Name)
		env := os.Getenv(envName)
		if env == "" {
			continue
		}

		// 根据类型解析环境变量值
		switch field.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(env)
			if err != nil {
				return fmt.Errorf("parse bool field %s from env %s failed: %w", fieldType.Name, envName, err)
			}
			field.SetBool(v)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(env, 10, 64)
			if err != nil {
				return fmt.Errorf("parse int field %s from env %s failed: %w", fieldType.Name, envName, err)
			}
			if field.OverflowInt(v) {
				return fmt.Errorf("int field %s overflow with value %s", fieldType.Name, env)
			}
			field.SetInt(v)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(env, 10, 64)
			if err != nil {
				return fmt.Errorf("parse uint field %s from env %s failed: %w", fieldType.Name, envName, err)
			}
			if field.OverflowUint(v) {
				return fmt.Errorf("uint field %s overflow with value %s", fieldType.Name, env)
			}
			field.SetUint(v)

		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(env, 64)
			if err != nil {
				return fmt.Errorf("parse float field %s from env %s failed: %w", fieldType.Name, envName, err)
			}
			if field.OverflowFloat(v) {
				return fmt.Errorf("float field %s overflow with value %s", fieldType.Name, env)
			}
			field.SetFloat(v)

		case reflect.String:
			field.SetString(env)

		default:
			return fmt.Errorf("unsupported field type %s for field %s", field.Kind(), fieldType.Name)
		}
	}

	return nil
}

// DefaultValue simply set default value to struct field by "default" tag
// current support struct field types [ float32 float64 uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 bool string slice array map ],current not support pointer types, not support unexported field(lower case named field), usage example see struct_test.go#TestDefaultValue()
// DefaultValue 简单的设置默认值给结构体字段，通过default标签
// 当前支持的结构体字段类型 [ float32 float64 uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 bool string slice array map ],当前不支持指针类型变量,不支持未导出的结构体字段（名字小写的字段）
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
			if structValue.Field(i).CanSet() {
				if err = DefaultValue(structValue.Field(i).Addr().Interface()); err != nil {
					return err
				}
			}

		case reflect.Slice:
			sliceData := strings.Split(defaultValue, ",")
			sliceLength := len(sliceData)
			sliceValue := reflect.MakeSlice(fieldType.Type, sliceLength, sliceLength)
			for k, v := range sliceData {
				if k < sliceLength {
					if err = BasicTypeReflectSetValue(sliceValue.Index(k), v); err != nil {
						return fmt.Errorf("slice field %s index %d: %w", fieldType.Name, k, err)
					}
				}
			}
			fieldValue.Set(sliceValue)

		case reflect.Array:
			arrayLen := fieldValue.Len()
			arrayData := strings.Split(defaultValue, ",")
			for k, v := range arrayData {
				if k < arrayLen {
					if err = BasicTypeReflectSetValue(fieldValue.Index(k), v); err != nil {
						return fmt.Errorf("array field %s index %d: %w", fieldType.Name, k, err)
					}
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
			if !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Struct {
				if err = DefaultValue(fieldValue.Interface()); err != nil {
					return err
				}
			}

		default:
			err = fmt.Errorf("unsupported struct field type %s for field %s", fieldType.Type.String(), fieldType.Name)
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
			return fmt.Errorf("convert tag default value to int64 failed: %w", err)
		}
		if rv.OverflowInt(intValue) {
			return fmt.Errorf("int overflow for value %s", value)
		}
		rv.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("convert tag default value to uint64 failed: %w", err)
		}
		if rv.OverflowUint(uintValue) {
			return fmt.Errorf("uint overflow for value %s", value)
		}
		rv.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("convert tag default value to float64 failed: %w", err)
		}
		if rv.OverflowFloat(floatValue) {
			return fmt.Errorf("float overflow for value %s", value)
		}
		rv.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("convert tag default value to bool failed: %w", err)
		}
		rv.SetBool(boolValue)

	default:
		err = fmt.Errorf("unsupported struct field type %s", rv.Type().String())
	}
	return err
}
