package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	KafkaConfig  KafkaConfig
	CacheConfig  CacheConfig
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
	PgSSLMode  string `env:"PGSSLMODE"`
}

type ServerConfig struct {
	HTTPPort string `env:"HTTP_PORT"`
}

type KafkaConfig struct {
	Brokers []string `env:"KAFKA_BROKERS"`
	Topic   string   `env:"KAFKA_TOPIC"`
	GroupID string   `env:"KAFKA_GROUP_ID"`
}

type CacheConfig struct {
	PreloadLimit int `env:"CACHE_PRELOAD_LIMIT"`
}

func (s *ServerConfig) Address() string {
	return fmt.Sprintf("localhost:%s", s.HTTPPort)
}

func New() *Config {
	cfg := &Config{}

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}
