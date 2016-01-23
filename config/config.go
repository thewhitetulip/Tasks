package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//Stores the main configuration for the application
type Configuration struct {
	ServerPort string
}

var err error
var config Configuration

//ReadConfig will read the configuration json file to read the parameters
//which will be passed in the config file
func ReadConfig(fileName string) Configuration {
	configFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Unable to read log file")
	}
	//log.Print(configFile)
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Print(err)
	}
	return config
}
