package model

type Criterion struct {
	Id          string  `bson:"id" json:"id"`
	Description string  `bson:"description" json:"description"`
	Ratio       float64 `bson:"ratio" json:"ratio"`
}
