package impl

import (
	. "../../model"
	m "github.com/common-go/mongo"
	. "github.com/common-go/search"
	. "github.com/common-go/service"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type BatchServiceImpl struct {
	database   *mongo.Database
	collection *mongo.Collection
	GenericService
	SearchService
}

func NewBatchServiceImpl(db *mongo.Database, searchResultBuilder m.SearchResultBuilder) *BatchServiceImpl {
	var model Batch
	modelType := reflect.TypeOf(model)
	collection := "batch"
	mongoService, searchService := m.NewMongoGenericSearchService(db, modelType, collection, searchResultBuilder, false, "")
	return &BatchServiceImpl{db, db.Collection(collection), mongoService, searchService}
}
