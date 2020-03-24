package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"../../model"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationServiceImpl struct {
	*m.DefaultGenericService
	*m.DefaultSearchService
}

func NewLocationServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *LocationServiceImpl {
	var model model.Location
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := m.NewMongoGenericSearchService(db, typeOfModel, model.CollectionName(), searchBuilder, false, "")
	return &LocationServiceImpl{genericService, searchService}
}
func (p *LocationServiceImpl) GetAll(ctx context.Context) (interface{}, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{"from": "locationInfo", "localField": "_id", "foreignField": "_id", "as": "locationInfo"}},
	}
	a, _ := p.Collection.Aggregate(ctx, pipeline)
	resp := []bson.M{}
	err := a.All(ctx, &resp)
	if err != nil {
		fmt.Println("ERROR", err)
		return nil, err
	}
	results := []model.Location{}
	for i := 0; i < len(resp); i++ {
		result := model.Location{}
		result = p.mapBsonToLocation(resp[i])
		results = append(results, result)
	}
	return results, nil
}

func (p *LocationServiceImpl) GetById(ctx context.Context, id interface{}) (interface{}, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": id.(string)}},
		{"$lookup": bson.M{"from": "locationInfo", "localField": "_id", "foreignField": "_id", "as": "locationInfo"}},
	}
	a, _ := p.Collection.Aggregate(ctx, pipeline)
	resp := []bson.M{}
	_ = a.All(ctx, &resp)
	result := model.Location{}
	result = p.mapBsonToLocation(resp[0])
	return result, nil
}

func (p *LocationServiceImpl) mapBsonToLocation(resp bson.M) model.Location {
	result := model.Location{}
	resultLocationInfo := model.LocationInfo{}
	// save the node location
	location := resp["location"]
	coordinates := location.(bson.M)["coordinates"].(primitive.A)
	valueMSI := []interface{}(coordinates)
	locationInfo := resp["locationInfo"].(primitive.A)
	locationInfoSlice := []interface{}(locationInfo)
	result.LocationId = resp["_id"].(string)
	resp["longitude"] = valueMSI[0]
	resp["latitude"] = valueMSI[1]
	delete(resp, "location")

	bsonBytes, _ := json.Marshal(resp)
	json.Unmarshal(bsonBytes, &result)
	bsonBytes1, _ := json.Marshal(locationInfoSlice[0])
	json.Unmarshal(bsonBytes1, &resultLocationInfo)
	result.LocationInfo = resultLocationInfo
	return result
}
