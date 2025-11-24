package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
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
}
