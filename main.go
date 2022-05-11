package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/smart7even/golang-do/internal/app"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	httpAddress := os.Getenv("HTTP_ADRESS")
	grpcAdress := os.Getenv("GRPC_ADRESS")

	app.Run(dbConnectionString, httpAddress, grpcAdress)
}
