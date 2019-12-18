package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LocationRate struct {
	RateId     primitive.ObjectID `json:"rateId" bson:"_id"`
	LocationId string             `json:"locationId" bson:"locationId"`
	UserId     string             `json:"userId" bson:"userId"`
	Rate       int                `json:"rate" bson:"rate"`
	RateTime   time.Time          `json:"rateTime" bson:"ratetime"`
	Review     string             `json:"review" bson:"review"`
}
