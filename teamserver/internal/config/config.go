package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

const defaultConfigPath = "config.toml"

type AuthorizedUser struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Config struct {
	Port            uint16           `toml:"port"`
	Https           bool             `toml:"https"`
	CertificatePath string           `toml:"cert"`
	KeyPath         string           `toml:"key"`
	AuthorizedUsers []AuthorizedUser `toml:"authorized"`
}

func ParseConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	var config Config
	_, err = toml.NewDecoder(configFile).Decode(&config)
	return &config, err
}

func ParseDefaultConfig() (*Config, error) {
	_, err := os.Stat(defaultConfigPath)
	if errors.Is(err, os.ErrNotExist) {
		return &Config{
			Port:  8080,
			Https: false,
			AuthorizedUsers: []AuthorizedUser{
				{Username: "osmium", Password: "I'm losing sanity"},
			},
		}, nil
	} else if err != nil {
		return nil, err
	}

	return ParseConfig(defaultConfigPath)
}

func ValidateConfig(config *Config) error {
	if config.Https {
		if len(config.CertificatePath) == 0 || len(config.KeyPath) == 0 {
			return errors.New("Certificate or key path not specified despite https mode enabled")
		}
	}

	return nil
}
