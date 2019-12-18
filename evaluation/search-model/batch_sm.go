package search_model

import "github.com/common-go/search"

type BatchSM struct {
	*search.SearchModel
	BatchId   string `json:"batchId"`
	BatchName string `json:"batchName"`
}
