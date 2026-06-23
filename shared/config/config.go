package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	ImageTag     string `mapstructure:"image_tag"`
	Name         string `mapstructure:"name"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	SSLMode      string `mapstructure:"ssl_mode"`
	PoolMaxConns int    `mapstructure:"pool_max_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	// 1. Set defaults from code
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "5432")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")

	// 2. Load from config.yaml (development base)
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		// It's OK if config.yaml doesn't exist (use defaults)
		if !strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("failed to read config.yaml: %v", err)
		}
	}

	// 3. Load from .env (overrides YAML)
	v.AddConfigPath(path)
	v.SetConfigName(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig() // Don't fail if .env doesn't exist

	// 4. Enable env vars (highest priority)
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall config: %v\n", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("Validation failed for cfg: %v\n", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("Database password cannot be empty")
	}
	return nil
}
