package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	helpers "backend/db-helpers"
	logger "backend/logger"
)

const BORROW_DAYS = 14

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

	r.POST("/returnwithfine", returnBookWithFine)

	r.DELETE("/delete", deleteBook)

	r.GET("/logs", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, logger.GetLogs())
	})

	port := os.Getenv("PORT")
	if port == "" {port = "8080"}
	r.Run(":" + port)
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
	err := helpers.Borrow(book, BORROW_DAYS)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	} else {
		err = logger.ParseData(book, "borrow")
		if err != nil {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}

func returnBook (c *gin.Context) {
	var book helpers.Book
	c.BindJSON(&book)
	err := helpers.ReturnBook(book, BORROW_DAYS)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	} else {
		err = logger.ParseData(book, "return")
		if err != nil {
			c.IndentedJSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
		}
	}
}

func returnBookWithFine (c *gin.Context) {
	var book helpers.Book
	c.BindJSON(&book)
	err := helpers.ReturnBookWithFine(book)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
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
