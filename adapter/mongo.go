package adapter

import (
	"log"
	"fmt"
	"context"

	"github.com/logpost/jobs-optimization-service/config"
	"github.com/logpost/jobs-optimization-service/models"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	COLLECTION_JOBS	=	"jobs"
)

type MongoClient struct {
	session	*mongo.Client
	config	config.DatabaseConfig
}

func CreateMongoConnection(config config.DatabaseConfig) MongoClient {

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

func changeStream(routineCtx context.Context, stream *mongo.ChangeStream) {

	defer stream.Close(routineCtx)
	// defer waitGroup.Done()

	for stream.Next(routineCtx) {
		var data bson.M
		
		if err := stream.Decode(&data); err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", data)
	}

}

func (client MongoClient) WatchJobs() { 

	coll			:=	client.session.Database(client.config.DatabaseName).Collection(COLLECTION_JOBS)
	cursor, err		:=	coll.Watch(context.TODO(), mongo.Pipeline{})
	
	if err != nil {
		panic(err)
	} 

	routineCtx, _	:= context.WithCancel(context.Background())
	go changeStream(routineCtx, cursor)

}

func (client MongoClient) GetAvailableJobs(jobID string) ([]models.Job, error) {
	
	var result []models.Job
	
	coll				:=	client.session.Database(client.config.DatabaseName).Collection(COLLECTION_JOBS)
	jobIDObjectID,	_	:=	primitive.ObjectIDFromHex(jobID)

	filter	:=	bson.M{ 
		"status": 100  ,
		"permission": "public",
		"job_id": bson.M{ 
			"$ne" : jobIDObjectID,
		},
	}

	opts				:=	options.Find().SetSort(bson.D{{ "created_at", 1 }})
	cursor, err			:=	coll.Find(context.TODO(), filter, opts)

	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		log.Fatal(err)
	}

	return result, err

}

func (client MongoClient) GetJobInformation(jobID string) (models.Job, error) {
	
	var result models.Job

	coll				:=	client.session.Database(client.config.DatabaseName).Collection(COLLECTION_JOBS)
	jobIDObjectID,	_	:=	primitive.ObjectIDFromHex(jobID)
	filter				:=	bson.M{ "job_id": jobIDObjectID }
	opts				:=	options.FindOne().SetSort(bson.D{{ "created_at", 1 }})
	err					:=	coll.FindOne(context.TODO(), filter, opts).Decode(&result)
	
	if err != nil {
		log.Fatal(err)
	}

	return result, err

}
