package service

import (
  . "../model"
  "context"
  search "github.com/common-go/search"
  service "github.com/common-go/service"
)

type CandidateService interface {
  service.GenericService
  search.SearchService
  ImportArrayObject(ctx context.Context, arr []Candidate) (int64, error)
  PatchMark(ctx context.Context, _id string, mark float64) error
}
