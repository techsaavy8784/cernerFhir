package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetObservation will return a careplan for a patient with id pid
func (c *Connection) GetObservation(pid string, code string) (*Observation, error) {
	bytes, err := c.Query(fmt.Sprintf("Observation?patient=%v&code=%v", pid, code))
	if err != nil {
		return nil, err
	}
	data := Observation{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Observation is a FHIR observation
type Observation struct {
	SearchResult
	Entry []struct {
		EntryPartial
		Resource struct {
			ResourceType      string    `json:"resourceType" bson:"resource_type"`
			EffectiveDateTime time.Time `json:"effectiveDateTime" bson:"effective_date_time"`
			Status            string    `json:"status"`
			ID                string    `json:"id"`
			Code              Concept   `json:"code"`
			ValueQuantity     Quantity  `json:"valueQuantity" bson:"value_quantity"`
			Subject           Thing     `json:"subject"`
			Performer         []Person  `json:"performer"`
			Category          Concept   `json:"category"`
		} `json:"resource"`
	} `json:"entry"`
}
