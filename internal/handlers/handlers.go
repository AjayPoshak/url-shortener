package handlers

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"net/http"
	"time"
)

type Handlers struct {
	db           *mongo.Client
	DatabaseName string
}

func NewHandlers(client *mongo.Client, dbName string) *Handlers {
	return &Handlers{
		db:           client,
		DatabaseName: dbName,
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

func (handler *Handlers) GetUrls(w http.ResponseWriter, r *http.Request) {
	// @TODO: Replace it with actual user ids
	userId := 1
	urlsCollection := handler.db.Database(handler.DatabaseName).Collection("urls")
	// Returns all the urls for this user
	filter := bson.M{"user_id": userId}
	cursor, err := urlsCollection.Find(r.Context(), filter)
	if err != nil {
		log.Printf("Could not retrieve urls for user %d: %v ", userId, err)
		http.Error(w, "Could not retrieve urls", http.StatusInternalServerError)
		return
	}
	// Unpack cursor into slice
	var results []url
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	defer cursor.Close(context.Background())
	fmt.Println("results ", results)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (handler *Handlers) CreateUrl(w http.ResponseWriter, r *http.Request) {
	const maxLength = 2048
	var urlReq URLRequest
	// Validate the request body
	if err := json.NewDecoder(r.Body).Decode(&urlReq); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	hash := md5.New()
	// @TODO: Replace it with actual user ids
	userId := 1
	urlsCollection := handler.db.Database(handler.DatabaseName).Collection("urls")
	// Read request body from POST
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error reading request body"))
		return
	}
	json.Unmarshal(requestBody, &urlReq)
	if urlReq.URL == "" || len(urlReq.URL) > maxLength {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
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
	fmt.Println("Hash: ", string(requestBody), shortURL)
	_, err = urlsCollection.InsertOne(r.Context(), bson.M{"original_url": urlReq.URL, "short_code": shortURL, "created_at": time.Now(), "user_id": userId})
	if err != nil {
		log.Printf("Could not insert new short_code to database %v", err)
		http.Error(w, "Could not insert new short_code to database", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Hash: " + string(hash.Sum(nil))))
}

func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (handler *Handlers) Redirect(response http.ResponseWriter, request *http.Request) {
	shortCode := request.PathValue("shortCode")
	// Get the correspond long URL
	urlsCollection := handler.db.Database(handler.DatabaseName).Collection("urls")
	filter := bson.M{"short_code": shortCode}
  var result url
	err := urlsCollection.FindOne(request.Context(), filter).Decode(&result)
	if err != nil {
    if err == mongo.ErrNoDocuments {
      http.Error(response, "This short URL does not exists", http.StatusNotFound)
      log.Printf("Short URL does not exists %v", shortCode)
      return
    }
		http.Error(response, "Something is not right here", http.StatusInternalServerError)
		return
	}
  // @TODO: Update click count
  http.Redirect(response, request, result.OriginalUrl, http.StatusFound)
}
