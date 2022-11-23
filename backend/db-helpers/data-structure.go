package db_helpers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID primitive.ObjectID `json:"-" bson:"_id"` // Don't send ObjectID to client
	Name string `json:"name"`
	Author string `json:"author"`
	Available bool `json:"available"`
	Student Student `json:"student" bson:"student"`
	ReturnDate string `json:"return-date" bson:"return-date"`
}

type Student struct {
	Name string `json:"name" bson:"name"`
	Std int `json:"std" bson:"std"`
	Sec string `json:"sec" bson:"sec"`
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
