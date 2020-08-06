package services

import (
	"encoding/json"
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
	confServiceOnce.Do(func() {
		confServiceInstance = &ConfService{}
		var fileName string
		if len(os.Args) == 1 || len(os.Args) > 1 && !strings.HasSuffix(os.Args[1], "-test.v") {
			fileName = "src\\config\\tsconfig.json"
		} else {
			fileName = "C:\\workspace\\dragontail_ex\\src\\config\\tsconfig.json"
		}
		var err error
		confServiceInstance.Config, err = confServiceInstance.ReadConfiguration(fileName)
		if err != nil {
			log.Fatal("Failed init configuration file", err)
		}
	})
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
