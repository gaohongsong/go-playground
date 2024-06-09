package reflects

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	//os.Setenv("CONFIG_SERVER_NAME", "core")
	os.Setenv("CONFIG_SERVER_IP", "127.0.0.1")
	os.Setenv("CONFIG_SERVER_PORT", "8080")
	os.Setenv("CONFIG_DEBUG", "true")

	cfg := readConfig()
	cfgJson, _ := json.MarshalIndent(cfg, "", "\t")
	t.Logf("config: \n%s\n", cfgJson)

	assert.Equal(t, "go-core", cfg.Name)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, true, cfg.Debug)

	cfg1 := GetConfig()
	t.Logf("cfg1: \n%+v\n", cfg1)
	cfg2 := GetConfig()
	t.Logf("cfg2: \n%+v\n", cfg2)
	assert.Equal(t, cfg1, cfg2)
}

// go test -bench .
// BenchmarkNew-16                 60195547                20.58 ns/op
// BenchmarkReflectNew-16          42175406                26.96 ns/op
func BenchmarkNew(b *testing.B) {
	var config *Config
	for i := 0; i < b.N; i++ {
		config = new(Config)
	}
	_ = config
}

func BenchmarkReflectNew(b *testing.B) {
	var config *Config
	typ := reflect.TypeOf(Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config, _ = reflect.New(typ).Interface().(*Config)
	}
	_ = config
}

// go test -bench="Set$" .
// BenchmarkSet-16                         1000000000          0.1167 ns/op
// BenchmarkReflect_FieldSet-16            113654162          10.46 ns/op
// BenchmarkReflect_FieldByNameSet-16       7451679          170.4 ns/op
// 对于一个普通的拥有 4 个字段的结构体 Config 来说，使用反射给每个字段赋值，相比直接赋值，性能劣化约 100 - 1000 倍。
// 其中，FieldByName 的性能相比 Field 劣化 10 倍
// 在反射的内部，字段是按顺序存储的，因此按照下标访问查询效率为 O(1)，而按照 Name 访问，则需要遍历所有字段，查询效率为 O(N)
// 结构体所包含的字段(包括方法)越多，那么两者之间的效率差距则越大
func BenchmarkSet(b *testing.B) {
	config := new(Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Name = "name"
		config.Ip = "127.0.0.1"
		config.Port = 123
		config.Debug = true
	}
}

func BenchmarkReflect_FieldSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(0).SetString("name")
		ins.Field(1).SetString("ip")
		ins.Field(2).SetInt(123)
		ins.Field(3).SetBool(true)
	}
}

func BenchmarkReflect_FieldByNameSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.FieldByName("Name").SetString("name")
		ins.FieldByName("Ip").SetString("ip")
		ins.FieldByName("Port").SetInt(123)
		ins.FieldByName("Debug").SetBool(true)
	}
}

// BenchmarkReflect_FieldByNameSet-16               7076823               170.8 ns/op
// BenchmarkReflect_FieldByNameCacheSet-16         29598568                38.55 ns/op
func BenchmarkReflect_FieldByNameCacheSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	//利用字典将 Name 和 Index 的映射缓存起来。避免每次反复查找，耗费大量的时间
	cache := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		cache[typ.Field(i).Name] = i
	}
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(cache["Name"]).SetString("name")
		ins.Field(cache["Ip"]).SetString("ip")
		ins.Field(cache["Port"]).SetInt(123)
		ins.Field(cache["Debug"]).SetBool(true)
	}
}
