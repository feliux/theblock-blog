package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tweet struct {
	Message string `bson:"mensaje" json:"mensaje"`
}

type SaveTweet struct {
	UserId  string    `bson:"userid" json:"userid,omitempty"`
	Message string    `bson:"mensaje" json:"mensaje,omitempty"`
	Date    time.Time `bson:"fecha" json:"fecha,omitempty"`
}

type RetrieveTweets struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserId  string             `bson:"userid" json:"userid,omitempty"`
	Message string             `bson:"mensaje" json:"mensaje,omitempty"`
	Date    time.Time          `bson:"fecha" json:"fecha,omitempty"`
}
