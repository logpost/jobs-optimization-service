package main

import (

	"github.com/logpost/jobs-optimization-service/utility"
	"github.com/logpost/jobs-optimization-service/config"
	"github.com/logpost/jobs-optimization-service/adapter"
	"github.com/logpost/jobs-optimization-service/routes"
	"github.com/logpost/jobs-optimization-service/logpost"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

)

var (
	databaseConfig	config.DatabaseConfig
	appConfig		config.AppConfig
	osrmConfig		config.OSRMConfig
	PORT			string
	logposter		logpost.LOGPOSTER
	// jobs			[]models.Job
	// pipeJobs		chan []models.Job
)

func init() {

	var secretConfig config.Config
	secretConfig.Read()

	configuration	:=	secretConfig.GetConfig()

	appConfig		=	config.AppConfig{
		Kind:			configuration.Kind,
		Port:			configuration.Port,
		OriginAllowed:	configuration.OriginAllowed,
	}

	databaseConfig	=	config.DatabaseConfig{
		DatabaseURI:	configuration.DatabaseURI,
		DatabaseName:	configuration.DatabaseName,
	}
	
	osrmConfig		=	config.OSRMConfig{
		OSRMBackendURI:	configuration.OSRMBackendURI,
	}

	PORT	=	utility.Getenv("PORT",	appConfig.Port)

}

func main() {

	logposter	=	logpost.CreateOSRMConnection(osrmConfig.OSRMBackendURI)
	mongoClient	:=	adapter.CreateMongoConnection(databaseConfig)
	mongoClient.WatchJobs()

	e	:=	route.Init(&mongoClient, &logposter)
	
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: appConfig.OriginAllowed,
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Logger.Fatal(e.Start(":" + PORT))
} 