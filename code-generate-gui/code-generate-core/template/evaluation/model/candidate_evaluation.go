package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CandidateEvaluation struct {
	Id          primitive.ObjectID    `bson:"_id" json:"_id"`
	BatchId     string                `bson:"batchId" json:"batchId"`
	Evaluator   string                `bson:"evaluator" json:"evaluator"`
	Candidate   string                `bson:"candidate" json:"candidate"`
	SchemeId    string                `bson:"schemeId" json:"schemeId"`
	Evaluations []CriterionEvaluation `bson:"evaluations" json:"evaluations"`
	Status      Status                `bson:"status" json:"status"`
	Mark        float64               `bson:"mark" json:"mark"` // auto calculate
	//note  string `bson:"mark" json:"mark"`
}
