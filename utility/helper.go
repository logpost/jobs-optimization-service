package utility

import (
	"os"
	"fmt"
	"github.com/logpost/jobs-optimization-service/models"
)

func Getenv(key, fallback string) string {

	value	:=	os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value

}

func TransformToAdjacencyList(jobs []models.Job) (map[string]*models.Job) {

	adjJobs	:=	make(map[string]*models.Job)

	for	_, job	:=	range (jobs)	{
		jobTemp	:= job
		adjJobs[job.JobID.Hex()] = &jobTemp
		fmt.Println("%+v", adjJobs[job.JobID.Hex()])
	}

	for _, v := range adjJobs {
		fmt.Printf("%+v\n", v)
	}

	return adjJobs

}