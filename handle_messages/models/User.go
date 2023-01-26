package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name                string             `bson:"username" json:"username,omitempty"`
	Score               float64            `bson:"score" json:"score,omitempty"`
	ImageUrl            string             `bson:"image_url" json:"image_url,omitempty"`
	IDUserTwitch        string             `bson:"userId" json:"userId,omitempty"`
	ColorPrimaryGlobe   string             `bson:"cpGlobe" json:"cpGlobe,omitempty"`
	ColorSecondaryGlobe string             `bson:"csGlobe" json:"csGlobe,omitempty"`
	OpacityGlobe        float64            `bson:"opacityGlobe" json:"opacityGlobe,omitempty"`
}
