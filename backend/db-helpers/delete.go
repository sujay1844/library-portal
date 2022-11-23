package db_helpers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
