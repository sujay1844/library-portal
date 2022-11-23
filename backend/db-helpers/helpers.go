package db_helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context
var collection *mongo.Collection

func init() {
	URL := getURL()
	ctx = context.Background()

	opts := options.Client().ApplyURI(URL)

	clnt, err := mongo.Connect(ctx, opts)
	handle(err)

	log.Printf("MongoDB connected on %s", URL)

	collection = clnt.Database("library").Collection("books")
	log.Printf("Collection '%s' is selected", "data")

	sampleBooks()
}

func getURL() string {
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	hostname := os.Getenv("MONGODB_HOSTNAME")
	var URL string
	if username == "" {
		URL = fmt.Sprintf("mongodb://%s:27017", hostname)
	} else {
		URL = fmt.Sprintf("mongodb://%s:%s@%s:27017", username, password, hostname)
	}
	if hostname == "" {
		URL = "mongodb://localhost:27017"
	}
	return URL
}

func handle(err error) {
	if err != nil { log.Print("Error: ", err) }
}
