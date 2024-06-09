package reflects

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Name  string `json:"server-name"`
	Ip    string `json:"server-ip"`
	Port  int    `json:"server-port"`
	Debug bool   `json:"debug"`
}

// 简陋版
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

// 单例模式保证config全局唯一，仅加载一次
var (
	config     *Config
	configOnce sync.Once
)

func GetConfig() *Config {
	configOnce.Do(loadConfig)
	return config
}

// 可以将readConfig直接改为loadConfig
func loadConfig() {
	_ = readConfig()
}

func readConfig() *Config {
	fmt.Println("READ CONFIG!!!")
	//config := Config{}
	config = &Config{}
	// 读配置文件
	if data, err := os.ReadFile("settings.json"); err == nil {
		_ = json.Unmarshal(data, config)
	}

	//反射类型
	typ := reflect.TypeOf(*config)
	// 反射值，兼容指针
	value := reflect.Indirect(reflect.ValueOf(config))
	// 遍历所有字段
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		// json tag 解析
		if v, ok := f.Tag.Lookup("json"); ok {
			// 环境变量名映射
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(v, "-", "_"))
			// 环境变量查找
			if env, exist := os.LookupEnv(key); exist {
				field := value.FieldByName(f.Name)
				// 可写判定
				if field.IsValid() && field.CanSet() {
					// 类型判定
					switch field.Kind() {
					case reflect.String:
						field.SetString(env)
					case reflect.Int:
						// 类型转换
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

	return config
}
