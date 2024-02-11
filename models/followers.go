package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResponseFollowerTweets struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	UserId         string             `bson:"userid" json:"userid,omitempty"`
	UserRelationId string             `bson:"userrelationid" json:"userrelationid,omitempty"`
	Tweet          struct {           // this is a new definition, not the same as models.Tweet struct
		Message string    `bson:"mensaje" json:"mensaje,omitempty"`
		Date    time.Time `bson:"fecha" json:"fecha,omitempty"`
		Id      string    `bson:"_id" json:"_id,omitempty"`
	}
}
