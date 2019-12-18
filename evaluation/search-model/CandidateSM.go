package search_model

import "github.com/common-go/search"

type CandidateSM struct {
	*search.SearchModel
	Email string `json:"email"`
}
