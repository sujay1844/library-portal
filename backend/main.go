package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	helpers "mongodb-test/db-helpers"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/all", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, helpers.FindAll())
	})

	r.GET("/find/:name", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, helpers.Find(c.Param("name")))
	})

	r.POST("/add", addBook)

	r.POST("/borrow", borrowBook)

	r.POST("/return", returnBook)

	r.DELETE("/delete", deleteBook)

	r.Run()
}

func addBook (c *gin.Context) {
	var book helpers.Book
	c.BindJSON(&book)
	book, err := helpers.Insert(book)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": book.Name + " already exists"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Added successfully"})
	}
}

func borrowBook (c *gin.Context) {
	var book helpers.Book
	c.BindJSON(&book)
	name := book.Name
	err := helpers.Borrow(name)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "An error occured",
			"error": err.Error(),
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": book.Name + " was borrowed"})
	}
}

func returnBook (c *gin.Context) {
	var book helpers.Book
	c.BindJSON(&book)
	name := book.Name
	err := helpers.ReturnBook(name)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": book.Name + " was returned"})
	}
}

func deleteBook (c *gin.Context) {
	var book helpers.Book
	c.BindJSON(&book)
	name := book.Name
	count, err := helpers.DeleteBook(name)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d books were deleted", count),
		})
	}
}
