package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	PipefyPipeID string
}

func Load() (Config, error) {
	values, err := godotenv.Read(".env")
	if err != nil {
		return Config{}, fmt.Errorf("erro ao ler .env: %w", err)
	}

	cfg := Config{
		DatabaseURL:  values["DATABASE_URL"],
		DBHost:       values["DB_HOST"],
		DBPort:       values["DB_PORT"],
		DBUser:       values["DB_USER"],
		DBPassword:   values["DB_PASSWORD"],
		DBName:       values["DB_NAME"],
		DBSSLMode:    values["DB_SSLMODE"],
		PipefyPipeID: values["PIPEFY_PIPE_ID"],
	}

	if cfg.PipefyPipeID == "" {
		return Config{}, fmt.Errorf("variável de ambiente PIPEFY_PIPE_ID não configurada")
	}

	if cfg.DatabaseURL != "" {
		return cfg, nil
	}

	required := map[string]string{
		"DB_HOST":     cfg.DBHost,
		"DB_PORT":     cfg.DBPort,
		"DB_USER":     cfg.DBUser,
		"DB_PASSWORD": cfg.DBPassword,
		"DB_NAME":     cfg.DBName,
		"DB_SSLMODE":  cfg.DBSSLMode,
	}
	for key, value := range required {
		if value == "" {
			return Config{}, fmt.Errorf("variável de ambiente %s não configurada", key)
		}
	}

	return cfg, nil
}

func (c Config) PostgresConnString() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}
