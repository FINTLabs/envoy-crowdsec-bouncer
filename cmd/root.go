package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "envoy-bouncer",
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

	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		println("Error binding flags")
		os.Exit(-1)
	}
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
