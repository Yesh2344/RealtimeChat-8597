package main

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration values.
type Config struct {
	Port     int    `mapstructure:"port"`
	LogLevel string `mapstructure:"log_level"`
}

// LoadConfig reads configuration from environment variables and optional config.yaml.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv() // read from .env

	// Default values
	viper.SetDefault("port", 8080)
	viper.SetDefault("log_level", "info")

	// Load config file if present (ignore error if not found)
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}