package model

type Scheme struct {
	SchemeId   string      `bson:"_id" json:"schemeId"`
	SchemeName string      `bson:"schemeName" json:"schemeName"` // promote SE to C
	Criteria   []Criterion `bson:"criteria" json:"criteria"`
}
