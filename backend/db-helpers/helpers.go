package db_helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	ID primitive.ObjectID `json:"-" bson:"_id"` // Don't send ObjectID to client
	Name string `json:"name"`
	Author string `json:"author"`
	Available bool `json:"available"`
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
	return URL
}

func sampleBooks() {
	books := []Book {
		{
			Name: "War and Peace",
			Author: "Leo Tolstoy",
		},
		{
			Name: "Macbeth",
			Author: "William Shakespeare",
		},
		{
			Name: "The Hobbit",
			Author: "J R R Tolkien",
		},
		{
			Name: "The Adventures of Sherlock Holmes",
			Author: "Sir Arthur Conan Doyle",
		},
	}
	for _, book := range books {
		Insert(book)
	}
}

func insert(ctx context.Context, collection *mongo.Collection, book Book) (Book, error) {
	// Check for duplicates
	if find(ctx, collection, book.Name) != nil {
		return book, errors.New(book.Name + " already exists")
	}

	res, err := collection.InsertOne(ctx, book)
	if err != nil {
		return book, err
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	return book, nil
}

func Insert(book Book) (Book, error) {
	book.ID = primitive.NewObjectID()
	book.Available = true
	return insert(ctx, collection, book)
}

func findAll(ctx context.Context, collection *mongo.Collection) []Book {
	csr, err := collection.Find(ctx, bson.D{})
	handle(err)
	var books []Book
	for csr.Next(ctx) {
		var book Book
		err := csr.Decode(&book)
		handle(err)
		books = append(books, book)
	}
	defer csr.Close(ctx)
	return books
}

func FindAll() []Book {
	return findAll(ctx, collection)
}

func find(ctx context.Context, collection *mongo.Collection, name string) []Book {
	filter := bson.M{
		"name": name,
	}
	csr, err := collection.Find(ctx, filter)
	handle(err)
	var books []Book
	for csr.Next(ctx) {
		var book Book
		err := csr.Decode(&book)
		handle(err)
		books = append(books, book)
	}
	defer csr.Close(ctx)
	return books
}

func Find(name string) []Book {
	return find(ctx, collection, name)
}

func deleteBook(ctx context.Context, collection *mongo.Collection, name string) (int, error) {
	filter := bson.M{
		"name": name,
	}
	csr, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	} else {
		return int(csr.DeletedCount), nil
	}
}

func DeleteBook(name string) (int, error) {
	return deleteBook(ctx, collection, name)
}

func borrow(ctx context.Context, collection *mongo.Collection, name string) error {
	filter := bson.M{
		"name": name,
	}
	update := bson.M{
		"$set": bson.M{
			"available": false,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func Borrow(name string) error {
	return borrow(ctx, collection, name)
}

func returnBook(ctx context.Context, collection *mongo.Collection, name string) error {
	filter := bson.M{
		"name": name,
	}
	update := bson.M{
		"$set": bson.M{
			"available": true,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func ReturnBook(name string) error {
	return borrow(ctx, collection, name)
}

func newID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func handle(err error) {
	if err != nil {
		log.Print("Error: ", err)
	}
}

