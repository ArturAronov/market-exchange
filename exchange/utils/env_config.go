package utils

import (
	"os"
)

type EnvConfig struct {
	Port      string
	Host      string
	UserAgent string
}

func EnvHandler() EnvConfig {
	envOutput := EnvConfig{
		Port:      os.Getenv("PORT"),
		Host:      os.Getenv("HOST"),
		UserAgent: os.Getenv("USER_AGENT"),
	}

	if envOutput.Port == "" {
		panic("Missing env: PORT")
	}

	if envOutput.Host == "" {
		panic("Missing env: HOST")
	}

	if envOutput.UserAgent == "" {
		panic("Missing env: USER_AGENT")
	}

	return envOutput
}
