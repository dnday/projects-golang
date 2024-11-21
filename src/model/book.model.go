package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Author       string             `bson:"author"`
	Published_at string             `bson:"published_at"`
	Updated_at   time.Time          `bson:"updated_at"`
	Created_at   time.Time          `bson:"created_at"`
}
