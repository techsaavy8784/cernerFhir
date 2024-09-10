package model

/*
import (
	"fmt"
	"log"
	"net/http"
	"time"

	fhir "github.com/dhf0820/cernerFhir/fhirongo"
)

type EMR struct {
	URL        string
	Connection *fhir.Connection
}

var emr EMR

func InitializeEMR(url) (*EMR, error) {
	emr.Connection = fhir.New(url)
	err := checkEMRConnection(emr.Connection)
	if err != nil {
		log.Fatal(err)
	}
	return &emr, nil
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
