package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2" 
	"github.com/joho/godotenv"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()
     err := godotenv.Load(".env")
	 if err != nil {
		log.Fatal(("Error loading .env file"))
	 }
	 PORT := os.Getenv("PORT")
	todos := []Todo{}
	//get route
	app.Get("/api/todos", func(c *fiber.Ctx) error {
  		return c.Status(200).JSON(todos)
	})
	//create or post route
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		err := c.BodyParser(todo)
		if err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"message": "Todo body is required"})
		}
		todo.Id = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(200).JSON(todo)
	})
	//update route
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos)
			}
		}
		return c.Status(400).JSON(fiber.Map{"error": "todo not found "})
	})
    //delete todo
	app.Delete("api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(400).JSON(fiber.Map{"succes": "true"})

			}
		}  
		return c.SendStatus(404)
	})

	log.Fatal(app.Listen(":"+PORT))

}
