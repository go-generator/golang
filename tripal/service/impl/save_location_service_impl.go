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

type SaveLocationServiceImpl struct {
	*m.DefaultGenericService
	*m.DefaultSearchService
}

func NewSaveLocationServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *SaveLocationServiceImpl {
	var model model.SaveLocation
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := m.NewMongoGenericSearchService(db, typeOfModel, "saveLocation", searchBuilder, false, "")
	return &SaveLocationServiceImpl{genericService, searchService}
}

func (s *SaveLocationServiceImpl) GetLocationsOfUser(ctx context.Context, userId string) ([]model.Location, error) {
	results := []model.Location{}
	pipeline := []bson.M{
		{"$match": bson.M{"userId": userId}},
		{"$lookup": bson.M{"from": "location", "localField": "locationId", "foreignField": "_id", "as": "ofLocation"}},
		{"$lookup": bson.M{"from": "locationInfo", "localField": "locationId", "foreignField": "_id", "as": "ofLocationInfo"}},
	}
	a, _ := s.Collection.Aggregate(ctx, pipeline)
	if a == nil {
		return nil, nil
	}
	resp := []bson.M{}
	_ = a.All(ctx, &resp)
	fmt.Println("aa", resp)
	for _, v := range resp {
		result := s.mapBsonToLocation(v)
		results = append(results, result)
	}
	return results, nil
}

func (s *SaveLocationServiceImpl) mapBsonToLocation(resp bson.M) model.Location {
	result := model.Location{}
	resultLocationInfo := model.LocationInfo{}
	ofLocations := make(map[string]interface{}, 0)
	aa := make(map[string][]float64, 0)
	// save the node location
	ofLocation := resp["ofLocation"]
	location := ofLocation.(bson.A)[0]
	bsonBytes1, _ := json.Marshal(location)
	json.Unmarshal(bsonBytes1, &ofLocations)
	coordinates := ofLocations["location"]
	bsonBytes3, _ := json.Marshal(coordinates)
	json.Unmarshal(bsonBytes3, &aa)
	locationInfo := resp["ofLocationInfo"].(primitive.A)
	locationInfoSlice := []interface{}(locationInfo)
	result.LocationId = resp["locationId"].(string)
	result.Longitude = aa["coordinates"][0]
	result.Latitude = aa["coordinates"][1]

	bsonBytes, _ := json.Marshal(location)
	json.Unmarshal(bsonBytes, &result)
	bsonBytes2, _ := json.Marshal(locationInfoSlice[0])
	json.Unmarshal(bsonBytes2, &resultLocationInfo)
	result.LocationInfo = resultLocationInfo
	return result
}
