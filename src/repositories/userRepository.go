package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/k0msak007/go-microservice/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Client *mongo.Client
}

func (r *UserRepository) FindOneUser(ctx context.Context, userId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 18*time.Second)
	defer cancel()

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	pipeline := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "users_items"},
					{"localField", "_id"},
					{"foreignField", "user_id"},
					{"as", "items"},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"_id", 1},
					{"username", 1},
					{"items", "$items.item_id"},
				},
			},
		},
		bson.D{{"$match", bson.D{{"_id", userObjectId}}}},
	}

	cursor, err := r.Client.Database("user_db").Collection("users").Aggregate(ctx, pipeline, nil)
	if err != nil {
		return nil, fmt.Errorf("aggregate users failed: %v", err)
	}
	defer cursor.Close(ctx)

	type result struct {
		Id       primitive.ObjectID   `bson:"_id"`
		Username string               `bson:"username"`
		Items    []primitive.ObjectID `bson:"items"`
	}
	temp := result{
		Items: make([]primitive.ObjectID, 0),
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&temp); err != nil {
			return nil, fmt.Errorf("decode user failed: %v", err)
		}
	}

	return &models.User{
		ObjectId: temp.Id,
		Username: temp.Username,
		Item: func() []models.Item {
			items := make([]models.Item, 0)
			for _, id := range temp.Items {
				items = append(items, models.Item{
					ObjectId: id,
				})
			}
			return items
		}(),
	}, nil
}
