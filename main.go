package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/smart7even/golang-do/internal/repository"
	"github.com/smart7even/golang-do/internal/service"
	pb "github.com/smart7even/golang-do/internal/transport/grpc_handler"
	"github.com/smart7even/golang-do/internal/transport/http_handler"
	"google.golang.org/grpc"
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

	go func() {
		router.Run(httpAddress)
	}()

	lis, err := net.Listen("tcp", grpcAdress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, pb.NewTodoGrpcServer(services))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
