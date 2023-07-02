package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smart7even/golang-do/internal/repository"
	"github.com/smart7even/golang-do/internal/service"
	pb "github.com/smart7even/golang-do/internal/transport/grpc_handler"
	"github.com/smart7even/golang-do/internal/transport/http_handler"
	"google.golang.org/grpc"

	firebase "firebase.google.com/go"

	"google.golang.org/api/option"
)

func initFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("firebase.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	return app, nil
}

func Run(dbConnectionString, httpAddress, grpcAdress, secret string) {
	db, err := sql.Open("mysql", dbConnectionString+"?parseTime=true")

	if err != nil {
		fmt.Printf("Unable to connect to db %v", err)
		return
	}

	todoRepo := repository.NewMySQLTodoRepo(db)
	todoService := service.NewTodoService(todoRepo)

	firebaseApp, err := initFirebase()

	if err != nil {
		fmt.Printf("Unable to initialize Firebase %v", err)
		return
	}

	eventRepo := repository.NewMySqlEventRepo(db)
	eventService := service.NewEventService(eventRepo)

	userRepo := repository.NewMySQLUserRepo(db, *firebaseApp)
	userService := service.NewUserService(userRepo)

	chatRepo := repository.NewMySQLChatRepo(db)
	chatService := service.NewChatService(chatRepo)

	messageRepo := repository.NewMySQLMessageRepo(db)
	messageService := service.NewMessageService(messageRepo, eventRepo)

	authService := service.NewAuthService(secret)

	services := service.Services{
		Todo:    *todoService,
		User:    *userService,
		Chat:    *chatService,
		Message: *messageService,
		Auth:    *authService,
		Event:   *eventService,
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

	srv := &http.Server{
		Addr:    httpAddress,
		Handler: router,
	}

	s := grpc.NewServer()
	pb.RegisterTodoServiceServer(s, pb.NewTodoGrpcServer(services))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", grpcAdress)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("failed to stop server: %v", err)
	}

	s.GracefulStop()
}
