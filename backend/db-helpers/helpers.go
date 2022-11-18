package db_helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

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
	BorrowDate string `json:"borrow-date" bson:"borrow-date"`
	BorrowedBy string `json:"borrowed-by" bson:"borrowed-by"`
	ReturnDate string `json:"return-date" bson:"return-date"`
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
	if hostname == "" {
		URL = "mongodb://localhost:27017"
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
	if err != nil { return book, err }

	book.ID = res.InsertedID.(primitive.ObjectID)
	return book, nil
}

func Insert(book Book) (Book, error) {
	book.ID = primitive.NewObjectID()
	book.Available = true
	book.BorrowedBy = "none"
	book.ReturnDate = "none"
	book.BorrowDate = "none"
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

func borrow(ctx context.Context, collection *mongo.Collection, book Book, days int) error {
	filter := bson.M{
		"name": book.Name,
	}
	borrowYear, borrowMonth, borrowDay := time.Now().Local().Date()
	returnYear, returnMonth, returnDay := time.Now().Local().AddDate(0, 0, days).Date()
	update := bson.M{
		"$set": bson.M{
			"available": false,
			"borrowed-by": book.BorrowedBy,
			"return-date": fmt.Sprintf("%d-%d-%d", returnYear, returnMonth, returnDay),
			"borrow-date": fmt.Sprintf("%d-%d-%d", borrowYear, borrowMonth, borrowDay),
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func Borrow(book Book, days int) error {
	err := validateBorrow(book)
	if err != nil { return err }

	return borrow(ctx, collection, book, days)
}

func validateBorrow(book Book) error {
	books := Find(book.Name)
	if book.BorrowedBy == "none" || book.BorrowedBy == "" {
		return errors.New("borrower name not provided")
	}
	if len(books) == 0 {
		return errors.New("book not found")
	} else if len(books) > 1 {
		return errors.New("given name matched to more than one book")
	} else {
		if !(books[0].Available) {
			return errors.New("book is not available")
		}
	}
	return nil
}

func returnBook(ctx context.Context, collection *mongo.Collection, book Book) error {
	filter := bson.M{
		"name": book.Name,
	}
	update := bson.M{
		"$set": bson.M{
			"available": true,
			"borrow-date": "none",
			"return-date": "none",
			"borrowed-by": "none",
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func ReturnBook(book Book, days int) error {
	err := validateReturn(book)
	if err != nil { return err }

	delayInDays, err := isReturnDelayed(book, days)
	if err != nil { return err }

	if delayInDays > 0 {
		return errors.New(fmt.Sprintf("Return delayed by %d", delayInDays))
	}
	return returnBook(ctx, collection, book)
}

func ReturnBookWithFine(book Book) error {
	err := validateReturn(book)
	if err != nil { return err }

	return returnBook(ctx, collection, book)
}

func validateReturn(book Book) error {
	books := Find(book.Name)
	if len(books) == 0 {
		return errors.New("book not found")
	} else if len(books) > 1 {
		return errors.New("given name matched to more than one book")
	} else {
		if books[0].Available {
			return errors.New("book is already available and doesn't need to be returned")
		}
	}
	return nil
}

func isReturnDelayed(book Book, days int) (int, error) {
	borrowObj, err := time.Parse("2006-1-2", book.BorrowDate)
	if err != nil { return 0, err }

	diff := int(time.Now().Local().Sub(borrowObj).Hours()/24)
	if diff > days { return diff-days, nil }

	return 0, nil
}

func newID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func handle(err error) {
	if err != nil { log.Print("Error: ", err) }
}
