package usecase

import (
	// "strconv"
	"fmt"
	"net/http"

	"github.com/logpost/jobs-optimization-service/adapter"
	"github.com/logpost/jobs-optimization-service/utility"
	"github.com/logpost/jobs-optimization-service/logpost"
	"github.com/logpost/jobs-optimization-service/models"

	"github.com/labstack/echo/v4"
	
)

type Body struct {
	JobID			string			`json:"job_id"`
	Hop				int				`json:"hop"`
	OriginLocation	models.Location	`json:"origin_location"`
}

func SuggestJobs(mongoClient *adapter.MongoClient, logposter *logpost.LOGPOSTER) echo.HandlerFunc {

	return	func(c echo.Context) (err error) {

		body := new(Body)

		if err = c.Bind(body); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Error : Body is invalid.")
		}
		
		jobPicked, err	:=	mongoClient.GetJobInformation(body.JobID)
		
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		jobs, err		:=	mongoClient.GetAvailableJobs(body.JobID)
		
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		jobsClone		:=	jobs
		jobsClone		=	append(jobsClone, jobPicked)
		adjJobs			:=	utility.TransformToAdjacencyList(jobsClone)

		originLocation	:=	models.CreateLocation(float64(body.OriginLocation.Latitude), float64(body.OriginLocation.Longitude))

		err				=	logposter.SuggestJobsByHop(originLocation, adjJobs, jobPicked, &jobs, body.Hop)

		if err != nil {
			fmt.Println(err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Error : No route in the system.")
		}

		return c.JSON(http.StatusOK, logposter.Result)

	}
	
}