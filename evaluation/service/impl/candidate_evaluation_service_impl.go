package impl

import (
	. "../../model"
	SM "../../search-model"
	S "../../service"
	"context"
	"errors"
	m "github.com/common-go/mongo"
	"github.com/common-go/search"
	"github.com/common-go/service"
	"go.mongodb.org/mongo-driver/mongo"
	//"log"
	"reflect"
)

type CandidateEvaluationServiceImpl struct {
	Database     *mongo.Database
	collection   *mongo.Collection
	QueryBuilder m.QueryBuilder
	modelType    reflect.Type
	service.GenericService
	search.SearchService
	canService    S.CandidateService
	schemeService S.SchemeService
}

func NewCandidateEvaluationServiceImpl(db *mongo.Database, searchResultBuilder m.SearchResultBuilder, candidateService S.CandidateService, schemeService S.SchemeService) *CandidateEvaluationServiceImpl {
	var model CandidateEvaluation
	modelType := reflect.TypeOf(model)
	collection := "candidateEvaluation"
	queryBuilder := &m.DefaultQueryBuilder{}
	mongoService, searchService := m.NewMongoGenericSearchService(db, modelType, collection, searchResultBuilder, true, "")
	return &CandidateEvaluationServiceImpl{db, db.Collection(collection), queryBuilder, modelType, mongoService, searchService, candidateService, schemeService}
}

func (s *CandidateEvaluationServiceImpl) Insert(ctx context.Context, obj interface{}) (int64, error) {
	caneva, ok := obj.(*CandidateEvaluation)
	if ok == false {
		return 0, errors.New("invalid data")
	}
	i, err := s.schemeService.GetById(ctx, caneva.SchemeId)
	if err != nil {
		return 0, errors.New("get Scheme failed")
	}
	scheme, ok1 := i.(*Scheme)
	if ok1 == false {
		return 0, errors.New("get Scheme failed")
	}
	caneva.Mark = s.calculateMark(caneva.Evaluations, scheme.Criteria)
	_, err2 := s.reCalculateCandidateMark(ctx, caneva.Candidate, caneva.SchemeId, caneva.BatchId)
	if err2 != nil {
		return 0, errors.New("patch candidate failed")
	}
	_, err3 := s.GenericService.Insert(ctx, caneva)
	if err3 != nil {
		return 1, errors.New("insert failed")
	}
	return 2, nil
}

func (s *CandidateEvaluationServiceImpl) Update(ctx context.Context, obj interface{}) (int64, error) {
	caneva, ok := obj.(*CandidateEvaluation)
	if ok == false {
		return 0, errors.New("invalid data")
	}
	i, err := s.schemeService.GetById(ctx, caneva.SchemeId)
	if err != nil {
		return 0, errors.New("get Scheme failed")
	}
	scheme, ok1 := i.(*Scheme)
	if ok1 == false {
		return 0, errors.New("get Scheme failed")
	}
	caneva.Mark = s.calculateMark(caneva.Evaluations, scheme.Criteria)
	_, err3 := s.GenericService.Update(ctx, caneva)
	if err3 != nil {
		return 0, errors.New("insert failed")
	}
	_, err2 := s.reCalculateCandidateMark(ctx, caneva.Candidate, caneva.SchemeId, caneva.BatchId)
	if err2 != nil {
		return 1, errors.New("patch candidate failed")
	}
	return 2, nil
}

func (s *CandidateEvaluationServiceImpl) calculateMark(evaluations []CriterionEvaluation, criteria []Criterion) float64 {
	var total float64
	for _, evaluation := range evaluations {
		for _, criterion := range criteria {
			if evaluation.CriterionId == criterion.Id {
				total += evaluation.Mark * criterion.Ratio
				break
			}
		}
	}
	return total
}

func (s *CandidateEvaluationServiceImpl) reCalculateCandidateMark(ctx context.Context, candidateId string, schemeId string, batchId string) (float64, error) {
	candidateEvaluations, err := s.findEvaluation(ctx, candidateId, schemeId, batchId)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, candidateEvaluation := range candidateEvaluations {
		total += candidateEvaluation.Mark
	}
	err = s.canService.PatchMark(ctx, candidateId, total/float64(len(candidateEvaluations)))
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *CandidateEvaluationServiceImpl) findEvaluation(ctx context.Context, candidate string, schemeId string, batchId string) ([]CandidateEvaluation, error) {
	sm := SM.CandidateEvaluationSM{
		BatchId:   batchId,
		SchemeId:  schemeId,
		Candidate: candidate,
	}
	query := s.QueryBuilder.BuildQuery(sm, s.modelType)
	caneva, err := m.Find(ctx, s.collection, query, s.modelType)
	if err != nil {
		return nil, err
	}
	canevas := caneva.(*[]CandidateEvaluation)
	return *canevas, nil
}
