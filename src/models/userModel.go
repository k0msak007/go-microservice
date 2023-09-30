package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ObjectId primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Username string             `bson:"username" json:"username"`
	Item     []Item             `json:"item"`
}

type UserItem struct {
	UserId string `bson:"user_id" json:"user_id"`
	ItemId string `bson:"item_id" json:"item_id"`
}
