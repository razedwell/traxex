package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Auth     AuthConfig     `mapstructure:"auth"`
}
type KafkaConfig struct {
	Brokers           []string `mapstructure:"brokers"`
	GroupID           string   `mapstructure:"group_id"`
	ReplicationFactor int      `mapstructure:"replication_factor"`
}
type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // "debug" | "info" | "warn" | "error"
	Format string `mapstructure:"format"` // "json" | "console"
}
type ServerConfig struct {
	Port           int    `mapstructure:"port"`
	Host           string `mapstructure:"host"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type DatabaseConfig struct {
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	Host         string `mapstructure:"host"`
	Port         uint16 `mapstructure:"port"`
	SSLMode      string `mapstructure:"ssl_mode"`
	PoolMaxConns int32  `mapstructure:"pool_max_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type AuthConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	GRPCPort  int    `mapstructure:"grpc_port"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	// 1. Set defaults from code
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.pool_max_conns", 10)
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.pool_size", 10)
	v.SetDefault("auth.grpc_port", 9090)
	// 2. Load from config.yaml (development base)
	v.AddConfigPath(path)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		// It's OK if config.yaml doesn't exist (use defaults)
		if !strings.Contains(err.Error(), "not found") {
			return nil, fmt.Errorf("failed to read config.yaml: %w", err)
		}
	}

	// 3. Load .env into process env (overrides YAML via BindEnv below).
	// Do NOT use ReadInConfig for .env — it flattens keys and wipes nested YAML like database.password.
	_ = gotenv.Load(filepath.Join(path, ".env"))

	// 4. Env vars (highest priority): APP_* plus docker-style names from root .env
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, binding := range []struct{ key, envKey string }{
		{"database.password", "POSTGRES_PASSWORD"},
		{"database.user", "POSTGRES_USER"},
		{"database.host", "POSTGRES_HOST"},
		{"database.port", "POSTGRES_PORT"},
		{"database.database", "POSTGRES_DB"},
		{"redis.port", "REDIS_PORT"},
	} {
		_ = v.BindEnv(binding.key, binding.envKey)
	}
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshall config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed for cfg: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("database password cannot be empty")
	}

	if c.Database.PoolMaxConns == 0 {
		return fmt.Errorf("database pool max connections must be greater than 0")
	}

	if c.Database.Port == 0 {
		return fmt.Errorf("database port cannot be 0")
	}

	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	if c.Redis.DB < 0 || c.Redis.DB > 15 {
		return fmt.Errorf("redis DB must be between 0 and 15")
	}

	if c.Redis.Port == 0 {
		return fmt.Errorf("redis port cannot be 0")
	}

	if c.Redis.PoolSize == 0 {
		return fmt.Errorf("redis pool size must be greater than 0")
	}

	if c.Auth.GRPCPort < 1 || c.Auth.GRPCPort > 65535 {
		return fmt.Errorf("auth grpc port must be between 1 and 65535")
	}

	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("auth jwt secret cannot be empty")
	}

	return nil
}
