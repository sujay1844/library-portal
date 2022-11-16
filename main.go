package main

import (
	"context"
	"fmt"
	// "encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const URL = "mongodb://localhost:27017"

type Book struct {
	ID primitive.ObjectID `json:"-" bson:"_id"` // Don't send ObjectID to client
	Name string `json:"name"`
	Author string `json:"author"`
	Available bool `json:"available"`
}

var ctx context.Context
var collection *mongo.Collection
func init() {
	ctx = context.Background()

	opts := options.Client().ApplyURI(URL)

	clnt, err := mongo.Connect(ctx, opts)
	handle(err)

	log.Printf("MongoDB connected on %s", URL)

	collection = clnt.Database("test").Collection("data")
	log.Printf("Collection '%s' is selected", "data")
}

func main() {
	macbeth := Book{
		ID: newID(), 
		Name: "Macbeth", 
		Author: "Shakespeare",
		Available: true,
	}
	macbeth, err := insert(ctx, macbeth)
	handle(err)
	warAndPeace := Book{
		ID: newID(), 
		Name: "War and Peace", 
		Author: "Leo Tolstoy",
		Available: false,
	}
	warAndPeace, err = insert(ctx, warAndPeace)
	handle(err)
	// log.Print(macbeth)
	// log.Print(findAll(ctx))
	// log.Print(find(ctx, "Macbeth"))
	// jsonObj, err := json.Marshal(macbeth)
	// handle(err)
	// log.Print(string(jsonObj))

	ok := http.StatusOK
	r := gin.Default()
	r.GET("/", func (c *gin.Context) {
		c.IndentedJSON(ok, macbeth)
	})
	r.GET("/all", func(c *gin.Context) {
		c.IndentedJSON(ok, findAll(ctx))
	})
	r.GET("/find/:name", func(c *gin.Context) {
		c.IndentedJSON(ok, find(ctx, c.Param("name")))
	})
	r.POST("/add", func(c *gin.Context) {
		var book Book
		c.BindJSON(&book)
		book.ID = primitive.NewObjectID()
		book.Available = true
		book, err := insert(ctx, book)
		if err != nil {
			c.IndentedJSON(ok, gin.H{"message": book.Name + " already exists"})
		} else {
			c.IndentedJSON(ok, gin.H{"message": "Added successfully"})
		}
	})
	r.POST("/borrow", func(c *gin.Context) {
		var book Book
		c.BindJSON(&book)
		name := book.Name
		err := borrow(ctx, name)
		if err != nil {
			c.IndentedJSON(ok, gin.H{
				"message": "An error occured",
				"error": err.Error(),
			})
		} else {
			c.IndentedJSON(ok, gin.H{"message": book.Name + " was borrowed"})
		}
	})
	r.POST("/return", func(c *gin.Context) {
		var book Book
		c.BindJSON(&book)
		name := book.Name
		err := returnBook(ctx, name)
		if err != nil {
			c.IndentedJSON(ok, gin.H{
				"message": err.Error(),
			})
		} else {
			c.IndentedJSON(ok, gin.H{"message": book.Name + " was returned"})
		}
	})
	r.DELETE("/delete", func(c *gin.Context) {
		var book Book
		c.BindJSON(&book)
		name := book.Name
		count, err := deleteBook(ctx, name)
		if err != nil {
			c.IndentedJSON(ok, gin.H{
				"message": err.Error(),
			})
		} else {
			c.IndentedJSON(ok, gin.H{
				"message": fmt.Sprintf("%d books were deleted", count),
			})
		}
	})
	r.Run()
}

func insert(ctx context.Context, book Book) (Book, error) {
	// Check for duplicates
	if find(ctx, book.Name) != nil {
		return book, errors.New(book.Name + " already exists")
	}

	res, err := collection.InsertOne(ctx, book)
	if err != nil {
		return book, err
	}

	book.ID = res.InsertedID.(primitive.ObjectID)
	return book, nil
}

func findAll(ctx context.Context) []Book {
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

func find(ctx context.Context, name string) []Book {
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

func deleteBook(ctx context.Context, name string) (int, error) {
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

func borrow(ctx context.Context, name string) error {
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

func returnBook(ctx context.Context, name string) error {
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

func newID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func handle(err error) {
	if err != nil {
		log.Print("Error: ", err)
	}
}
