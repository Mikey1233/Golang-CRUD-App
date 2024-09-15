package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection


type Todo struct {
	Body  string  `json:"body"`
	Completed bool  `json:"completed"`
	Id  primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
}
func main(){
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error fetching env variables")
	}
	//creates an app instance
	app := fiber.New()

	//fetches the mongo_uri string from the env file
	MONGO_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGO_URI)
	//create a mongodb client
	client,err := mongo.Connect(context.Background(),clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//closes the connection once the main func is executed
	defer client.Disconnect(context.Background())
	// sends a signal  thea the client can connect to DB/returns an error if client can't 
	err = client.Ping(context.Background(),nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mongodb database")

	//gets the collection from the mongo DB
	collection = client.Database("golang_db").Collection("todos")

	app.Get("/api/todos",getTodos)
	app.Post("/api/todos",createTodo)
	app.Patch("/api/todos/:id",updateTodo)
	app.Delete("/api/todos/:id",deleteTodo)


	//creates a PORT 5000
	PORT := os.Getenv("PORT")
	if PORT == ""{
		PORT = "5000"
	}
	
	//listens to the PORT 5000, logs 
	log.Fatal(app.Listen("0.0.0.0:" + PORT)) //error and exits the main if any error occurs
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	//when you execute a query in mongodb it returns a cursor and an error
	cursor,err := collection.Find(context.Background(),bson.M{})
	if err != nil {
		return err
	}
	//closes the function once the surrounding func ends (getTodos) 
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()){
		var todo Todo
	 err := cursor.Decode(&todo)
	 if err != nil {
		return err
	 }
	 todos = append(todos, todo)
	}
	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx)error{
	todo := new(Todo)
   
   err := c.BodyParser(todo)
   if err != nil {
	return err
   }

   if todo.Body ==  ""{
	return c.Status(400).JSON(fiber.Map{"error":"Todo body cannot be empty"})
   }
  insertResult ,err := collection.InsertOne(context.Background(),todo)
  if err != nil {
	return err
  }
  todo.Id = insertResult.InsertedID.(primitive.ObjectID)
  return c.Status(201).JSON(todo)
}
func updateTodo(c *fiber.Ctx)error{
  //collect id of the doc
  var todo Todo
  id := c.Params("id")
  err := c.BodyParser(&todo)
  if err != nil {
	return err
  }
  
  
  //converts string id to OjectId
  ObjectID,err := primitive.ObjectIDFromHex(id)
  if err != nil {
	return c.Status(400).JSON(fiber.Map{"error":"Invalid todo ID"})
  }

  filter := bson.M{"_id": ObjectID}
  update := bson.M{"$set":bson.M{"completed":true,"body":todo.Body}}
  //edit the doc collected
 _ , err = collection.UpdateOne(context.Background(),filter,update)

 if err!=nil{
	return err
 }

 return c.Status(200).JSON(fiber.Map{"success":true})
  
}
func deleteTodo(c *fiber.Ctx)error{
  id := c.Params("id")
  ObjectID,err := primitive.ObjectIDFromHex(id)
  if err != nil {
	return c.Status(400).JSON(fiber.Map{"error":"Invalid todo ID"})
  }
filter := bson.M{"_id":ObjectID}
_,err =  collection.DeleteOne(context.Background(),filter)
if err != nil {
	return err
}
return c.Status(200).JSON(fiber.Map{"deleted": "success"})
}
