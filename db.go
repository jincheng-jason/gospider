package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Item struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	ItemId string `json:"itemId" bson:"itemId"`
	Title string `json:"title" bson:"title"`
	Price float64 `json:"price" bson:"price"`
	ImgUrl string `json:"imgUrl" bson:"imgUrl"`
	BrandImgUrl string `json:"brandImgUrl" bson:"brandImgUrl"`
}

func InitMongo(ctx context.Context, url string, timeout int) (*mongo.Collection, error)  {
	//mongodb client init
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// defer client.Disconnect(ctx)

	collection := client.Database("price").Collection("items")
	return collection, nil
}