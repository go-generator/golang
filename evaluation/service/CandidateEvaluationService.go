package service

import(
	search "github.com/common-go/search"
	service "github.com/common-go/service"
)

type CandidateEvaluationService interface{
	service.GenericService
	search.SearchService
}
