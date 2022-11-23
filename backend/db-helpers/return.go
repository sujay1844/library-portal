package db_helpers

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func returnBook(ctx context.Context, collection *mongo.Collection, book Book) error {
	filter := bson.M{
		"name": book.Name,
	}
	update := bson.M{
		"$set": bson.M{
			"available": true,
			"return-date": "none",
			"student": bson.M{
				"name": "",
				"std": 0,
				"sec": "",
			},
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
		return errors.New(fmt.Sprintf("Return delayed by %d days", delayInDays))
	}
	return returnBook(ctx, collection, book)
}

func ReturnBookWithFine(book Book) error {
	err := validateReturn(book)
	if err != nil { return err }

	return returnBook(ctx, collection, book)
}
