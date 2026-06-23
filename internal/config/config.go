// Package config provides configuration loading for the auth service.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
type Config struct {
	HTTP  HTTPConfig  `mapstructure:"http"`
	GRPC  GRPCConfig  `mapstructure:"grpc"`
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
	Log   LogConfig   `mapstructure:"log"`
	JWT   JWTConfig   `mapstructure:"jwt"`
}

// HTTPConfig holds HTTP server configuration.
type HTTPConfig struct {
	Port string `mapstructure:"port"`
}

// GRPCConfig holds gRPC server configuration.
type GRPCConfig struct {
	Port string `mapstructure:"port"`
}

// DBConfig holds database connection configuration.
type DBConfig struct {
	URL string `mapstructure:"url"`
}

// RedisConfig holds Redis connection configuration.
type RedisConfig struct {
	URL string `mapstructure:"url"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
}

// JWTConfig holds JWT token configuration.
type JWTConfig struct {
	PrivateKeyPath string `mapstructure:"private_key_path"`
	PublicKeyPath  string `mapstructure:"public_key_path"`
	AccessTTL      string `mapstructure:"access_ttl"`
	RefreshTTL     string `mapstructure:"refresh_ttl"`
}

// Load reads configuration from a YAML file and environment variables.
// If the file does not exist, defaults are used.
func Load(path string) (*Config, error) {
	v := viper.New()

	// Set defaults.
	v.SetDefault("http.port", "8080")
	v.SetDefault("grpc.port", "50051")
	v.SetDefault("db.url", "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable")
	v.SetDefault("redis.url", "redis://localhost:6379/0")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.pretty", true)
	v.SetDefault("jwt.access_ttl", "15m")
	v.SetDefault("jwt.refresh_ttl", "720h")

	// Configure Viper to read the YAML file.
	v.SetConfigFile(path)
	v.SetConfigType("yml")

	// Configure environment variable overrides.
	v.SetEnvPrefix("AUTH")
	v.AutomaticEnv()

	// Read the config file. If it doesn't exist, that's fine — defaults apply.
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
