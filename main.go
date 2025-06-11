package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/vostelmakh/mixturka/src/application/processor/brew"
	"github.com/vostelmakh/mixturka/src/application/processor/recipe"
	"github.com/vostelmakh/mixturka/src/application/server"
	"github.com/vostelmakh/mixturka/src/infrastructure/db"
	mixturkaGrpc "github.com/vostelmakh/mixturka/src/infrastructure/grpc"
	"github.com/vostelmakh/mixturka/src/infrastructure/kafka"
	"github.com/vostelmakh/mixturka/src/infrastructure/repository"
	"github.com/vostelmakh/mixturka/src/infrastructure/rest/middlewares"
	"github.com/vostelmakh/mixturka/src/infrastructure/rest/routes"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	database, err := db.InitDB()
	if err != nil {
		panic(fmt.Errorf("error initializing the database: %w", err))
	}

	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.GinBodyLogMiddleware)
	router.Use(middlewares.CommonHeaders)

	routes.ApplicationRouter(router)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Server running at http://localhost:%s\n", port)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			panic(strings.ToLower(err.Error()))
		}
	}()

	// Инициализация grpc сервера
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Инициализация процессоров
	repo := repository.NewRecipeRepository(database)

	recipeProcessor := recipe.NewRecipeProcessor(repo)
	brewProcessor := brew.NewGRPCProcessor(repo)

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	mixturkaServer := server.NewMixturkaServer(recipeProcessor, brewProcessor)
	mixturkaGrpc.RegisterMixturkaServer(grpcServer, mixturkaServer)

	go func() {
		fmt.Printf("grpc server running at localhost:%s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Println("Start loading kafka brokers...")

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")

	consumer, err := kafka.NewConsumer(
		[]string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)},
		kafka.TopicBabushkaRecipeV1,
		recipeProcessor,
	)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")

	grpcServer.GracefulStop()
}
