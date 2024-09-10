package model

import (
	//"errors"
	"fmt"

	//"github.com/davecgh/go-spew/spew"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"

	// "log"
	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	//"strings"
	"time"
)

//PatientQuery is a structur to contain all the query items possible to find a patient
type PatientQuery struct {
	ID         string    `schema:"id"`
	MRN        string    `schema:"mrn"`
	SSN        string    `schema:"ssn"`
	Given      string    `schema:"given"`
	Family     string    `schema:"family"`
	BirthDate  time.Time `schema:"birthdate"` // SHould this be a string?
	Email      string    `schema:"email"`
	PatientGPI string    `schema:"patient_gpi"`
	Encounter  string    `schema:"encounter"`
	//queryString string
}

// CaPatientSearch: Searh for all fhir patients matching the query in the PatientFilter, converts them to
// aa slice of CAPatients and returns the slice, numberInPage, PagesHeldInCache, NumberInCache, err
func (pf *PatientFilter) CaPatientSearch() ([]*CAPatient, int, int64, int64, error) {

	fhirC = config.Fhir()
	activePatientFilter = pf
	pf.makeCacheQueryFilter() // sets up queryString, queryMap, and queryFilter

	// pages, err := pf.PatientPagesInCache()
	// if err != nil {
	// 	return nil, 0, 0, 0, err
	// }
	fmt.Printf("\n\n\n\n\n###### Requesting Page: %d\n", pf.Page)
	// Page == 0 indicates to fill the cache with patients from FHIR Source
	// Page > 0 indicates to use the cache as source of patients
	// if pf.Page > 0 &&  pf.Cache == ""{ //&& pf.Page <= pages {

	// 	fmt.Printf("\n\n\n   ### Reseting Cache. Requesting Cached Page %d\n", pf.Page)
	// 	pats, inPage, pages, inCache, err := pf.GetPatientPage()
	// 	return pats, inPage, pages, inCache, err
	// }

	// if Page == 0 fill the cache querying FHIR using the values in PatientFilter.

	pfs, err := pf.Session.UpdatePatSessionId()
	if err != nil {
		log.Errorf("CaPatientSearch:61 - %s", err.Error())
		return nil, 0, 0, 0, err
	}
	pf.Session = pfs
	//pf.Session = session
	// sessId := CreateSessionId()
	// pf.JWToken = SetTokenSession(pf.JWToken, sessId)

	_, err = pf.FhirCaPatients() // Caches what it finds.
	if err != nil {
		log.Errorf("")
	}

	//fmt.Printf("Search Patients: %s\n", spew.Sdump(pats))
	//pf.CountCachedCaPatients()
	pats, inPage, pages, inCache, err := pf.GetPatientPage()
	return pats, inPage, pages, inCache, err
}

func (pf *PatientFilter) Search() ([]*fhir.Patient, error) {
	fmt.Printf("In Patient Search string: %s\n", pf.queryString)
	// spew.Dump(f)
	// println()
	fhirC = config.Fhir()
	//fmt.Printf("FhirConfig: %s\n", spew.Sdump(fhirC))
	//f.UseCache = "true"

	activePatientFilter = pf
	if pf.Page > 1 {
		//Use the cache  if one repeat the for updates
		fmt.Printf("\n\n\n   ### Requesting Page %d\n", pf.Page)
		return pf.QueryCache()
	}

	pf.makeCacheQueryFilter() // sets up queryString, queryMap, and queryFilter

	pats, err := pf.FhirPatients()
	//fmt.Printf("Search Patients: %s\n", spew.Sdump(pats))
	return pats, err
}

func (f *PatientFilter) FhirCaPatients() ([]*CAPatient, error) {
	fmt.Printf("Filter.FhirPatients using Cache: %s  query: %s\n", f.UseCache, f.queryString)
	startTime := time.Now()
	fhirPatientsResults, err := fhirC.FindFhirPatients(f.queryString)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elapsed time to query FHIR: %f\n", time.Since(startTime).Seconds())

	//CaPatientsFromFhirResults(fhirPatientResults, f)
	caPatients, err := FhirResultsToCAPatients(fhirPatientsResults, f)
	if len(caPatients) == 0 {
		err = fmt.Errorf("404|no patients found for %s", f.queryString)
	}

	log.Debug("Start following next links")
	go f.FollowNextLinks(fhirPatientsResults.SearchResult.Link)

	// patients := PatientsFromResults(fhirPatientsResults, f.SessionId)
	// if len(patients) == 0 {
	// 	err = fmt.Errorf("404|no patients found for %s", f.queryString)
	// }

	return caPatients, err
}

func (f *PatientFilter) FhirPatients() ([]*fhir.Patient, error) {
	fmt.Printf("Filter.FhirPatients using Cache: %s  query: %s\n", f.UseCache, f.queryString)
	startTime := time.Now()
	fhirPatientsResults, err := fhirC.FindFhirPatients(f.queryString)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elapsed time to query FHIR: %f\n", time.Since(startTime).Seconds())

	//CaPatientsFromFhirResults(fhirPatientResults, f)
	caPatients, err := FhirResultsToCAPatients(fhirPatientsResults, f)
	if len(caPatients) == 0 {
		err = fmt.Errorf("404|no patients found for %s", f.queryString)
	}

	log.Debug("Start following next links")
	go f.FollowNextLinks(fhirPatientsResults.SearchResult.Link)

	patients := PatientsFromResults(fhirPatientsResults, f.SessionId)
	if len(patients) == 0 {
		err = fmt.Errorf("404|no patients found for %s", f.queryString)
	}

	return patients, err
}

func (f *PatientFilter) FindPatients() ([]*fhir.Patient, error) {
	//fmt.Printf("Filter.FindPatients -Always get from FHIR \n")

	pats, err := f.FhirPatients()
	if err != nil {
		return nil, err
	}
	//fmt.Printf("\n\n Finished finding from fhir err: %v\n\n", err)

	return pats, err
}

// func (pq *PatientQuery) Search() (*[]Patient, error) {
// 	// println("In Search")
// 	// spew.Dump(*pq)
// 	// println()

// 	// //queryMap := new(map[string]string)
// 	// queryMap, err := pq.toMap()
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// query := com.StringFromMap(queryMap)

// 	// pq.queryString = query
// 	// if queryMap["encounter"] != "" {
// 	// 	println("calling pq.byEncounter")
// 	// 	patient, err := pq.byEncounter()
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// patients := new([]Patient)
// 	// 	*patients = append(*patients, *patient)
// 	// 	return patients, err
// 	// }

// 	// fmt.Printf("Calling FindPatients with query: %s\n", query)
// 	// patients, err := FindPatients(query)

// 	// filter, _ := com.FilterFromMap(queryMap)
// 	// fmt.Printf("Filter: %v\n", filter)
// 	// patients, err := GetCachedPatients(filter)
// 	// if err == nil {
// 	// 	log.Printf("Cached Patient Error; %v\n", err)
// 	// 	log.Println("Cached Patients found: ")
// 	// 	spew.Dump(patients)
// 	// 	return patients, err
// 	// }
// 	// log.Println("No Cached Patients, get from FHIR error")
// 	// q := com.StringFromMap(queryMap)
// 	// fmt.Printf("SearchQuery: %s\n", q)
// 	// patients, err = GetPatients(q)
// 	// spew.Dump(patients)
// 	// //err := fmt.Errorf("First test nothing returned %s", "")
// 	err := fmt.Errorf("Not Implemented")
// 	return nil, err
// }

// func (pq *PatientQuery) toMap() (map[string]string, error) {
// 	m := make(map[string]string)
// 	mrn := strings.Trim(pq.MRN, " ")
// 	given := strings.Trim(pq.Given, " ")
// 	family := strings.Trim(pq.Family, " ")
// 	encounter := strings.Trim(pq.Encounter, " ")
// 	enterpriseID := strings.Trim(pq.EnterpriseID, " ")
// 	email := strings.Trim(pq.Email, " ")
// 	id := strings.Trim(pq.ID, " ")

// 	if given != "" {
// 		if family != "" {
// 			m["given"] = given
// 		} else {
// 			return m, fmt.Errorf("400|Invalid search: given alone is invalid")
// 		}
// 	}
// 	if family != "" {
// 		m["family"] = family
// 	}
// 	if mrn != "" {
// 		m["mrn"] = mrn
// 	}
// 	if encounter != "" {
// 		m["encounter"] = encounter
// 	}
// 	if enterpriseID != "" {
// 		m["enterpriseid"] = enterpriseID
// 	}
// 	if id != "" {
// 		m["enterpriseid"] = id
// 	}
// 	if email != "" {
// 		m["email"] = email
// 	}
// 	// Need Birthdate
// 	return m, nil
// }

// func (pq *PatientQuery) byEncounter() (*Patient, error) {
// 	fmt.Printf("pq:Find Patient by encounter query: %s\n", pq.queryString)
// 	var encounter = new(Encounter)
// 	encounter.EncounterID = pq.Encounter
// 	err := encounter.ForEncounterID()
// 	if err != nil {
// 		return nil, err
// 	}
// 	patID := encounter.PatientID
// 	patient := new(Patient)
// 	patient.EnterpriseID = patID
// 	err = patient.ForID(true)

// 	return patient, err
// }
