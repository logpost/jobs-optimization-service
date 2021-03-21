package logpost

import (
	"strconv"
	"sync"
	"time"
	"fmt"
	"container/heap"
	"github.com/logpost/jobs-optimization-service/pqueue"
	"github.com/logpost/jobs-optimization-service/utility"
	"github.com/logpost/jobs-optimization-service/models"
	"github.com/logpost/jobs-optimization-service/osrm"
)

// LOGPOSTER.var	OSRMClient	osrm.OSRM
type LOGPOSTER	struct {
	OSRMClient	osrm.OSRM
	adjJobs		map[string]*models.Job
	Result		Result
}

type Result	struct	{
	Summary		map[string]Summary		`json:"summary"`
	History		map[string]models.Job	`json:"history"`
}

type Summary struct {
	SumCost				float64			`json:"sum_cost"`
	SumOffer			float64			`json:"sum_offer"`
	Profit				float64			`json:"profit"`
	DistanceToOrigin	float64			`json:"distance_to_origin"`
	StartDate			string			`json:"start_date"`
	EndDate				string			`json:"end_date"`
}

// MinimumCostBuffer struct for sending to minimum pipe
type MinimumCostBuffer struct {
	minimumJobID				string
	minimumEndingCost			float64
	minimumDistanceToOrigin		float64
	minimumCost					float64
	minimumPrepare				float64
}

func timeTrack(start time.Time) {
	elapsed	:=	time.Since(start)
	fmt.Printf("\nTOOK:\t\t%s\n", elapsed)
}

func (logposter *LOGPOSTER) getJobMinimumCost(curentLocation *models.Location, originLocation *models.Location, minimumCostPipe chan MinimumCostBuffer, waitGroup *sync.WaitGroup, jobs *[]models.Job, startIndex int, endIndex int) {
	
	defer waitGroup.Done()
	
	minimumJobID			:=	""
	minimumCost				:=	9999999.999
	minimumPrepare			:=	0.0
	minimumEndingCost		:=	0.0
	minimumDistanceToOrigin	:=	0.0
 
	for index := startIndex; index < endIndex; index++ {
		
		if	!logposter.adjJobs[(*jobs)[index].JobID.Hex()].Visited {
		
			predictingPickUpLocation	:=	models.CreateLocation(logposter.adjJobs[(*jobs)[index].JobID.Hex()].PickUpLocation.Latitude,	logposter.adjJobs[(*jobs)[index].JobID.Hex()].PickUpLocation.Longitude)
			predictingDropOffLocation	:=	models.CreateLocation(logposter.adjJobs[(*jobs)[index].JobID.Hex()].DropOffLocation.Latitude,	logposter.adjJobs[(*jobs)[index].JobID.Hex()].DropOffLocation.Longitude)

			prepareRouting				:=	logposter.OSRMClient.GetRouteInfo(curentLocation,	&predictingPickUpLocation)
			endingRouting				:=	logposter.OSRMClient.GetRouteInfo(&predictingDropOffLocation,	originLocation)

			if prepareRouting != nil && endingRouting != nil {
				prepareRoutingDistance	:=	prepareRouting.Routes[0].Distance
				endingRoutingDistance	:=	endingRouting.Routes[0].Distance

				preparingCost			:=	utility.GetDrivingCostByDistance(prepareRoutingDistance, 0)
				endingCost				:=	utility.GetDrivingCostByDistance(endingRoutingDistance, 0)
				
				sumaryPredictingCost	:=	preparingCost + logposter.adjJobs[(*jobs)[index].JobID.Hex()].Cost + endingCost

				if	minimumCost > sumaryPredictingCost {
					minimumCost				=	sumaryPredictingCost
					minimumPrepare			=	preparingCost
					minimumDistanceToOrigin	=	endingRoutingDistance
					minimumEndingCost		=	endingCost
					minimumJobID			=	(*jobs)[index].JobID.Hex()
				}
			}
		} 
	}

	fmt.Printf("\n### MINIMUM PREDICT:\nINDEX:\t\t%s\nCOST:\t\t%f\nCOST_PREPARE:\t%f\nCOST_ENDING:\t%f\n", minimumJobID, minimumCost, minimumPrepare, minimumEndingCost)
	 
	buffer	:=	MinimumCostBuffer{
		minimumJobID, minimumEndingCost, minimumDistanceToOrigin, minimumCost, minimumPrepare,
	}

	minimumCostPipe	<- buffer

}

func getActualJobMinimumCost(minimumCostPipe chan MinimumCostBuffer, actualJobMinimumCostPipe chan MinimumCostBuffer, waitGroup *sync.WaitGroup) {

	defer waitGroup.Done()

	var minimumCostOne,	minimumCostTwo	MinimumCostBuffer
	goI, goII, goIII, goIV := <-minimumCostPipe, <-minimumCostPipe, <-minimumCostPipe, <-minimumCostPipe

	if goI.minimumCost > goII.minimumCost {
		minimumCostOne	=	goI
	} else {
		minimumCostOne	=	goII
	}

	if goIII.minimumCost > goIV.minimumCost {
		minimumCostTwo	=	goIII
	} else {
		minimumCostTwo	=	goIV
	}
	
	if minimumCostOne.minimumCost < minimumCostTwo.minimumCost {
		actualJobMinimumCostPipe <- minimumCostOne
	} else {
		actualJobMinimumCostPipe <- minimumCostTwo
	}
}

func (logposter *LOGPOSTER) getJobMinimumCostParallel(jobPickedLocation *models.Location, originLocation *models.Location, jobs *[]models.Job) (string, float64, float64, float64) {

	minimumCostPipe				:=	make(chan MinimumCostBuffer)
	actualJobMinimumCostPipe	:=	make(chan MinimumCostBuffer)

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go logposter.getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &waitGroup, jobs, 0, len(*jobs)/4)

	waitGroup.Add(1)
	go logposter.getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &waitGroup, jobs, len(*jobs)/4, len(*jobs)/4 * 2)

	waitGroup.Add(1)
	go logposter.getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &waitGroup, jobs, len(*jobs)/4 * 2, len(*jobs)/4 * 3)

	waitGroup.Add(1)
	go logposter.getJobMinimumCost(jobPickedLocation, originLocation, minimumCostPipe, &waitGroup, jobs, len(*jobs)/4 * 3, len(*jobs))

	waitGroup.Add(1)
	go getActualJobMinimumCost(minimumCostPipe, actualJobMinimumCostPipe, &waitGroup)

	actualJobMinimumCost		:=	<- actualJobMinimumCostPipe
	minimumJobID				:=	actualJobMinimumCost.minimumJobID
	minimumEndingCost			:=	actualJobMinimumCost.minimumEndingCost
	minimumDistanceToOrigin		:=	actualJobMinimumCost.minimumDistanceToOrigin
	minimumCost					:=	actualJobMinimumCost.minimumCost
	
	waitGroup.Wait()

	return minimumJobID, minimumEndingCost, minimumDistanceToOrigin, minimumCost
	
}

// LOGPOSTER.OSRMClient.CreateOSRM("http://osrm:5001/")
// CreateOSRMClient("http://127.0.0.1:5001/")

func CreateOSRMConnection(URL string) LOGPOSTER {

	var logposter	LOGPOSTER
	logposter.OSRMClient	=	osrm.OSRM{}
	logposter.OSRMClient.CreateOSRM(URL)
	
	fmt.Println("Connected to OSRM backend 🎃")

	return logposter

}

func (logposter *LOGPOSTER)	SuggestJobsByHop(adjJobs map[string]*models.Job, jobFirstPicked models.Job, jobs *[]models.Job, maxHop int) {

	var minimumJobID			string
	var minimumEndingCost		float64
	var minimumDistanceToOrigin	float64
	var Queue					pqueue.PriorityQueue
	
	heap.Init(&Queue)

	sumCost			:=	0.0
	sumOffer		:=	0.0
	currentHop		:=	0
	workingDays 	:=	1
	maxWorkingDays	:=	-1
	startDay	 	:=	time.Now()
	endDay			:=	time.Now()

	start						:=	time.Now()

	logposter.Result.Summary	=	make(map[string]Summary)
	logposter.Result.History	=	make(map[string]models.Job)
	logposter.adjJobs			=	adjJobs

	
	jobsFiltered, _	:=	utility.JobsFiltering(jobFirstPicked, jobs)
	jobs			=	&jobsFiltered

	for _, job := range (*jobs) {
		adjJobs[job.JobID.Hex()].Cost = utility.GetDrivingCostByDistance(adjJobs[job.JobID.Hex()].Distance, adjJobs[job.JobID.Hex()].Weight)
	}

	adjJobs[jobFirstPicked.JobID.Hex()].Cost = utility.GetDrivingCostByDistance(adjJobs[jobFirstPicked.JobID.Hex()].Distance, adjJobs[jobFirstPicked.JobID.Hex()].Weight)

	// Initial data selected by user
	originLocation	:=	models.CreateLocation(float64(14.7995081), float64(100.6533706))
	curentLocation	:=	originLocation

	// Starting suggestion algorithm
	heap.Push(&Queue, &pqueue.Item{
		JobID:	jobFirstPicked.JobID.Hex(),
	})

	for Queue.Len() > 0 {
		
		currentHop++
		
		jobPicked	:=	heap.Pop(&Queue).(*pqueue.Item)
		logposter.adjJobs[jobPicked.JobID].Visited	=	true

		jobsFiltered, _		:=	utility.JobsFiltering(*logposter.adjJobs[jobPicked.JobID], jobs)
		
		jobs				=	&jobsFiltered
		
		sumCost				+=	logposter.adjJobs[jobPicked.JobID].Cost
		sumOffer			+=	logposter.adjJobs[jobPicked.JobID].OfferPrice
		endDay				=	logposter.adjJobs[jobPicked.JobID].DropoffDate

		logposter.Result.History[strconv.Itoa(currentHop)]	=	*logposter.adjJobs[jobPicked.JobID]
		logposter.Result.Summary[strconv.Itoa(currentHop)]	=	Summary{
			SumCost:			sumCost,
			SumOffer:			sumOffer,
			Profit:				sumOffer - sumCost,
			DistanceToOrigin:	jobPicked.DistanceToOrigin,
			StartDate:			jobFirstPicked.PickupDate.String(),
			EndDate:			adjJobs[jobPicked.JobID].DropoffDate.String(),
		}

		jobPickedLocation	:=	models.CreateLocation(logposter.adjJobs[jobPicked.JobID].PickUpLocation.Latitude, logposter.adjJobs[jobPicked.JobID].PickUpLocation.Longitude)
		prepareRouting		:=	logposter.OSRMClient.GetRouteInfo(&curentLocation, &jobPickedLocation)

		if prepareRouting	!=	nil {
			preparingDistance	:=	prepareRouting.Routes[0].Distance
			preparingCost		:=	utility.GetDrivingCostByDistance(preparingDistance, 0)
			sumCost				+=	preparingCost
		}

		if currentHop	<=	maxHop {
			
			minimumJobID, minimumEndingCost, minimumDistanceToOrigin, _	=	logposter.getJobMinimumCostParallel(&jobPickedLocation, &originLocation, jobs)

			if minimumJobID	!=	"" {

				heap.Push(&Queue, &pqueue.Item{
					JobID:	minimumJobID,
					DistanceToOrigin:	minimumDistanceToOrigin,
				})
				
				curentLocation	=	models.CreateLocation(logposter.adjJobs[minimumJobID].DropOffLocation.Latitude, logposter.adjJobs[minimumJobID].DropOffLocation.Latitude)

			}
		}

		if currentHop > maxHop	||	Queue.Len() == 0 || len((*jobs)) == 0 {
			
			if currentHop == 1 {
				predictingDropOffLocation	:=	models.CreateLocation(logposter.adjJobs[jobPicked.JobID].PickUpLocation.Latitude, logposter.adjJobs[jobPicked.JobID].PickUpLocation.Longitude)
				endingRouting				:=	logposter.OSRMClient.GetRouteInfo(&predictingDropOffLocation, &originLocation)

				if endingRouting	!=	nil {
					endingRoutingDistance	:=	endingRouting.Routes[0].Distance
					minimumDistanceToOrigin	=	endingRoutingDistance
					endingCost				:=	utility.GetDrivingCostByDistance(endingRoutingDistance, 0)
					sumCost					+=	endingCost
				}

			} else {
				sumCost	+=	minimumEndingCost
			}

			// fmt.Println("\nCURRENT_HOP: ", currentHop)

			break

		}

		fmt.Println("\nCURRENT_HOP: ", currentHop)

	}

	fmt.Printf("\n## SUMARY ##\n")

	timeTrack(start)
	
	fmt.Printf("SUM_OFFER:\t\t%f\n",		sumOffer)
	fmt.Printf("SUM_COST:\t\t%f\n",			sumCost)
	fmt.Printf("SUM_PROFIT:\t\t%f\n",		sumOffer - sumCost)
	fmt.Printf("START_DATE:\t\t%s\n",		startDay.String())
	fmt.Printf("END_DATE:\t\t%s\n",			endDay.String())
	fmt.Printf("DISTANCE_TO_ORIGIN:\t%f\n",	minimumDistanceToOrigin)

	fmt.Println("DEBUG: ", Queue, workingDays, maxWorkingDays, startDay, endDay)


}
