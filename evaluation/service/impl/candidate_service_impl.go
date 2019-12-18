package impl

import (
	. "../../model"
	"context"
	m "github.com/common-go/mongo"
	"github.com/common-go/search"
	"github.com/common-go/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type CandidateServiceImpl struct {
	database   *mongo.Database
	collection *mongo.Collection
	service.GenericService
	search.SearchService
}

func NewCandidateServiceImpl(db *mongo.Database, searchResultBuilder m.SearchResultBuilder) *CandidateServiceImpl {
	var model Candidate
	modelType := reflect.TypeOf(model)
	collection := "candidate"
	mongoService, searchService := m.NewMongoGenericSearchService(db, modelType, collection, searchResultBuilder, false, "")
	return &CandidateServiceImpl{db, db.Collection(collection), mongoService, searchService}
}

func (s *CandidateServiceImpl) ImportArrayObject(ctx context.Context, arr []Candidate) (int64, error) {
	_, err := m.UpsertMany(ctx, s.collection, arr)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *CandidateServiceImpl) PatchMark(ctx context.Context, _id string, mark float64) error {
	var status string
	if mark >= 7 {
		status = "PA"
	} else {
		status = "FA"
	}
	object := Candidate{
		Id:     _id,
		Status: status,
		Mark:   mark,
	}
	body := s.NewModelMap(object)
	query := bson.M{"_id": _id}
	_, err := m.PatchOne(ctx, s.collection, body, query)
	if err != nil {
		return err
	}
	return nil
}

func (c *CandidateServiceImpl) NewModelMap(body interface{}) (out map[string]interface{}) {
	queryModel := make(map[string]interface{})
	value := reflect.Indirect(reflect.ValueOf(body))
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		key1 := value.Type().Field(i).Name
		if key1 != "ControlModel" && (((value.Field(i).Kind().String() == "string" || value.Field(i).Kind().String() == "array") && value.Field(i).Len() > 0) || (value.Field(i).Kind().String() == "float64" && value.Field(i).Float() != 0)) || (value.Field(i).Kind().String() == "ptr" && value.Field(i).Pointer() != 0) {
			val := value.Field(i).Interface()
			queryModel[key1] = val
		}
	}
	return queryModel
}
