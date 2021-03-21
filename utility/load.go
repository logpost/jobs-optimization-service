package utility

import (
	"encoding/json"
	"io/ioutil"
	"github.com/logpost/jobs-optimization-service/models"
)


// LoadJSON is method for loading JSON file 
func LoadJSON() []models.JobExpected {
	var data models.GetterExpected
	
	readFile, _ := ioutil.ReadFile("./google-maps-response-raw.json")
	_ = json.Unmarshal([]byte(readFile), &data)

	saveFile, _ := json.MarshalIndent(data.Getter, "", " ")
	_ = ioutil.WriteFile("google-maps-response-parsed.json", saveFile, 0644)

	return data.Getter
}