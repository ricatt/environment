package environment

import (
	"fmt"
	"github.com/joho/godotenv"
	"reflect"
	"strconv"
)

const tagName = "env"

var content map[string]string

func LoadFile[T any](path string, target *T, config Config) (err error) {
	content, err = godotenv.Read(path)
	if err != nil {
		return err
	}

	tp := reflect.ValueOf(target)
	proxy := reflect.ValueOf(*target)
	for i := 0; i < proxy.NumField(); i++ {
		field := proxy.Type().Field(i)
		tag := field.Tag.Get(tagName)

		if v, err := strconv.ParseBool(content[tag]); err == nil {
			tp.Elem().FieldByName(field.Name).SetBool(v)
		} else if v, err := strconv.ParseInt(content[tag], 10, 0); err == nil {
			tp.Elem().FieldByName(field.Name).SetInt(v)
		} else if v, err := strconv.ParseFloat(content[tag], 0); err == nil {
			tp.Elem().FieldByName(field.Name).SetFloat(v)
		} else if content[tag] != "" {
			tp.Elem().FieldByName(field.Name).SetString(content[tag])
		}

		if content[tag] == "" && config.Force {
			return fmt.Errorf("missing value for %s", field.Name)
		}
	}
	return
}
