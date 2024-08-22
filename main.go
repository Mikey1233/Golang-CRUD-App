package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	Body      string `json:"Body"`
	Completed bool   `json:"Completed"`
	Id        int    `json:"id" bson:"_id"`
}

var collection *mongo.Collection
func main() {
	fmt.Println("Hello world 4")
	 err := godotenv.Load(".env")
	 if err != nil {
		log.Fatal("Erro fetching .env file")
	 }
	 MONGODB_URI := os.Getenv("MONGODB_URI")
	 clientOptions := options.Client().ApplyURI(MONGODB_URI)
	 client,err := mongo.Connect(context.Background(),clientOptions)

	 if err != nil {
		log.Fatal(err)
	 }
	 err = client.Ping(context.Background(),nil)
	 if err != nil {
		log.Fatal(err)
	 }
	 fmt.Println("Connected to mongoDb")
}