package usecase

import (
	// "strconv"
	"log"
	"fmt"
	"net/http"

	"github.com/logpost/jobs-optimization-service/adapter"
	"github.com/logpost/jobs-optimization-service/utility"
	"github.com/logpost/jobs-optimization-service/logpost"

	"github.com/labstack/echo/v4"
	
)

type JobInfo struct {
	JobID	string		`json:"job_id"`
	Hop		int			`json:"hop"`
}

func SuggestJobs(mongoClient *adapter.MongoClient, logposter *logpost.LOGPOSTER) echo.HandlerFunc {

	return	func(c echo.Context) (err error) {

		jobInfo := new(JobInfo)

		if err = c.Bind(jobInfo); err != nil {
			return
		}

		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jobPicked, err	:=	mongoClient.GetJobInformation(jobInfo.JobID)
		
		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jobs, err		:=	mongoClient.GetAvailableJobs(jobInfo.JobID)
		
		if err != nil {
			log.Fatal(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		jobsClone		:=	jobs
		jobsClone		=	append(jobsClone, jobPicked)
		adjJobs			:=	utility.TransformToAdjacencyList(jobsClone)

		logposter.SuggestJobsByHop(adjJobs, jobPicked, &jobs, jobInfo.Hop)
		fmt.Printf("\n\n\n\n#### RESULT \n\n %+v", logposter.Result)
		return c.JSON(http.StatusOK, logposter.Result)

	}
	
}