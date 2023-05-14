package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
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

func main() {
	app := fiber.New()

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load database")
	}

	r := &Repository{
		DB: db,
	}

	r.SetupRoutes(app)
	app.Listen(":8000")

}
