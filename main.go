package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/smart7even/golang-do/internal/app"
)

func init() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Can't load .env file")
		}
	}
}

func main() {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	httpAddress := os.Getenv("HTTP_ADRESS")
	grpcAdress := os.Getenv("GRPC_ADRESS")
	secret := os.Getenv("SECRET")

	if secret == "" {
		log.Fatal("SECRET is not set")
	}

	time.Sleep(time.Second * 2)

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("mysql://%s", dbConnectionString))

	if err != nil {
		log.Fatalf("Unable to prepare migrations: %v", err)
	}

	log.Printf("Trying to apply migrations")
	err = m.Up()

	if err != nil {
		if strings.Contains(err.Error(), "connect: connection refused") {
			log.Fatalf("Connection refused: %v", err)
		} else if err == migrate.ErrNoChange {
			log.Printf("Database already up-to-date")
		} else {
			log.Fatalf("Unable to make migrations: %v", err)
		}
	} else {
		log.Printf("Migrations to db applied")
	}

	app.Run(dbConnectionString, httpAddress, grpcAdress, secret)
}
