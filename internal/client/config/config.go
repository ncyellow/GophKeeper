// Package config реализует необходимые структуры и парсинг конфигурации сервера
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
)

// Config структура для работы с конфигурацией клиента
type Config struct {
	Address string `env:"ADDRESS"`
}

// ParseConfig парсинг ENV и командной строки для получения конфигурации
func ParseConfig() *Config {
	var cfg Config

	flag.StringVar(&cfg.Address, "a", "localhost:8085", "address in the format host:port")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &cfg
}
