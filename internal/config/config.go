package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  Server  `yaml:"server" json:"server"`
	Bouncer Bouncer `yaml:"bouncer" json:"bouncer"`
	WAF     WAF     `yaml:"waf" json:"waf"`
}

type Server struct {
	GRPCPort int    `yaml:"grpcPort" json:"grpcPort"`
	LogLevel string `yaml:"logLevel" json:"logLevel"`
}

type Bouncer struct {
	Enabled         bool          `yaml:"enabled" json:"enabled"`
	Metrics         bool          `yaml:"metrics" json:"metrics"`
	TickerInterval  time.Duration `yaml:"tickerInterval" json:"tickerInterval"`
	MetricsInterval time.Duration `yaml:"metricsInterval" json:"metricsInterval"`
	ApiKey          string        `yaml:"apiKey" json:"apiKey"`
	LAPIURL         string        `yaml:"lapiUrl" json:"lapiUrl"`
	BanStatusCode   int           `yaml:"banStatusCode" json:"banStatusCode"`
}

type WAF struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	AppSecURL string `yaml:"appSecURL" json:"appSecURL"`
	ApiKey    string `yaml:"apiKey" json:"apiKey"`
}

func New(v *viper.Viper) (Config, error) {
	c := Config{}
	if v == nil {
		return c, errors.New("viper not initialized")
	}
	if v.ConfigFileUsed() != "" {
		err := v.ReadInConfig()
		if err != nil {
			return c, err
		}
	}
	err := v.Unmarshal(&c)
	return c, err
}
