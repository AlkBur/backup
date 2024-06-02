package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Notification struct {
	Token string
	URL   string
}

func NewNotification() (cfg Notification) {
	mydir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error get current working directory: ", err)
	}

	f, err := os.ReadFile(filepath.Join(mydir, "configs", "notification.json"))
	if err != nil {
		log.Fatal("Error when opening config file: ", err)
	}
	err = json.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return
}
