package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/gbrayhan/microservices-go/src/application/processor"
	"github.com/gbrayhan/microservices-go/src/infrastructure/db"
	"github.com/gbrayhan/microservices-go/src/infrastructure/kafka"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/routes"
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

	log.Println("Start loading kafka brokers...")

	repo := repository.NewRecipeRepository(database)

	recipeProcessor := processor.NewRecipeProcessor(repo)

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
}
