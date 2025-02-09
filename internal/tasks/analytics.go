package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// List of task types
const (
	TypeRedirectionAnalytics = "redirection:analytics"
)

type AnalyticsPayload struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ShortCode string             `bson:"short_code,omitempty"`
	UserAgent string             `bson:"user_agent,omitempty"`
	Referer   string             `bson:"timestamp,omitempty"`
	Timestamp time.Time          `bson:"referer,omitempty"`
}

// Task Creators
func NewAnalyticsTask(analyticsPayload AnalyticsPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(AnalyticsPayload{
		ShortCode: analyticsPayload.ShortCode,
		UserAgent: analyticsPayload.UserAgent,
		Referer:   analyticsPayload.Referer,
		Timestamp: analyticsPayload.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeRedirectionAnalytics, payload), nil
}

// Task Handlers
func HandleAnalyticsDeliveryTask(ctx context.Context, t *asynq.Task, db *mongo.Client) error {
	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		log.Fatal().Msg("MONGODB_DATABASE is not set")
	}

	var p AnalyticsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.UnMarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Info().Msgf("Recording analytics for payload %v", p)
	analyticsCollection := db.Database(databaseName).Collection("analytics")
	_, err := analyticsCollection.InsertOne(ctx, bson.M{"short_code": p.ShortCode, "timestamp": p.Timestamp, "user_agent": p.UserAgent, "referer": p.Referer})
	if err != nil {
		log.Error().Msgf("Error inserting redirection analytics data %v", err)
	}
	log.Info().Msg("Successfully recorded redireciton analytics data")
	return nil
}
