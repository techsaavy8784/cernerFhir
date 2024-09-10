package fhirongo

import (
	"encoding/json"
	"fmt"
)

// GetFamilyMemberHistory will return a family history for a patient with id pid
func (c *Connection) GetFamilyMemberHistory(pid string) (*FamilyMemberHistory, error) {
	res, err := c.Query(fmt.Sprintf("FamilyMemberHistory?patient=%v", pid))
	if err != nil {
		return nil, err
	}
	data := FamilyMemberHistory{}
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// FamilyMemberHistory is a FHIR family hx
type FamilyMemberHistory struct {
	SearchResult
	Entry []struct {
		EntryPartial
		Resource struct {
			ResourcePartial
			Date            string     `json:"date"`
			Name            string     `json:"name"`
			DeceasedBoolean bool       `json:"deceasedBoolean" bson:"deceased_boolean"`
			Relationship    Concept    `json:"relationship"`
			Condition       []CodeText `json:"condition"`
		} `json:"resource"`
	} `json:"entry"`
}
