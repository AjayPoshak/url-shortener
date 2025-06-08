package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/AjayPoshak/url-shortener/internal/tasks"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handlers struct {
	db           *mongo.Client
	DatabaseName string
	redis        *redis.Client
	queue        *asynq.Client
}

func NewHandlers(client *mongo.Client, dbName string, redis *redis.Client, queueClient *asynq.Client) *Handlers {
	return &Handlers{
		db:           client,
		DatabaseName: dbName,
		redis:        redis,
		queue:        queueClient,
	}
}

type URLRequest struct {
	URL string `json:"url"`
}

type url struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	OriginalUrl string             `bson:"original_url,omitempty" json:"original_url,omitempty"`
	ShortCode   string             `bson:"short_code,omitempty" json:"short_code,omitempty"`
	CreatedAt   primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UserId      int                `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func JSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	errResponse := ErrorResponse{
		Message: err,
		Code:    code,
	}
	json.NewEncoder(w).Encode(&errResponse)
}

func (handler *Handlers) GetUrls(w http.ResponseWriter, r *http.Request) {
	// @TODO: Replace it with actual user ids
	userId := 1
	urlsCollection := handler.db.Database(handler.DatabaseName).Collection("urls")
	// Returns all the urls for this user
	filter := bson.M{"user_id": userId}
	cursor, err := urlsCollection.Find(r.Context(), filter)
	if err != nil {
		log.Error().Msgf("Could not retrieve urls for user %d: %v ", userId, err)
		JSONError(w, "Could not retrieve urls", http.StatusInternalServerError)
		return
	}
	// Unpack cursor into slice
	var results []url
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Error().Msgf("Could not unpack cursor into slice: %v ", err)
		JSONError(w, "Something is not right here. We are looking into it", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

type CreateUrlResponse struct {
	ShortCode string
}

func (handler *Handlers) CreateUrl(w http.ResponseWriter, r *http.Request) {
	const maxLength = 2048
	var urlReq URLRequest
	// Validate the request body
	if err := json.NewDecoder(r.Body).Decode(&urlReq); err != nil {
		log.Error().Msgf("Invalid request body: %v", err)
		JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	hash := md5.New()
	// @TODO: Replace it with actual user ids
	userId := 1
	urlsCollection := handler.db.Database(handler.DatabaseName).Collection("urls")
	// Read request body from POST
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Msgf("Error reading request body: %v", err)
		w.Write([]byte("Error reading request body"))
		return
	}
	json.Unmarshal(requestBody, &urlReq)
	if urlReq.URL == "" || len(urlReq.URL) > maxLength {
		log.Error().Msgf("Invalid URL: %v", urlReq.URL)
		JSONError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	io.WriteString(hash, urlReq.URL+string(userId)) // Added userId along with URL to make the hash unique
	generatedHash := hash.Sum(nil)
	const dictionary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // We want to use lowercase, uppercase letters and numbers
	// Convert hash to base62
	var shortURL string
	for i := 0; i < 8; i++ {
		index := int(generatedHash[i]) % 62
		shortURL += string(dictionary[index])
	}
	log.Info().Str("requestBody ", string(requestBody)).Str("shortURL", shortURL).Msg("Generated Hash ")
	_, err = urlsCollection.InsertOne(r.Context(), bson.M{"original_url": urlReq.URL, "short_code": shortURL, "created_at": time.Now(), "user_id": userId})
	if err != nil {
		log.Error().Msgf("Could not insert new short_code to database %v", err)
		JSONError(w, "Already have this URL in database", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := CreateUrlResponse{
		ShortCode: shortURL,
	}
	encodingErr := json.NewEncoder(w).Encode(&response)
	if encodingErr != nil {
		log.Error().Msgf("Error encoding short URL %v", encodingErr)
		JSONError(w, "Error encoding short URL", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (handler *Handlers) InsertRedirection(redirectionData tasks.AnalyticsPayload, request *http.Request) {
	analyticsCollection := handler.db.Database(handler.DatabaseName).Collection("analytics")
	_, err := analyticsCollection.InsertOne(request.Context(), bson.M{"short_code": redirectionData.ShortCode, "timestamp": redirectionData.Timestamp, "user_agent": redirectionData.UserAgent, "referer": redirectionData.Referer})
	if err != nil {
		log.Error().Msgf("Error inserting redirection analytics data %v", err)
	}
}

var MAX_QUEUE_RETRY = 10
var QUEUE_TIMEOUT = 5 * time.Minute

func (handler *Handlers) Redirect(response http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	redirectionAnalytics := tasks.AnalyticsPayload{
		ShortCode: shortCode,
		UserAgent: request.Header.Get("User-Agent"),
		Referer:   "",
		Timestamp: time.Now(),
	}
	cachedValue, redisErr := handler.redis.Get(context.Background(), shortCode).Result()
	if redisErr != nil {
		log.Error().Msgf("Error in Redis fetching short URL %v", redisErr)
	}
	if cachedValue != "" {
		task, err := tasks.NewAnalyticsTask(redirectionAnalytics)
		if err != nil {
			log.Error().Msgf("Could not enqueue task: %v", err)
		}
		// If the request is a GET request, enqueue the task, HEAD requests should not enqueue tasks because they are for returning header to monitoring services
		if request.Method == http.MethodGet {
			handler.queue.Enqueue(task, asynq.MaxRetry(MAX_QUEUE_RETRY), asynq.Timeout(QUEUE_TIMEOUT))
		}
		log.Info().Msgf("Value found in redis %v", cachedValue)
		http.Redirect(response, request, cachedValue, http.StatusFound)
		return
	}
	// Get the correspond long URL
	urlsCollection := handler.db.Database(handler.DatabaseName).Collection("urls")
	filter := bson.M{"short_code": shortCode}
	var result url
	err := urlsCollection.FindOne(request.Context(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Error().Msgf("Short URL does not exists %v", shortCode)
			JSONError(response, "This short URL does not exists", http.StatusNotFound)
			return
		}
		log.Error().Msgf("Error fetching short URL %v", err)
		JSONError(response, "Something is not right here", http.StatusInternalServerError)
		return
	}
	task, err := tasks.NewAnalyticsTask(redirectionAnalytics)
	if err != nil {
		log.Error().Msgf("Could not enqueue task: %v", err)
	}

	if request.Method == http.MethodGet {
		handler.queue.Enqueue(task, asynq.MaxRetry(MAX_QUEUE_RETRY), asynq.Timeout(QUEUE_TIMEOUT))
	}

	err = handler.redis.Set(context.Background(), shortCode, result.OriginalUrl, 0).Err()
	if err != nil {
		log.Error().Msgf("Error in Redis setting short URL %v", err)
	}
	http.Redirect(response, request, result.OriginalUrl, http.StatusFound)
}
