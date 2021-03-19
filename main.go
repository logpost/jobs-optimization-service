package main

import (
	// "os"
	"fmt"
	"github.com/logpost/jobs-optimization-service/config"
	"github.com/logpost/jobs-optimization-service/mongo"
	// "github.com/logpost/logpost-suggestion-algorithm"
)

var DatabaseConfig config.DatabaseConfig
	
func init() {

	var secretConfig config.Config
	secretConfig.Read()

	configuration	:=	secretConfig.GetConfig()

	databaseConfig	=	config.DatabaseConfig{
		DatabaseURI:	configuration.DatabaseURI,
		DatabaseName:	configuration.DatabaseName,
	}

	fmt.Println(databaseConfig)

}

func main() {

	client	:=	mongo.Connection(databaseConfig)
	client.WatchCollection("jobs")
	
}