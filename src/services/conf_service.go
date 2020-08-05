package services

import (
	"encoding/json"
	"log"
	"models"
	"os"
	"sync"
)

type ConfService struct {
	Config models.Configuration
}

var confServiceInstance *ConfService
var confServiceOnce sync.Once

func NewConfService() *ConfService {
	confServiceOnce.Do(func() {
		confServiceInstance = &ConfService{}
		err := confServiceInstance.readConfiguration()
		if err != nil {
			log.Fatal("Failed init configuration file", err)
		}
	})
	return confServiceInstance
}

func (c *ConfService) readConfiguration() error {
	log.Println("Start init configuration file")
	filename := "src\\config\\tsconfig.json"
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	var conf models.Configuration
	err = decoder.Decode(&conf)
	if err != nil {
		return err
	}
	c.Config = conf
	log.Println("Init of the configuration file was completed successfully")
	return nil
}
