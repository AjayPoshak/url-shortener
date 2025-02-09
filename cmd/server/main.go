package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/AjayPoshak/url-shortener/config"
	"github.com/AjayPoshak/url-shortener/internal/handlers"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config := config.NewConfig()
	config.SetEnv(os.Getenv("GO_ENV"))
	config.SetMongoURI(os.Getenv("MONGODB_URI"))
	config.SetDatabaseName(os.Getenv("MONGODB_DATABASE"))
	config.SetRedisURI(os.Getenv("REDIS_URI"))
	if err := config.Validate(); err != nil {
		log.Fatal().Msgf("Failed to validate config: %v", err)
	}
	if config.GetEnv() == "production" {

		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else {

		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	mongoURI := config.GetMongoURI()
	databaseName := config.GetDatabaseName()
	redisURI := config.GetRedisURI()

	redis := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})

	queueClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisURI})
	defer queueClient.Close()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	opts.SetDirect(true)

	// Create a new client and connect to mongo
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal().Msgf("Failed to connect to MongoDB: %v", err)
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal().Msgf("Failed to disconnect from MongoDB: %v", err)
			panic(err)
		}
	}()

	// Send a ping to confirm successful connection
	var result bson.M

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatal().Msgf("Failed to ping MongoDB: %v", err)
		panic(err)
	}
	log.Info().Msg("Succesfully connected to MongoDB")
	// Create a new router
	router := http.NewServeMux()

	// Initialize the handlers
	handlers := handlers.NewHandlers(client, databaseName, redis, queueClient)

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
	log.Info().Msgf("Starting server on %v", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Msgf("Server failed to start: %v", err)
	}
}

// test comment
