package model

//import "time"

type Candidate struct {
	Id              string      `bson:"_id" json:"id"`
	BatchId         string      `bson:"batchId" json:"batchId"`
	BadgeID         string      `bson:"badgeID" json:"badgeID"`
	FullName        string      `bson:"fullName" json:"fullName"`
	Email           string      `bson:"email" json:"email"`
	Phone           string      `bson:"phone" json:"phone"`
	Project         string      `bson:"project" json:"project"`
	Role            string      `bson:"role" json:"role"`
	Dc              string      `bson:"dc" json:"dc"`
	CurrentPosition string      `bson:"currentPosition" json:"currentPosition"`
	PromoteTo       string      `bson:"promoteTo" json:"promoteTo"`
	Reviewer        string      `bson:"reviewer" json:"reviewer"`
	Topic           string      `bson:"topic" json:"topic"`
	Outlines        string      `bson:"outlines" json:"outlines"`
	Highlight       string      `bson:"highlight" json:"highlight"`
	Priority        string      `bson:"priority" json:"priority"`
	Status          string      `bson:"status" json:"status"`
	PresentDate     interface{} `bson:"presentDate" json:"presentDate"`
	Mark            float64     `bson:"mark" json:"mark"`
	//JudgeHead string `bson:"badgeID" json:"badgeID"`
}
