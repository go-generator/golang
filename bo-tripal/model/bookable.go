package model

type BookableType string

const (
	Room      BookableType = "R"
	Projector BookableType = "P"
)

type Bookable struct {
	BookableId          string       `json:"bookableId" bson:"_id"`
	LocationId          string       `json:"locationId" bson:"locationId"`
	BookableType        BookableType `json:"bookableType" bson:"bookableType"`
	BookableName        string       `json:"bookableName" bson:"bookableName"`
	BookableDescription string       `json:"bookableDescription" bson:"bookableDescription"`
	BookableCapacity    int          `json:"bookableCapacity" bson:"bookableCapacity"`
	Image               string       `json:"image" bson:"image"`
}
