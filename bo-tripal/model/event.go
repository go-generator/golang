package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Event struct {
	EventId    primitive.ObjectID `json:"eventId" bson:"_id"`
	EventName  string             `json:"eventName" bson:"eventName"`
	StartTime  *time.Time         `json:"startTime" bson:"startTime"`
	EndTime    *time.Time         `json:"endTime" bson:"endTime"`
	LocationId string             `json:"locationId" bson:"locationId"`
	Lat        float64            `json:"lat" bson:"lat"`
	Long       float64            `json:"long" bson:"long"`
}
