package model

type CriterionEvaluation struct {
	CriterionId string  `bson:"criterionId" json:"criterionId"`
	Mark        float64 `bson:"mark" json:"mark"`
}
