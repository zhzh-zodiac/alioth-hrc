package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName  string
	AppEnv   string
	HTTPPort int

	MySQL MySQLConfig
	Redis RedisConfig
}

type MySQLConfig struct {
	DSN                string
	MaxIdleConns       int
	MaxOpenConns       int
	ConnMaxLifetimeMin int
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	_ = v.ReadInConfig()

	cfg := &Config{
		AppName:  v.GetString("APP_NAME"),
		AppEnv:   v.GetString("APP_ENV"),
		HTTPPort: v.GetInt("HTTP_PORT"),
		MySQL: MySQLConfig{
			DSN:                v.GetString("MYSQL_DSN"),
			MaxIdleConns:       v.GetInt("MYSQL_MAX_IDLE_CONNS"),
			MaxOpenConns:       v.GetInt("MYSQL_MAX_OPEN_CONNS"),
			ConnMaxLifetimeMin: v.GetInt("MYSQL_CONN_MAX_LIFETIME_MIN"),
		},
		Redis: RedisConfig{
			Addr:     v.GetString("REDIS_ADDR"),
			Username: v.GetString("REDIS_USERNAME"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "alioth-hrc")
	v.SetDefault("APP_ENV", "dev")
	v.SetDefault("HTTP_PORT", 8080)

	v.SetDefault("MYSQL_MAX_IDLE_CONNS", 10)
	v.SetDefault("MYSQL_MAX_OPEN_CONNS", 50)
	v.SetDefault("MYSQL_CONN_MAX_LIFETIME_MIN", 30)

	v.SetDefault("REDIS_ADDR", "127.0.0.1:6379")
	v.SetDefault("REDIS_DB", 0)
}

func (c *Config) Validate() error {
	if c.MySQL.DSN == "" {
		return fmt.Errorf("MYSQL_DSN is required")
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}
	if c.HTTPPort <= 0 {
		return fmt.Errorf("HTTP_PORT must be positive")
	}
	return nil
}

func (c *Config) MySQLConnMaxLifetime() time.Duration {
	return time.Duration(c.MySQL.ConnMaxLifetimeMin) * time.Minute
}
