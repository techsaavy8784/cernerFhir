package fhirongo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	//"github.com/tidwall/pretty"
	//"github.com/davecgh/go-spew/spew"

	log "github.com/sirupsen/logrus"
)

// GetPatient will return patient information for a patient with id pid
func (c *Connection) GetPatient(pid string) (*Patient, error) {
	if pid == "" {
		msg := "fhir GetPatient param can not be blank"
		log.Errorf(msg)
		return nil, errors.New(msg)
	}
	//log.Infof("FHIR GetPatient url: %s/Patient/%v", c.BaseURL, pid)
	startTime := time.Now()
	qry := fmt.Sprintf("Patient/%v", pid)
	bytes, err := c.Query(qry)
	log.Infof("Query took %s", time.Since(startTime))
	if err != nil {
		msg := fmt.Sprintf("c.Query failed with [%s] err: %s\n",qry, err.Error() )
		log.Error(msg)
		return nil, fmt.Errorf(msg)
	}
//b := *body
	//fmt.Printf("\n\n\n@@@ Patient:22 RAW Patient: %s\n\n\n", pretty.Pretty(b))
	log.Infof("Length of bytes response: %d\n", len(bytes))
	data := Patient{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		fmt.Printf("UnMarshal GetPatient:25\n")
		return nil, err
	}
	//log.Infof("Query returning: %s", spew.Sdump(data))
	return &data, nil
}

func (c *Connection) FindFhirPatient(qry string) (*PatientResult, error) {
	//fmt.Printf("QRY: %s\n", qry)
	//fmt.Printf("With v: Patient%v\n", qry)
	//fmt.Printf("Patient%s\n", qry)
	//fmt.Printf("FHIR FindPatient url: %sPatient?%s\n", c.BaseURL, qry)
	query := fmt.Sprintf("Patient?%s", qry)
	bytes, err := c.Query(query)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("\n\n\n@@@ Patient 15 RAW Patient: %s\n\n\n", pretty.Pretty(b))
	data := PatientResult{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Connection) FindFhirPatients(qry string) (*PatientResult, error) {
	//fmt.Printf("QRY: %s\n", qry)
	//fmt.Printf("With v: Patient%v\n", qry)
	//fmt.Printf("Patient%s\n", qry)
	//fmt.Printf("FHIR FindPatient url: %sPatient?%s\n", c.BaseURL, qry)
	query := fmt.Sprintf("Patient%s", qry)  // The query has the correect seperator(/, ?)
	bytes, err := c.Query(query)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("\n\n\n@@@ Patient 15 RAW Patient: %s\n\n\n", pretty.Pretty(b))
	data := PatientResult{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Connection) NextFhirPatients(url string) (*PatientResult, error) {
	//fmt.Printf("Next retrieving : %s\n", url)
	bytes, err := c.GetFhir(url)
	if err != nil {
		msg := fmt.Sprintf("NextPatient returned error: %s", err.Error())
		log.Errorf("%s", msg)
		return nil, err
	}

	data := PatientResult{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}


// Patient is a FHIR patient
type Patient struct {
	CacheID 		primitive.ObjectID `json:"-" bson:"_id"`
	SessionId		string 			`json:"-" bson:"sessiopn_id"`
	ResourceType    string          `json:"resourceType" bson:"resource_type"`
	ID 				string 			`json:"id" bson:"id"`
	Meta 			MetaData		`json:"meta" bson:"meta"`
	Text            TextData        `json:"text" bson:"text"`
	Identifier      []Identifier    `json:"identifier" bson:"identifier"`
	Active          bool            `json:"active" bson:"active"`
	BirthDate       string          `json:"birthDate" bson:"birth_date"`
	Gender          string          `json:"gender" bson:"gender"`
	DeceasedBoolean bool            `json:"deceasedBoolean" bson:"deceased"`
	CareProvider    []Person        `json:"careProvider" bson:"care_provider"`
	Name            []Name          `json:"name" bson:"name"`
	Address         []Address       `json:"address" bson:"address"`
	Telecom         []Telecom       `json:"telecom" bson:"telecom"`
	MaritalStatus   Concept         `json:"maritalStatus" bson:"marital_status"`
	Communication   []Communication `json:"communication" bson:"communication"`
	Extension       []Extension     `json:"extension" bson:"extension"`
	LastAccess 		time.Time       `json:"-" bson:"last_access"`
}

type PatientBundle struct {
	SearchResult
	Entry []struct {
		FullURL 		string 			`json:"fullUrl" bson:"full_url"`
		Resource struct {
			Patient
		}  `json:"resource"`
	} `json:"entry"`
}
