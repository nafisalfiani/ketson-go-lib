package configreader

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nafisalfiani/ketson-go-lib/codes"
	"github.com/nafisalfiani/ketson-go-lib/errors"
)

type dotEnvReader struct {
	opt Options
}

func initDotEnv(opt Options) Interface {
	return &dotEnvReader{
		opt: opt,
	}
}

func (d *dotEnvReader) ReadConfig(cfg any) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return loadNestedStruct(cfg, "")
}

func (d *dotEnvReader) AllSettings() map[string]any {
	return map[string]any{}
}

func loadNestedStruct(obj any, prefix string) error {
	configValue := reflect.ValueOf(obj).Elem()
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envKey := field.Tag.Get("env")

		if envKey == "" {
			continue
		}

		if prefix != "" {
			envKey = fmt.Sprintf("%s_%s", prefix, envKey)
		}

		envValue := os.Getenv(envKey)
		if envValue == "" {
			// If the environment variable is not set, skip it
			continue
		}

		fieldName := field.Name
		fieldValue := configValue.FieldByName(fieldName)

		switch field.Type.Kind() {
		case reflect.Int, reflect.Int64:
			intValue, err := strconv.Atoi(envValue)
			if err != nil {
				return errors.NewWithCode(codes.CodeInvalidValue, fmt.Sprintf("invalid value for %s: %s", envKey, err))
			}
			fieldValue.SetInt(int64(intValue))
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Struct:
			err := loadNestedStruct(fieldValue.Addr().Interface(), envKey)
			if err != nil {
				return err
			}
		default:
			return errors.NewWithCode(codes.CodeInvalidValue, "unsupported field type: %s", field.Type.Kind())
		}
	}

	return nil
}
