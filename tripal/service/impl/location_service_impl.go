package impl

import (
	"../../model"
	"../../service"
	"context"
	"encoding/json"
	"fmt"
	m "github.com/common-go/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strings"
)

type LocationServiceImpl struct {
	*m.DefaultGenericService
	*m.DefaultSearchService
	SaveLocationService service.SaveLocationService
}

func NewLocationServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder, saveLocationService service.SaveLocationService) *LocationServiceImpl {
	var model model.Location
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := m.NewMongoGenericSearchService(db, typeOfModel, model.CollectionName(), searchBuilder, false, "")
	return &LocationServiceImpl{genericService, searchService, saveLocationService}
}
func (p *LocationServiceImpl) GetAll(ctx context.Context) (interface{}, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{"from": "locationInfo", "localField": "_id", "foreignField": "_id", "as": "locationInfo"}},
	}
	a, _ := p.Collection.Aggregate(ctx, pipeline)
	if a == nil {
		return nil, nil
	}
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
	if a == nil {
		return nil, nil
	}
	resp := []bson.M{}
	_ = a.All(ctx, &resp)
	result := model.Location{}
	result = p.mapBsonToLocation(resp[0])
	return result, nil
}

func (p *LocationServiceImpl) GetByUrlId(ctx context.Context, urlId string) (*model.Location, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"urlId": urlId}},
		{"$lookup": bson.M{"from": "locationInfo", "localField": "_id", "foreignField": "_id", "as": "locationInfo"}},
	}
	a, _ := p.Collection.Aggregate(ctx, pipeline)
	if a == nil {
		return nil, nil
	}
	resp := []bson.M{}
	_ = a.All(ctx, &resp)
	result := model.Location{}
	result = p.mapBsonToLocation(resp[0])
	return &result, nil
}

func (p *LocationServiceImpl) RateLocation(ctx context.Context, objRate model.LocationRate) (bool, error) {
	result, err := m.InsertOne(ctx, p.Database.Collection("locationRate"), objRate)
	if result == 1 {
		return true, nil
	}
	return false, err
}

func (p *LocationServiceImpl) SaveLocation(ctx context.Context, userId string, locationId string) (bool, error) {
	id := uuid.New()
	saveLocationId := strings.Replace(id.String(), "-", "", -1)
	obj := model.SaveLocation{saveLocationId, userId, locationId}

	result, err := p.SaveLocationService.Insert(ctx, obj)
	if result == 1 {
		return true, nil
	}
	return false, err
}

func (p *LocationServiceImpl) RemoveLocation(ctx context.Context, userId string, locationId string) (bool, error) {
	result := []model.SaveLocation{}
	query := bson.M{
		"userId":     userId,
		"locationId": locationId,
	}
	v, er0 := m.FindAndDecode(ctx, p.Database.Collection("saveLocation"), query, &result)
	if v {
		query1 := bson.M{
			"_id": result[0].SaveLocationId,
		}
		n, er1 := m.DeleteOne(ctx, p.Database.Collection("saveLocation"), query1)
		if er1 != nil {
			return false, er1
		}
		return true, nil
		fmt.Println("result", n, er1)
	}
	return false, er0
}

func (p *LocationServiceImpl) GetLocationsOfUser(ctx context.Context, userId string) ([]model.Location, error) {
	return p.SaveLocationService.GetLocationsOfUser(ctx, userId)
}

func (p *LocationServiceImpl) GetLocationByTypeInRadius(ctx context.Context, type1 string, raidus int) ([]model.Location, error) {
	results := []model.Location{}
	pipeline := []bson.M{
		//{"$geoNear": bson.M{
		//	"near": bson.M{
		//		"type": "Point", "coordinates": []float64{106.624352931976, 10.8528483653576},
		//	},
		//	"key": "location",
		//	"distanceField": "distance",
		//	"maxDistance": raidus,
		//	"spherical": true},
		//},
		{"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{106.624352931976, 10.8528483653576},
				},
				"$maxDistance": raidus,
			},
		}},
		{"$match": bson.M{"type": strings.ToLower(type1)}},
		{"$lookup": bson.M{"from": "locationInfo", "localField": "_id", "foreignField": "_id", "as": "locationInfo"}},
	}
	a, _ := p.Collection.Aggregate(ctx, pipeline)
	if a == nil {
		return nil, nil
	}
	resp := []bson.M{}
	_ = a.All(ctx, &resp)
	for _, v := range resp {
		result := p.mapBsonToLocation(v)
		results = append(results, result)
	}
	fmt.Println("result", results)
	return results, nil

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
