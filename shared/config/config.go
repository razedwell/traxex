package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	ImageTag string `mapstructure:"image_tag"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetDefault("server.port", 8080)
	v.SetDefault("database.host", "localhost")

	v.AddConfigPath(path)
	v.SetConfigName(".env")
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %v\n", err)
	}

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
