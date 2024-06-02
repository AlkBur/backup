package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	Port int
}

func Server() (cfg config) {
	cfg.Port = 8000

	mydir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error get current working directory: ", err)
	}

	f, err := os.ReadFile(filepath.Join(mydir, "configs", "server.json"))
	if err != nil {
		log.Fatal("Error when opening config file: ", err)
	}
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return
}
