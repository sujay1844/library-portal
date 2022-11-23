package logger

import (
	// "context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
)

type Query struct {
	Name string
	Date string
	Action string
	Student string
	Std int
	Sec string
	StartDate string `json:"start-date" bson:"start-date"`
	EndDate string `json:"end-date" bson:"end-date"`
}

func FilterByDate(query Query) []Entry {
	startDate, err := time.Parse("2006-01-02", query.StartDate)
	handle(err)
	endDate, err := time.Parse("2006-01-02", query.EndDate)
	handle(err)
	filter := []bson.M{
		{
			"$match": bson.M{
				"date": bson.M{
					"$gte" : primitive.NewDateTimeFromTime(startDate),
					"$lte" : primitive.NewDateTimeFromTime(endDate),
				},
			},
		},
		// {
		// 	"$sort": bson.M{
		// 		"date": 1,
		// 	},
		// },
	}
	csr, err := collection.Aggregate(ctx, filter)
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
