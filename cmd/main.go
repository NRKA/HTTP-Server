package main

import (
	"context"
	"errors"
	"github.com/NRKA/HTTP-Server/internal/db"
	"github.com/NRKA/HTTP-Server/internal/handlers"
	"github.com/NRKA/HTTP-Server/internal/kafka"
	"github.com/NRKA/HTTP-Server/internal/repository/postgresql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	port       = "PORT"
	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbName     = "DB_NAME"
	brokerAddr = "BROKER_ADDRESS"
	topic      = "TOPIC"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	brokerAddress := os.Getenv(brokerAddr)

	producer, err := kafka.NewKafkaProducer(brokerAddress)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("Failed to close producer: %v", err)
		}
	}()

	consumer, err := kafka.NewKafkaConsumer(brokerAddress)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}
	defer func() {
		err := consumer.Close()
		if err != nil {
			log.Printf("Failed to close consumer: %v", err)
		}
	}()

	port := os.Getenv(port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	dbConfig := db.DatabaseConfig{
		Host:     os.Getenv(dbHost),
		Port:     os.Getenv(dbPort),
		User:     os.Getenv(dbUser),
		Password: os.Getenv(dbPassword),
		DBName:   os.Getenv(dbName),
	}
	database, err := db.NewDB(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.GetPool().Close()

	server := &http.Server{Addr: port}
	go func() {
		err := consumer.Consume(ctx, os.Getenv(topic))
		if err != nil {
			log.Printf("failed to consume: %v", err)

			//здесь не закрыл сервер, так как кафка не является главной задачи сервиса
			//сервер будет работать выполняя crud операции
			return
		}

		err = server.Shutdown(context.Background())
		if err != nil {
			log.Printf("failed to stop server: %v", err)
		}
		return
	}()

	articleRepo := postgresql.NewArticleRepo(database)
	handler := handlers.NewArticleHandler(articleRepo, producer)
	http.Handle("/", handlers.CreateRouter(handler))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start the HTTP server: %v", err)
	}
}
