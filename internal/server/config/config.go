// Package config implements the necessary structures and parsing for server configuration
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
)

// Config structure for working with server configuration
type Config struct {
	Address      string `env:"RUN_ADDRESS"`
	GRPCAddress  string `env:"GRPC_ADDRESS"`
	DatabaseConn string `env:"DATABASE_URI"`
	SigningKey   string `env:"SUPER_KEY"`
	CryptoCrt    string `env:"CRYPTO_CERT"`
	CryptoKey    string `env:"CRYPTO_KEY"`
}

// ParseConfig parsing ENV + command line for reading configuration
func ParseConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Address, "addr", ":443", "address in the format host:port")
	flag.StringVar(&cfg.GRPCAddress, "grpc-addr", ":3200", "grpc address in the format host:port")
	flag.StringVar(&cfg.DatabaseConn, "dns", "", "connection string to postgresql")
	flag.StringVar(&cfg.CryptoCrt, "crypto-crt", "", "*.crt filepath for tls")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "*.key filepath for tls")

	// First, we parse the command line
	flag.Parse()

	// Next, prioritize arguments from ENV
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
