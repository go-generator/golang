package model

type TripLocation struct {
	LocationId string `json:"locationId" bson:"locationId"`
	Visited    bool   `json:"visited" bson:"visited"`
}
