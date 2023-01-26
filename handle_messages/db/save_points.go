package db

import (
	"context"
	"os"
	"time"

	"github.com/StkngEsk/handle_twitch_chat/handle_messages/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SavePoints(user models.User) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoC.Database(os.Getenv("DATABASE_NAME"))
	users := db.Collection("users")

	result, err := users.InsertOne(ctx, user)

	if err != nil {
		return "", false, err
	}

	ObjectID, _ := result.InsertedID.(primitive.ObjectID)
	return ObjectID.String(), true, nil

}
