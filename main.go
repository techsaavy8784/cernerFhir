package main

import (
	"fmt"

	"net/http"
	"os"

	//"docs/docs.go"
	//"github.com/swaggo/http-swagger" // http-swagger middleware

	// http-swagger middleware
	h "github.com/dhf0820/cernerFhir/pkg/handler"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	//m "pkg/model"
)

// http-swagger middleware

// @title FHIR API for ChartArchive
// @version .2
// @description fhirgo creates an interface to various FHIR systems. initially it is Cerner.
// There is no security at this time so only fake PHI is being used.
// Supports patient/encounters/documents(Diagnostic Reports)

var (
	config  *m.Config
	name    string
	mode    string
	port    string
	mainErr error
)

func main() {
	fmt.Println("Start app")
	godotenv.Load()

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true

	log.SetFormatter(Formatter)
	log.SetLevel(log.DebugLevel)
	log.Infof("Getting MONGODB url")
	mongoURL, err := os.LookupEnv("MONGODB")
	if err != true {
		log.Errorf("Error in MONGODB lookup: %v", err)
		log.Fatal(err)
	}

	log.Infof("Calling Initialize all")
	config = m.InitializeAll(mongoURL)

	port := config.Port()
	mode := config.Mode()
	router := h.NewRouter()
	config.SetRouter(router)
	log.Infof("Serving %s FHIR interface VERSION %s in %s mode on port: %s", config.Source(), config.ServerVersion(), mode, port)
	log.Infof("URL to use: %s", config.BaseUrl())
	log.Infof("Image URL: %s", config.ImageURL())
	log.Infof("Ca Image URL: %s", config.Env("caImageURL"))
	log.Infof("Ca Server url: %s", config.Env("caServerURL"))
	mainErr = http.ListenAndServe(port, router)
	if mainErr != nil {
		log.Errorf("Main error: %v", mainErr)
	}

}
