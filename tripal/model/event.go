package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	EventId    primitive.ObjectID `json:"eventId,omitempty" bson:"_id, omitempty" gorm:"type:varchar(500);column:_id;primary_key"`
	EventName  string             `json:"eventName" bson:"eventName, omitempty"`
	StartTime  *time.Time         `json:"startTime" bson:"startTime"`
	EndTime    *time.Time         `json:"endTime" bson:"endTime"`
	LocationId string             `json:"locationId" bson:"locationId"`
	Lat        float64            `json:"lat" bson:"lat"`
	Long       float64            `json:"long" bson:"long"`
}
