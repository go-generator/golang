package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"../model"
	. "github.com/common-go/mongo"
	. "github.com/common-go/search"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LocationSearchResultBuilder struct {
	Context      context.Context
	Database     *mongo.Database
	QueryBuilder QueryBuilder
	SortBuilder  SortBuilder
}

func (b *LocationSearchResultBuilder) BuildSearchResult(ctx context.Context, collection *mongo.Collection, m interface{}, modelType reflect.Type) (*SearchResult, error) {
	query := b.QueryBuilder.BuildQuery(m, modelType)

	var sort = bson.M{}
	var searchModel *SearchModel

	if sModel, ok := m.(*SearchModel); ok {
		searchModel = sModel
		sort = b.SortBuilder.BuildSort(*sModel, modelType)
	} else {
		value := reflect.Indirect(reflect.ValueOf(m))
		numField := value.NumField()
		for i := 0; i < numField; i++ {
			if sModel1, ok := value.Field(i).Interface().(*SearchModel); ok {
				searchModel = sModel1
				sort = b.SortBuilder.BuildSort(*sModel1, modelType)
			}
		}
	}
	return b.Build(ctx, collection, modelType, query, sort, int(searchModel.PageIndex), int(searchModel.PageSize))
}

func (b *LocationSearchResultBuilder) Build(ctx context.Context, collection *mongo.Collection, modelType reflect.Type, query bson.M, sort bson.M, pageIndex int, pageSize int) (*SearchResult, error) {
	pipeline := []bson.M{
		{"$match": query},
		{"$limit": int64(pageSize)},
		{"$skip": int64(pageSize * (pageIndex - 1))},
		{"$lookup": bson.M{"from": "locationInfo", "localField": "_id", "foreignField": "_id", "as": "locationInfo"}},
	}

	if len(sort) != 0 {
		pipeline = append(pipeline, bson.M{"$sort": sort})
	}

	databaseQuery, errFind := collection.Aggregate(ctx, pipeline)
	if errFind != nil {
		return nil, errFind
	}

	resp := []bson.M{}
	errAll := databaseQuery.All(ctx, &resp)
	if errAll != nil {
		fmt.Println(errAll)
	}

	results := []model.Location{}
	for i := 0; i < len(resp); i++ {
		result := model.Location{}
		result = b.mapBsonToLocation(resp[i])
		results = append(results, result)
	}

	var count int
	options := options.Count()
	countDB, errCount := collection.CountDocuments(ctx, query, options)
	if errCount != nil {
		count = 0
	}
	count = int(countDB)
	fmt.Println("count", count)

	searchResult := SearchResult{}
	searchResult.ItemTotal = int64(count)

	searchResult.LastPage = false
	lengthModels := reflect.Indirect(reflect.ValueOf(results)).Len()
	if pageSize*pageIndex+lengthModels >= count {
		searchResult.LastPage = true
	}

	searchResult.Results = results

	return &searchResult, nil
}

func (p *LocationSearchResultBuilder) mapBsonToLocation(resp bson.M) model.Location {
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
