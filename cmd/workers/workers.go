package main

import (
	"context"
	"os"

	"github.com/AjayPoshak/url-shortener/internal/tasks"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAnalyticsHandler(db *mongo.Client) asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		return tasks.HandleAnalyticsDeliveryTask(ctx, t, db)
	}
}
func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal().Msg("MONGO_URI is not set")
	}
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		log.Fatal().Msg("MONGODB_DATABASE is not set")
	}

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

	redisURI := os.Getenv("REDIS_URI")
	if redisURI == "" {
		log.Fatal().Msg("REDIS_URI is not set")
	}
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisURI},
		asynq.Config{
			Concurrency: 2,
		},
	)
	mux := asynq.NewServeMux()

	mux.HandleFunc(tasks.TypeRedirectionAnalytics, NewAnalyticsHandler(client))

	serverErr := server.Run(mux)
	if serverErr != nil {
		log.Fatal().Msgf("Could not run server %v", serverErr)
	}
}
