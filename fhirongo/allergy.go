package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"
)

// GetAllergyIntolerence will return patient allergy intolerence
func (c *Connection) GetAllergyIntolerence(pid string) (*AllergyResults, error) {
	body, err := c.Query(fmt.Sprintf("AllergyIntolerance?patient=%v", pid))
	if err != nil {
		return nil, err
	}
	data := AllergyResults{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

type AllergyResults struct {
	Bundle
	Entry []struct {
		FullURL string             `json:"fullUrl"`
		Allergy AllergyIntolerence `json:"resource"`
	} `json:"entry"`
}

// AllergyIntolerence is a FHIR allergy intolerence
// type AllergyIntolerence struct {
// 	SearchResult
// 	Entry []struct {
// 		EntryPartial
// 		Resource struct {
// 			ResourcePartial
// 			Criticality string     `json:"criticality"`
// 			Onset       time.Time  `json:"onset"`
// 			Substance   Concept    `json:"substance"`
// 			Reaction    []Reaction `json:"reaction"`
// 			Note        Note       `json:"note"`
// 		} `json:"resource"`
// 	} `json:"entry"`
// }

type AllergyIntolerence struct {
	Criticality string     `json:"criticality"`
	Onset       time.Time  `json:"onset"`
	Substance   Concept    `json:"substance"`
	Reaction    []Reaction `json:"reaction"`
	Note        Note       `json:"note"`
}
