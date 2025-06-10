package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gbrayhan/microservices-go/src/application/processor"
	"github.com/gbrayhan/microservices-go/src/infrastructure/db"
	"github.com/gbrayhan/microservices-go/src/infrastructure/kafka"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRecipeRepository(database)

	recipeProcessor := processor.NewRecipeProcessor(repo)

	consumer, err := kafka.NewConsumer(
		[]string{"localhost:9092"},
		"recipes",
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
