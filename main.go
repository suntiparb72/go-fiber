package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var testUser = User{
	Email:    "suntiparb",
	Password: "1234",
}

func checkMiddleware(c *fiber.Ctx) error {
	time := time.Now()

	fmt.Printf("URL: %s, Method: %s, Time:%s, \n", c.OriginalURL(), c.Method(), time)

	return c.Next()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Load .env error")
	}

	app := fiber.New()

	app.Post("/login", login)

	app.Use(checkMiddleware)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	books = append(books, Book{ID: 1, Title: "Suntiparb", Author: "Pae"})
	books = append(books, Book{ID: 2, Title: "goFiber", Author: "Pae"})

	app.Get("/books", getBooks)
	app.Get("/books/:id", getBook)
	app.Post("/books", createBook)
	app.Put("/books/:id", updateBook)
	app.Delete("/books/:id", deleteBook)

	app.Post("/upload", uploadFile)

	app.Get("/config", getEnv)

	app.Listen("localhost:8080")
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/"+file.Filename)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File upload complete!")
}

func getEnv(c *fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"secret": os.Getenv("SECRET"),
	})
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != testUser.Email || user.Password != testUser.Password {
		return fiber.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message": "login success",
		"token":   t,
	})
}
