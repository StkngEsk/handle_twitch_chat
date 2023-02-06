package db

import (
	"context"
	"os"
	"time"

	"github.com/StkngEsk/handle_twitch_chat/handle_messages/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers(userId string) (user models.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoC.Database(os.Getenv("DATABASE_NAME"))
	collectionUsers := db.Collection("users")

	var userData models.User
	condition := bson.M{
		"userId": userId,
	}

	_ = collectionUsers.FindOne(ctx, condition).Decode(&userData)

	return userData
}
