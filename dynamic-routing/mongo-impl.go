package main

import (
	"context"

	"github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func (t MongoType) GetById() error {
	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, t.Info.DataSourceName, t.Info.DbName)
	if err != nil {
		return err
	}
	collection := db.Collection(t.Info.Source)
	result := collection.FindOne(ctx, bson.M{"_id": t.IdParam})
	tmp, err := result.DecodeBytes()
	*t.OutputString = tmp.String()
	return err
}
func (t MongoType) GetAll() error {
	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, t.Info.DataSourceName, t.Info.DbName)
	if err != nil {
		return err
	}
	collection := db.Collection(t.Info.Source)
	result, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return err
	}
	var tmp bson.Raw
	for result.Next(ctx) {
		err = result.Decode(&tmp)
		*t.OutputString += tmp.String()
	}
	return err
}
func (t MongoType) Create() error {
	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, t.Info.DataSourceName, t.Info.DbName)
	if err != nil {
		return err
	}
	collection := db.Collection(t.Info.Source)
	_, err = collection.InsertOne(ctx, t.InputMap)
	if err != nil {
		return err
	}
	*t.OutputString = "Created Successfully"
	return err
}
func (t MongoType) Update() error {
	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, t.Info.DataSourceName, t.Info.DbName)
	if err != nil {
		return err
	}
	collection := db.Collection(t.Info.Source)
	_, err = collection.ReplaceOne(ctx, bson.D{{"_id", t.IdParam}}, t.InputMap)
	if err != nil {
		return err
	}
	*t.OutputString = "Updated Successfully"
	return err
}
func (t MongoType) Delete() error {
	ctx := context.Background()
	db, err := mongo.CreateConnection(ctx, t.Info.DataSourceName, t.Info.DbName)
	if err != nil {
		return err
	}
	collection := db.Collection(t.Info.Source)
	_, err = collection.DeleteOne(ctx, bson.D{{"_id", t.IdParam}})
	if err != nil {
		return err
	}
	*t.OutputString = "Deleted Successfully"
	return err
}
