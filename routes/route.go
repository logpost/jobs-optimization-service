package route

import (

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/logpost/jobs-optimization-service/usecase"
	"github.com/logpost/jobs-optimization-service/adapter"
	"github.com/logpost/jobs-optimization-service/logpost"

)

func Init(mongoClient *adapter.MongoClient, logposter *logpost.LOGPOSTER) *echo.Echo {

	e	:=	echo.New()

	r	:=	e.Group("/job-opts")
	{
		r.GET("/healthcheck", func(c echo.Context) error { 
			return c.String(http.StatusOK, "Service is Alive 🥳")
		})

		r.POST("/suggest", usecase.SuggestJobs(mongoClient, logposter))
	}

	return e

}