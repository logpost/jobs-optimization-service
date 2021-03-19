package config

import (
	"log"
	"os"
	"github.com/BurntSushi/toml"
)

type DatabaseConfig struct {
	DatabaseURI		string
	DatabaseName	string
}

type Configuration struct { 
	Kind			string
	Port			string
	OriginAllowed 	string
	DatabaseURI		string
	DatabaseName	string
}

// Config is struct for parse data from toml file
type Config struct {
	Stagging		Configuration
	Development		Configuration
}

// Read and parse the configuration file
func (c *Config) Read() {

	if _, err	:=	toml.DecodeFile("toml/config.toml", &c); err	!=	nil {
		log.Fatal(err)
	}
	
}

func (c *Config) GetConfig() Configuration {

	switch KIND := os.Getenv("KIND"); KIND {
		case "development":
			return	c.Development
		case "stagging":
			return	c.Stagging
		default:
			panic("required environment vairble !" + KIND)
	}

}


