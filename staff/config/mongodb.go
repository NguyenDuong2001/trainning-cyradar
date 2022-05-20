package config

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	clientOptions := options.Client().ApplyURI(os.Getenv("URL_MongoDB"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connect to MongoDB")

	}
	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("Team-Staff").Collection(collectionName)
	return collection
}
