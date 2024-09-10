package fhirongo

import (
	"encoding/json"
	"fmt"
)

// TODO more types of queries

// GetMedication will return a careplan for a patient with id pid
func (c *Connection) GetMedication(pid string) (*Medication, error) {
	res, err := c.Query(fmt.Sprintf("MedicationOrder?patient=%v", pid))
	if err != nil {
		return nil, err
	}
	data := Medication{}
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Medication is a FHIR medication
type Medication struct {
	SearchResult
	Entry []struct {
		EntryPartial
		Resource struct {
			ResourcePartial
			DateWritten         string              `json:"dateWritten" bson:"date_written"`
			Identifier          []Identifier        `json:"identifier"`
			Prescriber          Person              `json:"prescriber"`
			MedicationReference Thing               `json:"medicationReference" bson:"medication_reference"`
			DosageInstruction   []DosageInstruction `json:"dosageInstruction" bson:"dosage_instruction`
			DispenseRequest     DispenseRequest     `json:"dispenseRequest" bson:"dispense_request"`
		} `json:"resource"`
	} `json:"entry"`
}
