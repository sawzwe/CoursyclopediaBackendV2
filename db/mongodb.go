package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var Database *mongo.Database

func ConnectDB() {
	connectionString := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	DB = client
	Database = client.Database("coursyclopediadb")
}

func GetDB() *mongo.Database {
	return Database
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Database("coursyclopediadb").Collection(collectionName)
}

func DisconnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := DB.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
