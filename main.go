package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Books struct {
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
