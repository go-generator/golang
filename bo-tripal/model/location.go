package model

type Location struct {
	LocationId   string       `json:"locationId" bson:"_id"`
	LocationName string       `json:"locationName" bson:"locationName"`
	LocationInfo LocationInfo `json:"locationInfo" bson:"locationInfo"`
	Description  string       `json:"description" bson:"description"`
	Type         string       `json:"type" bson:"type"`
	Longitude    float64      `json:"longitude" bson:"longitude"`
	Latitude     float64      `json:"latitude" bson:"latitude"`
}

func (Location) CollectionName() string {
	return "location"
}
