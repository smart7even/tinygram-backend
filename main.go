package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/smart7even/golang-do/internal/repository"
	"github.com/smart7even/golang-do/internal/service"
	"github.com/smart7even/golang-do/internal/transport/http_handler"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	address := os.Getenv("ADRESS")

	db, err := sql.Open("mysql", dbConnectionString)

	todoRepo := repository.NewMySQLTodoRepo(db)
	todoService := service.NewTodoService(todoRepo)

	services := service.Services{
		Todo: *todoService,
	}

	if err != nil {
		fmt.Printf("Can't prepare driver to connect to db: %v", err)
		return
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Can't connect to db: %v", err)
		return
	}

	handler := http_handler.Handler{Services: services}
	router := handler.InitAPI()

	router.Run(address)
}
