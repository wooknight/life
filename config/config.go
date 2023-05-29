package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port     string `mapstructure:"PORT"`
	DBUrl    string `mapstructure:"DB_URL"`
	DBDriver string `mapstructure:"DB_DRIVER"`
}

func LoadAppConfig(path string) (AppConfig, error) {
	if path == "" {
		return AppConfig{}, nil
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return AppConfig{}, fmt.Errorf("config file not found: %s", path)
		}
		return AppConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return config, nil
}
