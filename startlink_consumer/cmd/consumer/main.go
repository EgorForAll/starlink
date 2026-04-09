package main

import (
	"log"

	"starlink_consumer/internal/app"
	"starlink_consumer/internal/config"
	"starlink_consumer/internal/container"
	"starlink_consumer/internal/infra/db"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	configData, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := db.InitDb(configData.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	db.InitSquirrel(dbConn)

	// инициализация зависимостей
	di := container.NewDiContainer(dbConn, configData)
	di.InitDependencies(configData)

	if err := app.InitApp(di); err != nil {
		log.Fatal(err)
	}
}
