package impl

import (
  . "../../model"
  m "github.com/common-go/mongo"
  . "github.com/common-go/search"
  . "github.com/common-go/service"
  "go.mongodb.org/mongo-driver/mongo"
  "reflect"
)

type SchemeServiceImpl struct {
  database   *mongo.Database
  collection *mongo.Collection
  GenericService
  SearchService
}

func NewSchemeServiceImpl(db *mongo.Database, searchResultBuilder m.SearchResultBuilder) *SchemeServiceImpl {
  var model Scheme
  modelType := reflect.TypeOf(model)
  collection := "scheme"
  mongoService, searchService := m.NewMongoGenericSearchService(db, modelType, collection, searchResultBuilder, false, "")
  return &SchemeServiceImpl{db, db.Collection(collection), mongoService, searchService}
}
