package models

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claim struct {
	Email string             `json:"email"`
	Id    primitive.ObjectID `bson:"_id" json:"_id,moitempty"`
	jwt.RegisteredClaims
}
