package search_model

import "github.com/common-go/search"

type CandidateEvaluationSM struct {
	*search.SearchModel
	BatchId string `json:"batchId" bson:"batchId"`
	//Evaluator string `json:"evaluator" bson:"evaluator"`
	SchemeId  string `json:"schemeId" bson:"schemeId"`
	Candidate string `json:"candidate" bson:"candidate"`
}
