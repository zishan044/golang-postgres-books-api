package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-posgres-books-api/models"
	"github.com/golang-posgres-books-api/storage"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/books", r.GetBooks)
	api.Get("/books/:id", r.GetBookByID)
	api.Post("/books", r.CreateBook)
	api.Delete("/books/:id", r.DeleteBook)
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}
	if err := c.BodyParser(&book); err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
			"msg": "request failed",
		})
		return err
	}

	if err := r.DB.Create(&book).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"msg": "could not create book",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"msg": "succesfully created book",
	})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	books := &[]models.Book{}
	if err := r.DB.Find(books).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"msg": "could not fetch books",
		})
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"msg":  "books fetched successfully",
		"data": books,
	})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	book := models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"msg": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(book, id).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"msg": "failed to delete book",
		})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"msg": "book deleted succesfully",
	})
	return nil
}

func (r *Repository) GetBookByID(c *fiber.Ctx) error {
	book := &models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"msg": "id cannot be empty",
		})
	}

	err := r.DB.Where("id = ?", id).First(book).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"msg": "could not find book",
		})
		return nil
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"msg":  "book fetched succesfully",
		"data": book,
	})

	return nil
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("could not load .env file")
		os.Exit(-1)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load database")
		os.Exit(-1)
	}

	if err := models.MigrateBooks(db); err != nil {
		log.Fatal("could not migrate database")
		os.Exit(-1)
	}

	r := &Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8000")

}
