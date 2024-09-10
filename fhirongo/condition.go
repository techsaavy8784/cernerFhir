package fhirongo

import (
	"encoding/json"
	"fmt"
)

// GetCondition will return a condition for a patient with id pid
func (c *Connection) GetCondition(pid string) (*Condition, error) {
	body, err := c.Query(fmt.Sprintf("Condition?patient=%v", pid))
	if err != nil {
		return nil, err
	}
	data := Condition{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Condition is a FHIR condition
type Condition struct {
	SearchResult
	Entry []struct {
		EntryPartial
		Resource struct {
			ResourcePartial
			ClinicalStatus     string  `json:"clinicalStatus" bson:"clinical_status"`
			OnsetDateTime      string  `json:"onsetDateTime" bson:"onset_date_time"`
			VerificationStatus string  `json:"verificationStatus" bson:"verification_status"`
			Asserter           Person  `json:"asserter"`
			Code               Concept `json:"code"`
			Category           Concept `json:"category"`
			Severity           Note    `json:"severity"`
		} `json:"resource"`
	} `json:"entry"`
}
