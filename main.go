package main

import (
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	shadow *mongo.Collection
)

func init() {
	client := connect()
	shadow = databaseCollection(client, "Shadow", "shadow")
}

func main() {
	println("Main Function")
}
