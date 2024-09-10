package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetProcedure will return a careplan for a patient with id pid
func (c *Connection) GetProcedure(pid string) (*Procedure, error) {
	bytes, err := c.Query(fmt.Sprintf("Procedure?patient=%v", pid))
	if err != nil {
		return nil, err
	}
	data := Procedure{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Procedure is a FHIR procedure
type Procedure struct {
	SearchResult
	Entry []struct {
		SearchResult
		Resource struct {
			ResourcePartial
			NotPerformed      bool         `json:"notPerformed" bson:"not_performed"`
			PerformedDateTime time.Time    `json:"performedDateTime" bson:"performed_date_time"`
			Identifier        []Identifier `json:"identifier"`
			Code              Concept      `json:"code"`
		} `json:"resource"`
	} `json:"entry"`
}
