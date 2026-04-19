package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName  string `env:"APP_NAME"`
	HttpPort string `env:"HTTP_PORT"`
	DbUrl    string `env:"DATABASE_URL"`
}

func LoadConfig() (*Config, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "starlink_producer"
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	return &Config{
		DbUrl:    dbUrl,
		AppName:  appName,
		HttpPort: httpPort,
	}, nil
}
