package db_helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func findAll(ctx context.Context, collection *mongo.Collection) []Book {
	csr, err := collection.Find(ctx, bson.M{})
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
