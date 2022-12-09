// Package config реализует необходимые структуры и парсинг конфигурации сервера
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address      string `env:"RUN_ADDRESS"`
	DatabaseConn string `env:"DATABASE_URI"`
	SigningKey   string `env:"SUPER_KEY"`
}

func ParseConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "localhost:8085", "address in the format host:port")
	flag.StringVar(&cfg.DatabaseConn, "d", "user=postgres password=12345 host=localhost port=5433 dbname=gophkeep", "connection string to postgresql")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &cfg
}
