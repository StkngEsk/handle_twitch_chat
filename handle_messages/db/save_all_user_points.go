package db

import (
	"context"
	"log"
	"math"
	"os"
	"time"

	"github.com/StkngEsk/handle_twitch_chat/handle_messages/models"
	"go.mongodb.org/mongo-driver/bson"
)

func SaveAllUserPoints(users []models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := MongoC.Database(os.Getenv("DATABASE_NAME"))
	collectionUsers := db.Collection("users")

	for _, user := range users {

		var userUpdate models.User
		record := make(map[string]interface{})
		condition := bson.M{
			"userId": user.IDUserTwitch,
		}

		_ = collectionUsers.FindOne(ctx, condition).Decode(&userUpdate)

		if len(userUpdate.IDUserTwitch) == 0 {
			_, err := collectionUsers.InsertOne(ctx, user)

			if err != nil {
				log.Fatal("Record could not be saved.")
				return err
			}
		} else {
			if user.Score > 0 {
				record["score"] = math.Floor((userUpdate.Score+user.Score)*100) / 100
			}

			filter := bson.M{
				"_id": bson.M{"$eq": userUpdate.ID},
			}

			updateString := bson.M{
				"$set": record,
			}
			_, err := collectionUsers.UpdateOne(ctx, filter, updateString)

			if err != nil {
				return err
			}
		}

	}
	return nil

}
