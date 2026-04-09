package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppName      string
	DbUrl        string
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
}

func LoadConfig() (*Config, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "starlink_consumer"
	}
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		kafkaTopic = "user.created"
	}
	kafkaGroupID := os.Getenv("KAFKA_GROUP_ID")
	if kafkaGroupID == "" {
		kafkaGroupID = "starlink_consumer"
	}

	return &Config{
		AppName:      appName,
		DbUrl:        dbUrl,
		KafkaBrokers: strings.Split(kafkaBrokers, ","),
		KafkaTopic:   kafkaTopic,
		KafkaGroupID: kafkaGroupID,
	}, nil
}
