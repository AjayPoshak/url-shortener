package main

import (
	"context"
	"fmt"
	"github.com/AjayPoshak/url-shortener/internal/handlers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		log.Fatal("MONGODB_DATABASE is not set")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	opts.SetDirect(true)

	// Create a new client and connect to mongo
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm successful connection
	var result bson.M

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Succesfully connected to MongoDB")
	// Create a new router
	router := http.NewServeMux()

	// Initialize the handlers
	handlers := handlers.NewHandlers(client, databaseName)

	// Register routes with middleware
	router.HandleFunc("GET /urls", handlers.GetUrls)
	router.HandleFunc("POST /urls", handlers.CreateUrl)
	router.HandleFunc("GET /health", handlers.HealthHandler)
	router.HandleFunc("GET /{shortCode}", handlers.Redirect)

	port := ":8095"
	// Create a new server
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server
	log.Println("Starting server on ", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// test comment
