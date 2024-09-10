package fhirongo

import (
	"encoding/json"
	"fmt"

	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	//"os"
	"time"
)


func (c *Connection) GetDiagnosticPDF(url string) ([]byte, error) {
	bytes, err := c.GetFhir(url )
	if err != nil {
		return nil, fmt.Errorf("error getting a single diagnistic document: %s", err.Error())
	}
	// if err := os.WriteFile("./debbie.pdf", bytes, 0666); err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Printf("Wrote the pdf file\n")
	return bytes, nil
}
	

// GetDiagnosticReport will return a diagnostic report for a patient with id pid
func (c *Connection) FindDiagnosticReports(query string) (*DocumentResults, error) {
	//query := fmt.Sprintf("patient=%s", patId)
	//fmt.Printf("\n\nFindDiagnosticReport:29 : DiagnosticReport?%s\n", query)
	res, err := c.Query(fmt.Sprintf("DiagnosticReport?%s", query))
	//fmt.Printf("GetDiag error: %v\n", err)
	if err != nil {
		msg := fmt.Sprintf("fhir:35 FindDiagnosticReports error: %s\n", err.Error())
		log.Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	// fmt.Printf("\n\n\n@@@ RAW DiagnosticReport: %s\n\n\n", pretty.Pretty(b))
	// spew.Dump(b)
	data := DocumentResults{}
	if err := json.Unmarshal(res, &data); err != nil {
		msg := fmt.Sprintf("FindDR Unmarshal: error: %s\n", err.Error())
		log.Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	return &data, nil
}

// GetDiagnosticReport will return a diagnostic report for a patient with id pid
func (c *Connection) GetDiagnosticReports(qry string) (*DocumentResults, error) {
	//res, err := c.Query(fmt.Sprintf("DiagnosticReport?patient=%v", pid))

//bytes, err := c.Query("https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/DiagnosticReport?patient=12724066")
	bytes, err := c.Query(fmt.Sprintf("DiagnosticReport%s", qry))
	//fmt.Printf("GetDiag error: %v\n", err)
	if err != nil {
		log.Infof("c.Query returned:51 error: %s\n", err.Error())
		return nil, err
	}
	// fmt.Printf("\n\n\n@@@ RAW DiagnosticReport: %s\n\n\n", pretty.Pretty(b))
	// spew.Dump(b)


	data := DocumentResults{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Infof("unmarshal:56 error: %s\n", err.Error())
		return nil, err
	}
	if len(data.Link) > 1 {  // There are additional pages.
			fmt.Printf("there are more pages\n")
	}
	//fmt.Printf("results: %s\n", spew.Sdump(data))
	return &data, nil
}

// GetPatientDiagnosticReport will return a diagnostic report for a patient with id pid
func (c *Connection) GetPatientDiagnosticReports(pid string) (*DocumentResults, error) {
	qry := fmt.Sprintf("?patient=%s", pid)
	return c.GetDiagnosticReports(qry)

}



func (c *Connection) NextFhirDiagRepts(url string) (*DocumentResults, error) {
	//fmt.Printf("Next retrieving : %s\n", url)
	bytes, err := c.GetFhir(url)
	if err != nil {
		msg := fmt.Sprintf("NextPatient returned error: %s", err.Error())
		log.Errorf("%s", msg)
		return nil, err
	}

	data := DocumentResults{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// DiagnosticReport is a FHIR report


type DiagnosticReport struct {
	CacheID 		  primitive.ObjectID `json:"cach_id" bson:"_id,omitempty"`
	ResourceType      string             `json:"resourceType" bson:"resource_type"`
	SessionId	  	  string   			 `json:"-" bson:"sessionid"`
	ID                string             `json:"id"`
	FullURL			  string 			 `json:"fullUrl" bson:"full_url"`
	EffectiveDateTime time.Time          `json:"effectiveDateTime" bson:"effectiive_datetime"`
	Meta              MetaData           `json:"meta" bson:"meta"`
	Text              TextData           `json:"text" bson:"text"`
	Status            string             `json:"status" bson:"status"`
	Category          CodeableConcept    `json:"category" bson:"category"`
	Code              CodeableConcept    `json:"code" bson:"code"`
	Subject           Person             `json:"subject" bson:"subject"`
	Encounter         EncounterReference `json:"encounter" bson:"encounter"`
	Issued            time.Time          `json:"issued" bson:"issued"`
	Performer         Person             `json:"performer" bson:"performer"`
	PresentedForm     []Attachment       `json:"presentedForm" bson:"presented_form"`
	Request           []Thing            `json:"request" bson:"request"`
	Result            []Thing            `json:"result" bson:"result"`
} 


// type DocumentResponse struct {
// 	//Bundle Bundle
// 	SearchResult
// 	Entry  []struct {
// 		FullURL  string `json:"fullUrl"`
// 		Document  Document `json:"resource"`
// 	} `json:"entry"`
// }
// type DiagnosticReportResponse struct {
// 	//Bundle Bundle
// 	SearchResult
// 	Entry  []struct {
// 		FullURL  string `json:"fullUrl"`
// 		DiagnosticReport  Document `json:"resource"`
// 		//DiagnosticReport  `json:"resource"`
// 	} `json:"entry"`
// }


// type ResourcePartial struct {
// 	ResourceType      string             `json:"resourceType"`
// 	EffectiveDateTime time.Time          `json:"effectiveDateTime"`
// 	RecordedDate      time.Time          `json:"recordedDate"`
// 	Status            string             `json:"status"`
// 	ID                string             `json:"id"`
// 	Subject           Person             `json:"subject"`
// 	Patient           Person             `json:"patient"`
// 	Performer         Person             `json:"performer"`
// 	Recorder          Person             `json:"recorder"`
// 	Encounter         EncounterReference `json:"Encounter"`
// }
// type ResourceEntry struct {
// 	EntryPartial EntryPartial
// 	Resource     DiagnosticReportResource `json:"resource"`
// }

// type DiagnosticReportResource struct {
// 	ResourcePartial   ResourcePartial
// 	EffectiveDateTime time.Time    `json:"effective_date_time"`
// 	Issued            time.Time    `json:"issued"`
// 	Identifier        []Identifier `json:"identifier"`
// 	Meta              MetaData     `json:"meta"`
// 	Text              TextData     `json:"text"`
// 	Category          CodeText     `json:"category"`
// 	Code              CodeText     `json:"code"`
// 	PresentedForm     []Attachment `json:"presentedForm"`
// 	Request           []Thing      `json:"request"`
// 	Encounter         Encounter    `json:"encounter"`
// 	Result            []Thing      `json:"result"`
// }

// // Return the actual decoded text attachment
// func (a *Attachment) DecodeImage() (string, error) {
// 	switch a.ContentType {
// 	case "text/html":
// 		data, err := decodeURL(a.URL)
// 		// ...
// 	case "application/pdf":
// 		data, err := decodeURL(a.URL)
// 		// ...
// 	}
// 	return data, err
// }

// func decodeURL(url string, filePath string) (string, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	out, err := os.Create(filePath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer out.Close()

// }
