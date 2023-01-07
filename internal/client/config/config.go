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

// Config структура для работы с конфигурацией клиента
type Config struct {
	Address     string `env:"ADDRESS"`
	GRPCAddress string `env:"GRPC_ADDRESS"`
	CryptoCrt   string `env:"CRYPTO_CERT"`
	CryptoKey   string `env:"CRYPTO_KEY"`
	CACertFile  string `env:"CA_CERT_KEY"`
}

// ParseConfig парсинг ENV и командной строки для получения конфигурации
func ParseConfig() *Config {
	var cfg Config

	// Тестовый код удалю в конце
	pwd, _ := os.Getwd()
	crtFileName := filepath.Join(pwd, "certs", "client.crt")
	keyFileName := filepath.Join(pwd, "certs", "client.key")
	caFileName := filepath.Join(pwd, "certs", "ExampleCA.crt")

	flag.StringVar(&cfg.Address, "a", "https://localhost", "address in the format host:port")

	flag.StringVar(&cfg.CryptoCrt, "c", crtFileName, "*.crt filepath for tls")
	flag.StringVar(&cfg.CryptoKey, "k", keyFileName, "*.key filepath for tls")
	flag.StringVar(&cfg.CACertFile, "ca", caFileName, "*.key filepath for tls")
	flag.StringVar(&cfg.GRPCAddress, "g", ":3200", "grpc address in the format host:port")

	//flag.StringVar(&cfg.CryptoCrt, "c", "", "crt filepath for tls")
	//flag.StringVar(&cfg.CryptoKey, "k", "", "key filepath for tls")
	//flag.StringVar(&cfg.CACertFile, "ca", "", "crt filepath for ca")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal().Err(err)
	}
	return &cfg
}
