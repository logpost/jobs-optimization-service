package mongo

import (
	"log"
	"fmt"
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/logpost/jobs-optimization-service/config"
)

type MongoClient struct {
	session	*mongo.Client
	config	config.DatabaseConfig
}

func Connection(config config.DatabaseConfig) MongoClient {

	var mongoClient MongoClient

	mongoClient.config	=	config
	mongoOptions		:=	options.Client().ApplyURI(config.DatabaseURI)
	client, err			:=	mongo.Connect(context.TODO(), mongoOptions)
	
	if err != nil {
		panic(err)
	}

	err	=	client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	mongoClient.session	= client
	fmt.Println("Connected to MongoDB ☄️")

	return mongoClient

}

func changeStream(routineCtx context.Context, waitGroup sync.WaitGroup, stream *mongo.ChangeStream) {

	defer stream.Close(routineCtx)
	defer waitGroup.Done()

	for stream.Next(routineCtx) {
		var data bson.M
		
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", data)
	}

}

func (client MongoClient) WatchCollection(collName string) { 

	coll			:=	client.session.Database(client.config.DatabaseName).Collection(collName)
	cursor, err		:=	coll.Watch(context.TODO(), mongo.Pipeline{})
	
	if err != nil {
		panic(err)
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	routineCtx, _	:= context.WithCancel(context.Background())
	go changeStream(routineCtx, waitGroup, cursor)

	waitGroup.Wait()

}