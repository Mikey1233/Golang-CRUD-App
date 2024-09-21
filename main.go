package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type Todo struct {
	Body string `json:"body"`
	Id primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool `json:"completed"`
}
var collection *mongo.Collection
func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	//  log.Fatal("error loading env file")
	// }
	if os.Getenv("ENV") != "production"{
	err := godotenv.Load(".env")
			if err != nil {
	 log.Fatal("error loading env file")
	}
	}
	MONGO_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
     client  ,err:= mongo.Connect(context.Background(),clientOptions)
	 if err != nil {
		log.Fatal(err)
	 }
	 defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(),nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mongoDB")
	app := fiber.New()
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	if os.Getenv("ENV") == "production"{
		app.Static("/","./client/dist")
	}
	// connection,err := mongo.c
	collection = client.Database("golang_db").Collection("todos")
	//routes
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins : "http://localhost:5173",
	// 	AllowHeaders : "Origin,Content-Type,Accept",
	// }))
	app.Get("/api/todos", getTodo)
	app.Post("/api/todos",createTodo)
	app.Patch("/api/todos/:id",updateTodo)
	app.Delete("/api/todos/:id",deleteTodo)
	log.Fatal(app.Listen("0.0.0.0:"+PORT))
}
func getTodo(c *fiber.Ctx) error{
	var todos []Todo
  cursor,err :=	collection.Find(context.Background(),bson.M{})
  if err != nil {
	return err
}
  defer cursor.Close(context.Background())
  for cursor.Next(context.Background()) {
	var todo Todo
	err := cursor.Decode(&todo)
	if err != nil {
		return err
	}
	todos = append(todos,todo )
  }

  return c.JSON(todos)
}
func createTodo(ctx *fiber.Ctx) error{
	todo := new(Todo)
	err := ctx.BodyParser(&todo)
	if err != nil {
		return err
	}
	if todo.Body == "" {
		return ctx.Status(404).JSON(fiber.Map{"message":"body cannot be empty"})
	}
	insertedResult,err := collection.InsertOne(context.Background(),todo)
	if err != nil {
		return err
	}
	todo.Id = insertedResult.InsertedID.(primitive.ObjectID)
	return ctx.Status(200).JSON(fiber.Map{"message":"success"})
  }
  func updateTodo(c *fiber.Ctx) error{
	var todo Todo
	id := c.Params("id")
	err := c.BodyParser(&todo)
	if err != nil{
		return err
	}
	ObjectId,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		 	return c.Status(400).JSON(fiber.Map{"error":"Invalid todo ID"})
		   }
		   filter := bson.M{"_id" : ObjectId}
		   update := bson.M{"$set":bson.M{"completed":todo.Completed,"body":todo.Body}}
		   _ , err = collection.UpdateOne(context.Background(),filter,update)
		   if err != nil {
			return err
		   }
	return c.Status(200).JSON(fiber.Map{"success":true})
  }
  func deleteTodo(c *fiber.Ctx) error{
	id := c.Params("id")
	ObjectId,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		 	return c.Status(400).JSON(fiber.Map{"error":"Invalid todo ID"})
		   }
		   filter := bson.M{"_id" : ObjectId}
		  
		   _ , err = collection.DeleteOne(context.Background(),filter)
		   if err != nil {
			return err
		   }
	return c.Status(204).JSON(fiber.Map{"deleted":true})
  }