package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetImmunization will return a careplan for a patient with id pid
func (c *Connection) GetImmunization(pid string) (*Immunization, error) {
	res, err := c.Query(fmt.Sprintf("Immunization?patient=%v", pid))
	if err != nil {
		return nil, err
	}
	data := Immunization{}
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Immunization is a FHIR immunization
type Immunization struct {
	SearchResult
	Entry []struct {
		EntryPartial
		Resource struct {
			ResourcePartial
			Date        time.Time `json:"date"`
			WasNotGiven bool      `json:"wasNotGiven" bson:"was_not_given"`
			Reported    bool      `json:"reported"`
			LotNumber   string    `json:"lotNumber" bson:"lot_number"`
			//ID          string    `json:"id"`
			VaccineCode Concept `json:"vaccineCode" bson:"vaccine_code"`
		} `json:"resource"`
	} `json:"entry"`
}
