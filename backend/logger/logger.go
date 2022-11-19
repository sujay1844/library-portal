package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	helpers "backend/db-helpers"
)

type Entry struct {
	Name string
	Date string
	Action string
	Student string
	Std int
	Sec string
}

var ctx context.Context
var collection *mongo.Collection

func init() {
	URL := getURL()
	ctx = context.Background()

	opts := options.Client().ApplyURI(URL)

	clnt, err := mongo.Connect(ctx, opts)
	handle(err)

	log.Printf("MongoDB connected on %s", URL)

	collection = clnt.Database("library").Collection("logs")
	log.Printf("Collection '%s' is selected", "logs")
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

func ParseData(book helpers.Book, action string) error {
	entry := Entry{
			Name: book.Name,
			Date: time.Now().Local().Format("2006-01-02"),
			Action: action,
			Student: book.Student.Name,
			Std: book.Student.Std,
			Sec: book.Student.Sec,
		}
	_, err := collection.InsertOne(ctx, entry)
	return err
}

func getLogs(ctx context.Context, collection *mongo.Collection) []Entry {
	csr, err := collection.Find(ctx, bson.M{})
	handle(err)
	var entries []Entry
	for csr.Next(ctx) {
		var entry Entry
		err := csr.Decode(&entry)
		handle(err)
		entries = append(entries, entry)
	}
	defer csr.Close(ctx)
	return entries
}

func GetLogs() []Entry {
	return getLogs(ctx, collection)
}

func handle(err error) {
	if err != nil { log.Print("Error: ", err) }
}
