package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	GRPC     GRPCConfig
	Postgres PostgresConfig
	Storage  StorageConfig
}

type GRPCConfig struct {
	Port string `env:"GRPC_PORT" env-default:":50052"`
}

type PostgresConfig struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
}

type StorageConfig struct {
	Endpoint        string        `env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
	AccessKeyID     string        `env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
	SecretAccessKey string        `env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
	UseSSL          bool          `env:"MINIO_USE_SSL" env-default:"false"`
	BucketName      string        `env:"MINIO_BUCKET_NAME" env-default:"soundwave-audio"`
	UploadExpiry    time.Duration `env:"MINIO_UPLOAD_EXPIRY" env-default:"15m"`
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