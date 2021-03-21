package usecase

import (
	"strconv"
	"log"
	"fmt"
	"net/http"

	"github.com/logpost/jobs-optimization-service/adapter"
	"github.com/logpost/jobs-optimization-service/utility"
	"github.com/logpost/jobs-optimization-service/logpost"

	"github.com/labstack/echo/v4"
	
)

func SuggestJobs(mongoClient *adapter.MongoClient, logposter *logpost.LOGPOSTER) echo.HandlerFunc {

	return	func(c echo.Context) (err error) {
		
		jobID			:=	c.Param("job_id")
		maxHop, err		:=	strconv.Atoi(c.QueryParam("hop"))

		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jobPicked, err	:=	mongoClient.GetJobInformation(jobID)
		
		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jobs, err		:=	mongoClient.GetAvailableJobs(jobID)
		
		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jobsClone		:=	jobs
		jobsClone		=	append(jobsClone, jobPicked)
		adjJobs			:=	utility.TransformToAdjacencyList(jobsClone)

		logposter.SuggestJobsByHop(adjJobs, jobPicked, &jobs, maxHop)
		fmt.Printf("\n\n\n\n#### RESULT \n\n %+v", logposter.Result)
		return c.JSON(http.StatusOK, logposter.Result)

	}
	
}