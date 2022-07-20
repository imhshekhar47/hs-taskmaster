/*
Copyright © 2022 Himanshu Shekhar <himanshu.kiit@gmail.com>
Code ownership is with Himanshu Shekhar. Use without modifications.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/imhshekhar47/go-api-core/config"
	"github.com/imhshekhar47/go-api-core/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmdLogger = logger.GetLogger("cmd/root")
	// args
	cfgFile     string
	argLogLevel string
	argPortGrpc uint16 = 50051
	argPortRest uint16 = 0
)

var rootCmd = &cobra.Command{
	Use:   "skill-api",
	Short: "Skill Management API",
	Long:  `Skill Management APIs`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&argLogLevel, "loglevel", "info", "log level trace, debug, info, error, warn, error")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.skill-api.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".skill-api")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	appConfig := config.ApplicationConfig{}
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err := viper.Unmarshal(&appConfig); err != nil {
			rootCmdLogger.Fatalf("Failed to unmarshal config: %v", err)
		}
	}

	//rootCmdLogger.Tracef("Config: %s", appConfig.Json())
}
