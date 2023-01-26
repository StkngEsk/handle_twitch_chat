package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoC = ConnectDB()

func ConnectDB() *mongo.Client {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	clientOptions := options.Client().ApplyURI(os.Getenv("DATABASE_URL"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	log.Println("Connected db success")
	return client
}

func CheckConnection() bool {
	err := MongoC.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}
	return true
}
