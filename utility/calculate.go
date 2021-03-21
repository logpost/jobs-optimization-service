package utility

const (
	tonPrice 		=	350.0
	oilPrice 		=	8.0
	driver 			=	400.0
	depreciation 	=	800.0
	rateFixedCost	=	0.05	// 5% of offer
	rateTax			=	0.01	// 1% of offer
	changeMToKM		=	1000.0
)

// GetEnvironmentCostByDay function for get environment cost.
func GetEnvironmentCostByDay(offer float64, day int) float64 {

	costTruck	:=	float64(depreciation * day)
	costDriver	:=	float64(driver * day)
	costFixed	:=	float64(offer) * rateFixedCost 
	costTax		:=	float64(offer) * rateTax
	
	return costDriver + costTruck + costFixed + costTax
}

// GetDrivingCostByDistance function for get driving cost in one time.
func GetDrivingCostByDistance(distance, weight float64) float64 {

	var costOilPrice	float64
	
	if weight > 0 { 
		costOilPrice = oilPrice		// TRUCK DRIVE WHEN LOADING.
	} else {
		costOilPrice = oilPrice / 2	// TRUCK DRIVE WHEN NO LOAD.
	}

	return (distance / changeMToKM ) * costOilPrice
	
}

// GetCostOneJob function for get cost from offer and distance.
func GetCostOneJob(offer, distance float64) float64 {

	day := 1

	oilPriceDeparture	:=	oilPrice
	oilPriceReturn		:=	oilPrice / 2	
	costDriving			:=	distance * (oilPriceDeparture + oilPriceReturn) 
	costEnvironment		:=	GetEnvironmentCostByDay(offer, day)

	return costDriving + costEnvironment
}

// GetOfferFromWeight function for get offer from weight.
func GetOfferFromWeight(weight float64) float64 {

	autoPrice	:=	weight * tonPrice
	return (autoPrice * 0.1) + autoPrice

}