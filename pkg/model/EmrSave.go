package model

import (
	log "github.com/sirupsen/logrus"
)

type EmrDocuments struct {
	DocumentIds []string `json:"document_ids"`
	EmrMRN      string   `json:"emr_mrn"`
	UserId      string   `json:"user_id"`
}

type EmrPatient struct {
	FhirPatientID string `json:"fhir_patient_id"`
	EmrMRN        string `json:"emr_mrn"`
	UserId        string `json:"user_id"`
}

func AddDocumentsToEMR(emrDocs *EmrDocuments) (int, error) {
	// mrn := emrDocs.EmrMRN
	// userId := emrDocs.UserId
	for _, docId := range emrDocs.DocumentIds {
		_, err := FhirDocumentById((docId))
		if err != nil {
			log.Errorf("Error getting DocumentId: %s", docId)
		}
		//If useId is numeric it is a chartarchive user
	}

	return 0, nil
}
