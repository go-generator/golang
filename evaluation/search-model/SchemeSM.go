package search_model

import "github.com/common-go/search"

type SchemeSM struct {
	*search.SearchModel
	SchemeId   string `json:"schemeId"`
	SchemeName string `json:"schemeName"`
}
