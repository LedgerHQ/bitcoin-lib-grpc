package config

import (
	"github.com/spf13/viper"
)

// LoadConfigProvider returns a configured viper instance
func LoadProvider(appName string) *viper.Viper {
	return readViperConfig(appName)
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	// global defaults

	v.SetDefault("json_logs", false)
	v.SetDefault("loglevel", "debug")

	return v
}
