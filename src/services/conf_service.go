package services

import (
	"encoding/json"
	"fmt"
	"log"
	"models"
	"os"
	"strings"
	"sync"
)

type ConfService struct {
	Config models.Configuration
}

var confServiceInstance *ConfService
var confServiceOnce sync.Once

func NewConfService() *ConfService {
	if len(os.Args) == 1 || len(os.Args) > 1 && !strings.HasSuffix(os.Args[1], "-test.v") {
		confServiceOnce.Do(func() {
			confServiceInstance = &ConfService{}
			fmt.Println("normal run")
			filename := "src\\config\\tsconfig.json"
			var err error
			confServiceInstance.Config, err = confServiceInstance.ReadConfiguration(filename)
			if err != nil {
				log.Fatal("Failed init configuration file", err)
			}
		})
	}
	return confServiceInstance
}

func (c *ConfService) ReadConfiguration(filename string) (models.Configuration, error) {
	log.Println("Start init configuration file")
	file, err := os.Open(filename)
	if err != nil {
		return models.Configuration{}, err
	}
	decoder := json.NewDecoder(file)
	var conf models.Configuration
	err = decoder.Decode(&conf)
	if err != nil {
		return models.Configuration{}, err
	}
	log.Println("Init of the configuration file was completed successfully")
	return conf, nil
}
