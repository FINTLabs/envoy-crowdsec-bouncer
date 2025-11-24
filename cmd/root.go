package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/mattis/envoy-crowdsec-bouncer/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "envoy-bouncer",

	RunE: func(cmd *cobra.Command, args []string) error {
		v := viper.GetViper()
		config, err := config.New(v)
		if err != nil {
			return err
		}

		data, err := json.Marshal(config)
		if err != nil {
			return err
		}
		println(string(data))

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (json or yaml)")

	rootCmd.Flags().Int("server.grpcPort", 8080, "")
	rootCmd.Flags().String("server.logLevel", slog.LevelInfo.String(), "")

	rootCmd.Flags().Bool("bouncer.enabled", true, "")
	rootCmd.Flags().Bool("bouncer.metrics", false, "")
	rootCmd.Flags().Duration("bouncer.tickerInterval", time.Second*10, "")
	rootCmd.Flags().Duration("bouncer.metricsInterval", time.Minute*10, "")
	rootCmd.Flags().Int("bouncer.banStatusCode", 403, "")

	rootCmd.Flags().Bool("waf.enabled", false, "")

	viper.BindPFlags(rootCmd.Flags())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvPrefix("ENVOY_BOUNCER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", ""))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
