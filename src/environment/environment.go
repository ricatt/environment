package environment

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "env"

func LoadFile[T any](path string, target *T, config Config) (err error) {
	var content map[string]string
	content, err = parseEnvFile(path)
	if err != nil {
		return err
	}

	tp := reflect.ValueOf(target)
	proxy := reflect.ValueOf(*target)
	for i := 0; i < proxy.NumField(); i++ {
		field := proxy.Type().Field(i)
		tag := field.Tag.Get(tagName)
		switch field.Type.Kind() {
		case reflect.Bool:
			v, _ := strconv.ParseBool(content[tag])
			tp.Elem().FieldByName(field.Name).SetBool(v)
		case reflect.Int:
			v, _ := strconv.ParseInt(content[tag], 10, 0)
			tp.Elem().FieldByName(field.Name).SetInt(v)
		case reflect.Float64, reflect.Float32:
			v, _ := strconv.ParseFloat(content[tag], 0)
			tp.Elem().FieldByName(field.Name).SetFloat(v)
		case reflect.String:
			tp.Elem().FieldByName(field.Name).SetString(content[tag])
		default:
			return fmt.Errorf("env: type \"%s\" not supported", field.Type.Kind())
		}

		if content[tag] == "" && config.Force {
			return fmt.Errorf("missing value for %s", field.Name)
		}
	}
	return
}

func parseEnvFile(path string) (content map[string]string, err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		var val string
		if len(line) == 2 {
			val = line[1]
		}
		content[line[0]] = val
	}
	return nil, scanner.Err()
}
