package types

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	User        string `env:"DBUSER"`
	Password    string `env:"DBPASSWORD"`
	Database    string `env:"DATABASE"`
	Environment string `env:"ENVIRONMENT" envDefault:"DEV"`
	LogDir      string `env:"LOGDIR"`
}

func (c Config) GetUser() string {
	return c.User
}

func (c Config) GetPassword() string {
	return c.Password
}

func (c Config) GetDatabase() string {
	return c.Database
}

func (c Config) GetEnvironment() string {
	return c.Environment
}

func (c Config) GetLogDir() string {
	return c.LogDir
}

// Load required environment variables and put them in the config struct.
func LoadEnv() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err.Error())
	}

	return &cfg
}
