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
)

func init() {

	var secretConfig config.Config
	secretConfig.Read()

	configuration	:=	secretConfig.GetConfig()

	appConfig		=	config.AppConfig{
		Kind:			configuration.Kind,
		Host:			configuration.Host,
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

	e	:=	route.Init(&mongoClient, &logposter)
	
	e.Debug				=	false
	e.HideBanner		=	true

	e.Use(middleware.Logger(), middleware.Recover(), middleware.Gzip())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: appConfig.OriginAllowed,
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Logger.Fatal(e.Start(appConfig.Host + ":" + PORT))
} 