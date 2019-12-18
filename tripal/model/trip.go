package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Trip struct {
	TripId    primitive.ObjectID `json:"tripId" bson:"_id"`
	StartTime *time.Time         `json:"startTime" bson:"startTime"`
	EndTime   *time.Time         `json:"endTime" bson:"endTime"`
	Locations []TripLocation     `json:"locations" bson:"locations"`
}
