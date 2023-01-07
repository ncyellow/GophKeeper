// Package config реализует необходимые структуры и парсинг конфигурации сервера
package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
)

// Config структура для работы с конфигурацией сервера
type Config struct {
	Address      string `env:"RUN_ADDRESS"`
	GRPCAddress  string `env:"GRPC_ADDRESS"`
	DatabaseConn string `env:"DATABASE_URI"`
	SigningKey   string `env:"SUPER_KEY"`
	CryptoCrt    string `env:"CRYPTO_CERT"`
	CryptoKey    string `env:"CRYPTO_KEY"`
}

// ParseConfig парсинг ENV + командной строки для чтения конфигурации
func ParseConfig() *Config {
	var cfg Config

	pwd, _ := os.Getwd()
	crtFileName := filepath.Join(pwd, "certs", "localhost.crt")
	keyFileName := filepath.Join(pwd, "certs", "localhost.key")

	flag.StringVar(&cfg.Address, "a", ":443", "address in the format host:port")
	flag.StringVar(&cfg.GRPCAddress, "g", ":3200", "grpc address in the format host:port")
	flag.StringVar(&cfg.DatabaseConn, "d", "user=postgres password=12345 host=localhost port=5433 dbname=gophkeep", "connection string to postgresql")
	//flag.StringVar(&cfg.CryptoCrt, "c", "", "*.crt filepath for tls")
	//flag.StringVar(&cfg.CryptoKey, "k", "", "*.key filepath for tls")

	flag.StringVar(&cfg.CryptoCrt, "c", crtFileName, "*.crt filepath for tls")
	flag.StringVar(&cfg.CryptoKey, "k", keyFileName, "*.key filepath for tls")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &cfg
}
