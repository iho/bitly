package conf

import (
	"os"

	"go.uber.org/fx"
)

type Config struct {
	Host, Port, User, Password, DbName string
	HostName                           string
}

func DefaultConfig() *Config {
	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		HostName: "0.0.0.0:8080",
	}
}

// Module for go fx
var Module = fx.Options(fx.Provide(DefaultConfig))
