package utility

import (
	"time"
	"github.com/logpost/jobs-optimization-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InTimeSpan func check timimg is in btw these time.
func CheckNotInTimeSpan(start, end, check time.Time) bool {
	return	!(check.After(start) && check.Before(end))
}

func CheckAvaliable(jobIterator *models.Job) bool {

	objDefault, _	:=	primitive.ObjectIDFromHex("000000000000000000000000")

	if jobIterator.CarrierID == objDefault && jobIterator.Permission == "public" && jobIterator.Status == 100 {
		return true
	}

	return false
}

func JobFilterAllConditions(jobPicked *models.Job, jobIterator *models.Job) bool {

	startAt			:= jobPicked.PickupDate
	endAt			:= jobPicked.DropoffDate

	isNotInTimeSpan	:=	CheckNotInTimeSpan(startAt, endAt, jobIterator.PickupDate)
	isAvaliable		:=	CheckAvaliable(jobIterator)

	return isNotInTimeSpan && isAvaliable
}

func JobsFiltering(jobPicked models.Job, jobs *[]models.Job) ([]models.Job, int) {
	
	var result []models.Job

	for _, currentJob := range (*jobs) {
		if ok := JobFilterAllConditions(&jobPicked, &currentJob); ok {
			result = append(result, currentJob)
		}
	}
	
	return result, len(result)
}
