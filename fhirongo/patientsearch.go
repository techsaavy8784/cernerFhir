package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"

	//"github.com/davecgh/go-spew/spew"

	log "github.com/sirupsen/logrus"
)

// PatientSearch will search for a patient based on the query string
// identifier, family, given, birthdate, gender, address, telecom
// i.e. family=Argonaut&given=Jason

//TODO: Figure out how to handle the limit and offset.  Previously cashed all in mongo then returned from there.
func (c *Connection) PatientSearch(query string) (*PatientResult, error) {
	log.Infof("queryString: %s\n", query)
	qry := fmt.Sprintf("/Patient?%s", query)
	log.Infof("Finual url to query: %s\n", qry)
	startTime := time.Now()
	//b, err := c.Query(fmt.Sprintf("/Patient?%s", query))
	log.Infof("Query time: %s", time.Since(startTime))
	bytes, err := c.Query(qry)

	if err != nil {

		return nil, fmt.Errorf("Query %s failed: %s", query, err.Error())
	}
	
	//fmt.Printf("\n\n\n@@@ RAW Patient: %s\n\n\n", pretty.Pretty(b))
	// prettyJSON, err := json.MarshalIndent(b, "", "    ")
	// if err != nil {
	// 	fmt.Printf("MarshalIndent failed: %s\n", err.Error())
	// 	return nil, err
	// }

	startTime = time.Now()
	data := PatientResult{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, fmt.Errorf("PatientSearch ummarshal : %s", err.Error())
	}
	log.Infof("Unmarshal time: %s", time.Since(startTime))
	//fmt.Printf("Response: %s\n", spew.Sdump(data))
	return &data, nil
}

// PatientResult is a patient search result
type PatientResult struct {
	SearchResult
	Entry []struct {
		EntryPartial
		Patient Patient `json:"resource"`
	} `json:"entry"`
}
