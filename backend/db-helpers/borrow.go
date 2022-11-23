package db_helpers

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func borrow(ctx context.Context, collection *mongo.Collection, book Book, days int) error {
	filter := bson.M{
		"name": book.Name,
	}
	update := bson.M{
		"$set": bson.M{
			"available": false,
			"student": book.Student,
			"return-date": time.Now().Local().AddDate(0, 0, days).Format("2006-01-02"),
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
