package source_manager_service

import (
	"bufio"
	"encoding/csv"
	"github.com/palantir/stacktrace"
	"io"
	"log"
	"models"
	"os"
	"services"
	"sync"
)

type SourceManagerService struct {
	ConfSrv *services.ConfService
	Dal     *services.DbService
}

var sourceManagerServiceInstance *SourceManagerService
var sourceManagerServiceOnce sync.Once

func NewSourceManagerService() *SourceManagerService {
	sourceManagerServiceOnce.Do(func() {
		sourceManagerServiceInstance = &SourceManagerService{
			ConfSrv: services.NewConfService(),
			Dal:     services.NewDbService(),
		}
	})
	return sourceManagerServiceInstance
}

func (s *SourceManagerService) LoadCsvSourceFile(fileName string) error {
	if fileName == "" {
		fileName = s.ConfSrv.Config.CSVSourceFile
	}
	csvFile, err := os.Open(fileName)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to open the source CSV file")
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	//  omit first line with captions
	_, err = reader.Read()
	if err != nil {
		return stacktrace.Propagate(err, "Failed to omit caption line in the source CSV file")
	}

	var restaurants []models.Restaurant
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return stacktrace.Propagate(err, "Failed to read the line in the source CSV file")
		}
		restaurants = append(restaurants, models.Restaurant{Name: line[0], Type: line[1],
			Phone: line[2], Location: line[3]})
	}

	err = s.Dal.InsertRestaurants(restaurants)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to insert the restaurants to the database")
	}

	log.Println("Init restaurants was completed successfully")
	return nil
}
