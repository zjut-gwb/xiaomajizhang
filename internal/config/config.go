package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type AppConfig struct {
	Env  string `mapstructure:"env"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	Charset  string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret                 string        `mapstructure:"secret"`
	AccessTokenExpiration  time.Duration `mapstructure:"access_token_expiration"`
	RefreshTokenExpiration time.Duration `mapstructure:"refresh_token_expiration"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("configs")
	v.AddConfigPath(".")
	v.SetEnvPrefix("XMJZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("app.env", "dev")
	v.SetDefault("app.host", "0.0.0.0")
	v.SetDefault("app.port", 8080)
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "root")
	v.SetDefault("database.name", "xiaomajizhang")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.access_token_expiration", "2h")
	v.SetDefault("jwt.refresh_token_expiration", "168h")
	v.SetDefault("log.level", "debug")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return &cfg, nil
}
