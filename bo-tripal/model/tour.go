package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tour struct {
	TourId    primitive.ObjectID `json:"tourId" bson:"_id"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   time.Time          `json:"endTime" bson:"endTime"`
	Locations []string           `json:"locations" bson:"locations"`
}
