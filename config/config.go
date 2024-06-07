package config

import (
	"encoding/json"
	"log"
	"os"
)

type Server struct {
	Port string
}

type Log struct {
	Level string
}

// Configuration Stores the main configuration for the application
type Config struct {
	Server Server
	Log    Log
}

// ReadConfig will read the configuration json file to read the parameters
// which will be passed in the config file
func ReadConfig(fileName string) (Config, error) {
	cfg := Config{}
	configFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Print("Unable to read config file, switching to flag mode")
		return cfg, err
	}
	err = json.Unmarshal(configFile, &cfg)
	if err != nil {
		log.Print("Invalid JSON, expecting port from command line flag")
		return cfg, err
	}
	return cfg, nil
}
