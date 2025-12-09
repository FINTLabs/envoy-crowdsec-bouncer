package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("nil viper returns error", func(t *testing.T) {
		c, err := New(nil)
		assert.Error(t, err)
		assert.Equal(t, "viper not initialized", err.Error())
		assert.Empty(t, c)
	})

	t.Run("valid viper config", func(t *testing.T) {
		v := viper.New()
		v.Set("server.grpcPort", 8080)
		v.Set("server.logLevel", "debug")
		v.Set("bouncer.apiKey", "test-key")
		v.Set("bouncer.lapiURL", "http://test.com")
		v.Set("bouncer.metrics", true)
		v.Set("bouncer.tickerInterval", "30s")
		v.Set("bouncer.metricsInterval", "100s")
		v.Set("waf.enabled", true)
		v.Set("waf.apiKey", "test-key")
		v.Set("waf.appSecURL", "http://test.com")

		c, err := New(v)
		assert.NoError(t, err)
		assert.Equal(t, 8080, c.Server.GRPCPort)
		assert.Equal(t, "debug", c.Server.LogLevel)

		assert.Equal(t, "test-key", c.Bouncer.ApiKey)
		assert.Equal(t, "http://test.com", c.Bouncer.LAPIURL)
		assert.True(t, c.Bouncer.Metrics)
		assert.Equal(t, time.Second*30, c.Bouncer.TickerInterval)
		assert.Equal(t, time.Second*100, c.Bouncer.MetricsInterval)
		assert.True(t, c.WAF.Enabled)
		assert.Equal(t, "test-key", c.WAF.ApiKey)
		assert.Equal(t, "http://test.com", c.WAF.AppSecURL)
	})

	t.Run("non existing config file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test-config.yaml")

		v := viper.New()
		v.SetConfigFile(path)

		c, err := New(v)
		assert.Error(t, err)
		assert.Empty(t, c)
	})

	t.Run("valid viper config from file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test-config.yaml")
		configContent := `
server:
  grpcPort: 8080
  logLevel: error
bouncer:
  apiKey: test-key
waf:
  enabled: true
`
		err := os.WriteFile(path, []byte(configContent), 0644)
		require.NoError(t, err)

		v := viper.New()
		v.SetConfigFile(path)

		cfg, err := New(v)
		require.NoError(t, err)

		assert.Equal(t, path, v.ConfigFileUsed())
		assert.Equal(t, 8080, cfg.Server.GRPCPort)
		assert.Equal(t, "error", cfg.Server.LogLevel)
		assert.Equal(t, "test-key", cfg.Bouncer.ApiKey)
		assert.True(t, cfg.WAF.Enabled)
	})
}
