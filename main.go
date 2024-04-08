package main

import (
	"github.com/gofiber/fiber/v2"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	books = append(books, Book{ID: 1, Title: "Suntiparb", Author: "Pae"})
	books = append(books, Book{ID: 2, Title: "MM", Author: "Pae"})

	// print(books)

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)

	app.Listen("localhost:8080")
}

func getBooks(c *fiber.Ctx) error {
	return c.JSON(books)
}

func getBook(c *fiber.Ctx) error {
	bookId := c.Params("id")
	return c.SendString(bookId)
	// return c.JSON(books[0])
}
