package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"routes_service"
	"source_manager_service"
)

const logFile = "DragonTailEx.log"

func main() {
	initLog()
	log.Println("Starting DragonTail system...")

	sms := source_manager_service.NewSourceManagerService()
	err := sms.LoadCsvSourceFile("")
	if err != nil {
		log.Fatal("Failed to load CSV source file", err)
	}

	rs := routes_service.NewRoutesService()
	err = rs.InitRoutes()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to init the routes: %v", err))
	}
}

func initLog() {
	fmt.Println("Start initializing the log")
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Failed to create log file")
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
