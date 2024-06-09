package reflects

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
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
}
