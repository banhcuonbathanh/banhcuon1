package util

// import (
// 	"os"
// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	DatabaseURL         string `mapstructure:"databaseURL"`
// 	GRPCAddress         string `mapstructure:"GRPCAddress"`
// 	ServerAddress       string `mapstructure:"ServerAddress"`
// 	AnthropicAPIKey    string `mapstructure:"anthropicAPIKey"`
// 	AnthropicAPIURL    string `mapstructure:"anthropicAPIURL"`
// 	QuanAnAddress      string `mapstructure:"QuanAnAddress"`
// 	QuanAnJWTsecretKey string `mapstructure:"QuanAnJWTsecretKey"`
// }

// func Load() (*Config, error) {
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("yaml")
// 	viper.AddConfigPath(".")
// 	viper.AutomaticEnv()

// 	if err := viper.ReadInConfig(); err != nil {
// 		return nil, err
// 	}

// 	var cfg Config
// 	if err := viper.Unmarshal(&cfg); err != nil {
// 		return nil, err
// 	}

// 	// Override database URL based on environment
// 	if os.Getenv("APP_ENV") == "docker" {
// 		cfg.DatabaseURL = "postgres://myuser:mypassword@mypostgres_ai:5432/mydatabase?sslmode=disable"
// 	} else {
// 		// Use the local configuration from config.yaml
// 		// This will use the localhost database URL
// 	}

// 	return &cfg, nil
// }
// // cfg, err := config.Load()
// // if err != nil {
// // 	log.Fatalf("Failed to load config: %v", err)
// // }

// new for docker

// config/config.go

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL         string `mapstructure:"databaseURL"`
	GRPCAddress         string `mapstructure:"GRPCAddress"`
	ServerAddress       string `mapstructure:"ServerAddress"`
	AnthropicAPIKey    string `mapstructure:"anthropicAPIKey"`
	AnthropicAPIURL    string `mapstructure:"anthropicAPIURL"`
	QuanAnAddress      string `mapstructure:"QuanAnAddress"`
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

	// Override database URL based on environment
	if os.Getenv("APP_ENV") == "docker" {
		// Construct database URL using environment variables

		fmt.Print("golang/util/config.go func Load() ")
		cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"))
			fmt.Print("golang/util/config.go func Load() 	cfg.DatabaseURL", 	cfg.DatabaseURL)
		// Override GRPC address for Docker
		cfg.GRPCAddress = ":50051"
	}

	return &cfg, nil
}