package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
) 

// Job Struct for create job instance
type Job struct { 
	// Attribute parsed from raw
	JobID				primitive.ObjectID	`json:"job_id" bson:"job_id"`	
	CarrierID			primitive.ObjectID	`json:"carrier_id" bson:"carrier_id"`
	OfferPrice 			float64 			`json:"offer_price" bson:"offer_price"`
	Weight				float64 			`json:"weight" bson:"weight"`
	Duration 			int 				`json:"duration"bson:"duration"`
	WaitingTime			int					`json:"waiting_time" bson:"waiting_time"`
	Distance 			float64 			`json:"distance" bson:"distance"`
	ProductType			string 				`json:"product_type" bson:"product_type"`
	Permission			string 				`json:"permission" bson:"permission"`
	PickupDate			time.Time 			`json:"pickup_date" bson:"pickup_date"`
	DropoffDate			time.Time	 		`json:"dropoff_date" bson:"dropoff_date"`
	PickUpLocation		Location			`json:"pickup_location" bson:"pickup_location"`
	DropOffLocation		Location			`json:"dropoff_location" bson:"dropoff_location"`
	Status				int					`json:"status" bson:"status"`
	// Attribute for running algorithm
	Visited				bool				`json:"visited" bson:"visited"`
	Cost				float64				`json:"cost" bson:"cost"`
}

// Location Struct for mapping location information
type Location struct {
	Latitude			float64				`json:"latitude" bson:"latitude"`
	Longitude			float64				`json:"longitude" bson:"longitude"`
	Address				string				`json:"address" bson:"address"`
	Province			string				`json:"province" bson:"province"`
	District			string				`json:"district" bson:"district"`
	Zipcode				string				`json:"zipcode" bson:"zipcode"`
}

// CreateLocation func do create location's struct
func CreateLocation(latitude float64, longitude float64) Location {
	return Location{
		Latitude:	latitude,
		Longitude:	longitude,
	}
}