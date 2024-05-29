package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectMongo() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
}

func SaveRequest(query, email string, domains []string) error {
	collection := Client.Database("brandscan").Collection("requests")

	doc := map[string]interface{}{
		"query":   query,
		"email":   email,
		"domains": domains,
		"created": time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), doc)
	return err
}
