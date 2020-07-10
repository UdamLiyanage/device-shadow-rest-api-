package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Operations interface {
		Create(model interface{}) (interface{}, error)
		Read(filter bson.M) (*mongo.Cursor, error)
		Index(filter bson.M) *mongo.SingleResult
		Update(urn string, update interface{}) (*mongo.UpdateResult, error)
		Delete(urn string) (bool, error)
	}

	Configuration struct {
		Collection *mongo.Collection
	}
)

func (config Configuration) Create(model interface{}) (interface{}, error) {
	doc, err := config.Collection.InsertOne(context.Background(), model)
	return doc, err
}

func (config Configuration) Read(filter bson.M) (*mongo.Cursor, error) {
	docs, err := config.Collection.Find(context.Background(), filter)
	return docs, err
}

func (config Configuration) Index(filter bson.M) *mongo.SingleResult {
	res := config.Collection.FindOne(context.Background(), filter)
	return res
}

func (config Configuration) Update(urn string, update interface{}) (*mongo.UpdateResult, error) {
	filter := bson.M{"urn": urn}
	doc, err := config.Collection.UpdateOne(context.Background(), filter, update)
	return doc, err
}

func (config Configuration) Delete(urn string) (bool, error) {
	filter := bson.M{"urn": urn}
	delRes, err := config.Collection.DeleteOne(context.Background(), filter)
	if err != nil || delRes.DeletedCount == 0 {
		err = errors.New("delete error: delete count less than 1")
		return false, err
	}
	return true, nil
}
