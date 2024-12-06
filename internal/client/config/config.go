// Package config implements the necessary structures and parsing of the server configuration
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
)

// Config structure for working with client configuration
type Config struct {
	Address     string `env:"ADDRESS"`
	GRPCAddress string `env:"GRPC_ADDRESS"`
	CryptoCrt   string `env:"CRYPTO_CERT"`
	CryptoKey   string `env:"CRYPTO_KEY"`
	CACertFile  string `env:"CA_CERT_KEY"`
}

// ParseConfig parsing ENV and command line arguments to get the configuration
func ParseConfig() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.Address, "addr", "https://localhost", "address in the format host:port")
	flag.StringVar(&cfg.CryptoCrt, "crypto-crt", "", "*.crt filepath")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "*.key filepath")
	flag.StringVar(&cfg.CACertFile, "crypto-ca", "", "*.key filepath for ca")
	flag.StringVar(&cfg.GRPCAddress, "grp-addr", ":3200", "grpc address in the format host:port")

	// First, we parse the command line arguments
	flag.Parse()

	// Then, we give priority to the arguments from ENV
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
