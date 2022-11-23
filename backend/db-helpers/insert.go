package db_helpers

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	book.Student = Student{}
	book.ReturnDate = "none"
	return insert(ctx, collection, book)
}
