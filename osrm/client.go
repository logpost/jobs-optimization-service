package osrm

import (
	"strconv"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"github.com/karmadon/gosrm"
	"github.com/paulmach/go.geo"
	"github.com/logpost/jobs-optimization-service/models"
)

// OSRM Struct is data for access gosrm.
type OSRM struct {
	client  *gosrm.OsrmClient
}

// CreateOSRM is function for create client.
func (osrm *OSRM) CreateOSRM(uri string) {
	options := &gosrm.Options{
		Url:            url.URL{ Host: uri },
		Service:        gosrm.ServiceRoute,
		Version:        gosrm.VersionFirst,
		Profile:        gosrm.ProfileDriving,
		RequestTimeout: 5,
	}

	osrm.client = gosrm.NewClient(options)
}

// GetRouteInfo function for get information routing.
func (osrm *OSRM) GetRouteInfo(source, dest *models.Location) *gosrm.OSRMResponse {	

	routeRequest := &gosrm.RouteRequest{
		Coordinates: geo.PointSet{
			{ source.Longitude, source.Latitude },
			{ dest.Longitude, dest.Latitude },
		},
	}
	
	response, err := osrm.client.Route(routeRequest)

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	loggingRouteToJSON(response, source, dest)

	return response
}

// loggingLastRoute function for logging response from route btw points.
func loggingRouteToJSON(response *gosrm.OSRMResponse, source, dest *models.Location) {
	sourceLatlong	:= strconv.FormatFloat(source.Latitude, 'f', 3, 64) + "," + strconv.FormatFloat(source.Longitude, 'f', 3, 64)
	destLatlong 	:= strconv.FormatFloat(dest.Latitude, 'f', 3, 64) + "," + strconv.FormatFloat(dest.Longitude, 'f', 3, 64)

	saveFile, _ := json.MarshalIndent(response, "", " ")
	outputPath	:= "output/" + sourceLatlong + "_" + destLatlong + "-response-coordinates-osrm.json"

	_ = ioutil.WriteFile(outputPath, saveFile, 0644)
}