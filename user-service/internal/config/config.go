package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GRPC     GRPCConfig
	JWT      JWTConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

type GRPCConfig struct {
	Port string `env:"GRPC_PORT" env-default:":50051"`
}

type JWTConfig struct {
	SecretKey  string        `env:"JWT_SECRET" env-required:"true"`
	AccessTTL  time.Duration `env:"JWT_ACCESS_TTL" env-default:"15m"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" env-default:"720h"`
}

type PostgresConfig struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
}

type RedisConfig struct {
	Host     string        `env:"REDIS_HOST" env-default:"localhost"`
	Port     string        `env:"REDIS_PORT" env-default:"6379"`
	Password string        `env:"REDIS_PASSWORD" env-default:""`
	DB       int           `env:"REDIS_DB" env-default:"0"`
	UserTTL  time.Duration `env:"REDIS_USER_CACHE_TTL" env-default:"15m"`
}

func Load() (*Config, error) {
	var cfg Config

	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, fmt.Errorf("failed to read .env file: %w", err)
		}
		return &cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read env variables: %w", err)
	}

	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("config error: %s", err))
	}
	return cfg
}
