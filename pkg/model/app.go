package model

/*
import (
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/dhf0820/cernerFhir/pkg/model"

	"github.com/dhf0820/cernerFhir/pkg/storage"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	"github.com/gorilla/mux"
)

type App struct {
	Router        *mux.Router
	DB            *storage.MongoDB
	EMR           *EMR
	EMRURL        string
	EMRConnection *fhir.Connection
	Port          string
}

var currentApp *App

func (a *App) Initialize(url string) {
	db, err := storage.Open("")
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
	a.Router = NewRouter()
	a.Port = ":9000" // will come from config
	//a.Router = mux.NewRouter()
	//var emr model.EMR

	a.EMRConnection = initializeEMR("https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/")
	currentApp = a
	// emr, err := model.InitializeEMR()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// a.EMRConnection = emr

}

func (a *App) Run() {
	if err := http.ListenAndServe(a.Port, a.Router); err != nil {
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}
}

func Current() *App {
	// check if app is nill and create app if necessary
	return currentApp
}

func initializeEMR(url string) *fhir.Connection {
	c := fhir.New(url)
	err := checkEMRConnection(c)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func checkEMRConnection(c *fhir.Connection) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, "metadata")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json+fhir")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("EMRInit: did not connect %s", resp.Status)
	}
	return nil

}
*/
