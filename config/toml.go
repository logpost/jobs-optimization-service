package config

import (
	"log"
	
	"github.com/logpost/jobs-optimization-service/utility"

	"github.com/BurntSushi/toml"
)

type OSRMConfig	struct {
	OSRMBackendURI	string
}

type DatabaseConfig struct {
	DatabaseURI		string
	DatabaseName	string
}

type AppConfig	struct {
	Kind			string
	Port			string
	OriginAllowed	[]string
}

type Configuration struct { 
	Kind			string
	Port			string
	OriginAllowed 	[]string
	DatabaseURI		string
	DatabaseName	string
	OSRMBackendURI	string
}

// Config is struct for parse data from toml file
type Config struct {
	Stagging		Configuration
	Development		Configuration
}

// Read and parse the configuration file
func (c *Config) Read() {

	if _, err	:=	toml.DecodeFile("conf/config.toml", &c); err	!=	nil {
		log.Fatal(err)
	}
	
}

func (c *Config) GetConfig() Configuration {

	switch KIND := utility.Getenv("KIND", "development"); KIND {
		case "development":
			return	c.Development
		case "stagging":
			return	c.Stagging
		default:
			panic("required environment vairble !" + KIND)
	}

}


