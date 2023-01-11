// Package config реализует необходимые структуры и парсинг конфигурации сервера
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
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
func ParseConfig() (*Config, error) {
	var cfg Config

	// Тестовый код удалю в конце
	// pwd, _ := os.Getwd()
	// crtFileName := filepath.Join(pwd, "certs", "client.crt")
	// keyFileName := filepath.Join(pwd, "certs", "client.key")
	// caFileName := filepath.Join(pwd, "certs", "ExampleCA.crt")

	flag.StringVar(&cfg.Address, "addr", "https://localhost", "address in the format host:port")
	flag.StringVar(&cfg.CryptoCrt, "crypto-crt", "", "*.crt filepath")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "*.key filepath")
	flag.StringVar(&cfg.CACertFile, "crypto-ca", "", "*.key filepath for ca")
	flag.StringVar(&cfg.GRPCAddress, "grp-addr", ":3200", "grpc address in the format host:port")

	// Сначала парсим командную строку
	flag.Parse()

	// Далее приоритетно аргументы из ENV
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
