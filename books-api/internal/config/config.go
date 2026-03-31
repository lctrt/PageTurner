package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
	Cache    CacheConfig
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
}

type AppConfig struct {
	Host string
	Port int
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type CacheConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Enabled  bool   `mapstructure:"enabled"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.BindEnv("database.host", "POSTGRES_HOST")
	viper.BindEnv("database.port", "POSTGRES_PORT")
	viper.BindEnv("database.user", "POSTGRES_USER")
	viper.BindEnv("database.password", "POSTGRES_PASSWORD")
	viper.BindEnv("database.database", "POSTGRES_DB")
	viper.BindEnv("database.sslmode", "POSTGRES_SSLMODE")
	viper.BindEnv("cache.host", "REDIS_HOST")
	viper.BindEnv("cache.port", "REDIS_PORT")
	viper.BindEnv("cache.password", "REDIS_PASSWORD")
	viper.BindEnv("cache.db", "REDIS_DB")
	viper.BindEnv("cache.enabled", "REDIS_ENABLED")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
