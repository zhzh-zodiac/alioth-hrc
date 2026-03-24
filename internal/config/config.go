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

	JWTSecret          string
	JWTAccessTTLMin    int
	JWTRefreshTTLHours int

	// CORSAllowOrigins 逗号分隔，例如 http://localhost:5173,http://127.0.0.1:5173；空则使用开发默认。
	CORSAllowOrigins string

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
		AppName:            v.GetString("APP_NAME"),
		AppEnv:             v.GetString("APP_ENV"),
		HTTPPort:           v.GetInt("HTTP_PORT"),
		JWTSecret:          v.GetString("JWT_SECRET"),
		JWTAccessTTLMin:    v.GetInt("JWT_ACCESS_TTL_MIN"),
		JWTRefreshTTLHours: v.GetInt("JWT_REFRESH_TTL_HOURS"),
		CORSAllowOrigins:  v.GetString("CORS_ALLOW_ORIGINS"),
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
	v.SetDefault("JWT_ACCESS_TTL_MIN", 15)
	v.SetDefault("JWT_REFRESH_TTL_HOURS", 168)

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
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWTSecret) < 16 {
		return fmt.Errorf("JWT_SECRET must be at least 16 characters")
	}
	if c.JWTAccessTTLMin <= 0 {
		return fmt.Errorf("JWT_ACCESS_TTL_MIN must be positive")
	}
	if c.JWTRefreshTTLHours <= 0 {
		return fmt.Errorf("JWT_REFRESH_TTL_HOURS must be positive")
	}
	return nil
}

func (c *Config) MySQLConnMaxLifetime() time.Duration {
	return time.Duration(c.MySQL.ConnMaxLifetimeMin) * time.Minute
}

// CORSOriginList 返回允许的浏览器 Origin 列表。
func (c *Config) CORSOriginList() []string {
	s := strings.TrimSpace(c.CORSAllowOrigins)
	if s == "" {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	return out
}
