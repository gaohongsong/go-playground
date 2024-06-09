package reflects

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Name  string `json:"server-name"`
	Ip    string `json:"server-ip"`
	Port  int    `json:"server-port"`
	Debug bool   `json:"debug"`
}

func readConfigV1() *Config {
	// TODO: read from settings.json
	config := Config{}

	typ := reflect.TypeOf(config)
	value := reflect.Indirect(reflect.ValueOf(&config))
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if v, ok := f.Tag.Lookup("json"); ok {
			key := fmt.Sprintf("CONFIG_%s", strings.ToUpper(strings.ReplaceAll(v, "-", "_")))
			if env, exist := os.LookupEnv(key); exist {
				value.FieldByName(f.Name).Set(reflect.ValueOf(env))
			}
		}
	}

	return &config
}

func readConfig() *Config {
	config := Config{}
	if data, err := os.ReadFile("settings.json"); err == nil {
		_ = json.Unmarshal(data, &config)
	}

	typ := reflect.TypeOf(config)
	value := reflect.Indirect(reflect.ValueOf(&config))
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if v, ok := f.Tag.Lookup("json"); ok {
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(v, "-", "_"))
			if env, exist := os.LookupEnv(key); exist {
				field := value.FieldByName(f.Name)
				if field.IsValid() && field.CanSet() {
					//field.Set(reflect.ValueOf(env))
					switch field.Kind() {
					case reflect.String:
						field.SetString(env)
					case reflect.Int:
						if intValue, err := strconv.Atoi(env); err == nil {
							field.SetInt(int64(intValue))
						}
					case reflect.Bool:
						if boolValue, err := strconv.ParseBool(env); err == nil {
							field.SetBool(boolValue)
						}
					default:
						fmt.Printf("Unsupported field type: %s\n", field.Kind())
					}
				}
			}
		}
	}

	return &config
}
