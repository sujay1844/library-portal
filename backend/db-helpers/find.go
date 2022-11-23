package db_helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func find(ctx context.Context, collection *mongo.Collection, name string) []Book {
	filter := bson.M{
		"name": bson.M{
			"$regex": ".*"+name+".*",
			"$options": "i",
		},
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
