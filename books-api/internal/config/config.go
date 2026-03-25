package config

import "github.com/spf13/viper"

type Config struct {
	Database DatabaseConfig
	App      AppConfig
	JWT      JWTConfig
	Cache    CacheConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
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
	Host     string
	Port     int
	Password string
	DB       int
	Enabled  bool
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
