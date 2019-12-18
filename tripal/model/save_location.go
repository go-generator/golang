package model

type SaveLocation struct {
	SaveLocationId string `json:"saveLocationId" bson:"_id"`
	UserId         string `json:"userId" bson:"userId"`
	LocationId     string `json:"locationId" bson:"locationId"`
}
