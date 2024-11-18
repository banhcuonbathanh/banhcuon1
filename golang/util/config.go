package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL  string `mapstructure:"databaseURL"`
	GRPCAddress  string `mapstructure:"GRPCAddress"`
	ServerAddress  string `mapstructure:"ServerAddress"`
	AnthropicAPIKey string `mapstructure:"anthropicAPIKey"`
	AnthropicAPIURL string `mapstructure:"anthropicAPIURL"`
	QuanAnAddress string `mapstructure:"QuanAnAddress"`
	QuanAnJWTsecretKey string `mapstructure:"QuanAnJWTsecretKey"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// cfg, err := config.Load()
// if err != nil {
// 	log.Fatalf("Failed to load config: %v", err)
// }