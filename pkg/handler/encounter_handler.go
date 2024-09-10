package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"time"
	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	//"github.com/gorilla/mux"
	ca "github.com/dhf0820/cernerFhir/pkg/ca"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	"github.com/gorilla/schema"

	fhir "github.com/dhf0820/cernerFhir/fhirongo"
)

type EncounterResponse struct {
	StatusCode   int               `json:"status_code"`
	Message      string            `json:"message"`
	CacheStatus  string            `json:"cache_status"`
	TotalInCache int64             `json:"total_encounters"`
	PagesInCache int64             `json:"pages_in_cache"`
	NumberInPage int64             `json:"encounters_in_page"`
	Page         int64             `json:"page"`
	SessionId    string            `json:"session_id"`
	Encounters   []*fhir.Encounter `json:"encounters"`
	Encounter    fhir.Encounter    `json:"encounter"`
}

type EncounterWithPatientResponse struct {
	Encounter m.Encounter
	Patient   m.Patient
}

func WriteEncounterResponse(w http.ResponseWriter, resp *EncounterResponse) error {
	w.Header().Set("Content-Type", "application/json")

	// resp.StatusCode = status_code
	// resp.Data = data
	w.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	return nil
}

//############################################ Route Handlers ##########################################
// getEncounter finds the Encounter by their patentID
// findEncounter finds Encounters by information in their record (mrn, ssn.,,,)
// createEncounter
// updateEncounter
// deleteEncounter
//Route Handlers
func SessionEncounters(w http.ResponseWriter, r *http.Request) {
}

func FindEncounters(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("==================FindEncounters========\n")
	//println("EncounterHandler:FindEncounter")
	resp := EncounterResponse{}
	ef, err := SetupEncounterFilter(r)
	if err != nil {
		resp.Message = err.Error()
		resp.StatusCode = 400

		WriteEncounterResponse(w, &resp)
		return
	}

	if ef.Cache == "stats" { // Retrieve from cache only
		log.Debugf("EncounterHandler:143 -- querying stats")
		pagesInCache, totalInCache, cacheStatus, _ := ef.FhirEncounterCacheStats()
		resp.CacheStatus = cacheStatus
		resp.TotalInCache = totalInCache
		resp.PagesInCache = pagesInCache
		resp.Page = ef.Page
		resp.StatusCode = 200
		resp.Message = "Ok"
		WriteEncounterResponse(w, &resp)
		return
	}
	if ef.Page == 0 {
		ef.Page = 1
	}
	// if ef.Cache != "reset" {
	startTime := time.Now()
	encs, numInPage, pagesInCache, totalInCache, cacheStatus, err := ef.GetFhirEncounterPage()
	if err != nil {
		if err.Error() == "notFound" {
			log.Errorf("FindEncounters:170  --  Error: %s", err.Error())
			FillEncounterCache(ef, w)
			log.Infof("GetPage took: %f seconds", time.Since(startTime).Seconds())
			return
		} else {
			log.Errorf("FindEncounters:181 -- Error: %s", err.Error())
			resp.StatusCode = 400
			resp.Message = err.Error()
			WriteEncounterResponse(w, &resp)
			return
		}
	}

	if ef.ResultFormat == "ca" {
		cacheStatus := ef.Session.GetEncounterStatus()
		ca.FhirEncountersToCA(w, totalInCache, pagesInCache, int64(numInPage), ef.Page, cacheStatus, encs)
		log.Infof("GetPageCA took: %f seconds", time.Since(startTime).Seconds())
		return
	}
	resp.CacheStatus = cacheStatus
	resp.Encounters = encs
	resp.NumberInPage = numInPage
	resp.PagesInCache = pagesInCache
	resp.TotalInCache = totalInCache
	if err != nil {
		resp.Message = err.Error()
		resp.StatusCode = 400
	}
	WriteEncounterResponse(w, &resp)
}

func FillEncounterCache(ef *m.EncounterFilter, w http.ResponseWriter) {
	log.Debug("#####     FillEncounterCache")
	startTime := time.Now()
	//ef.Session.UpdateEncSessionId()  // reset force the EncSessionId to be new so the old will delete
	resp := EncounterResponse{}
	_, err := ef.SearchEncounters() // fill the cache
	if err != nil {
		err = WriteGenericResponse(w, 404, err.Error())
		if err != nil {
			log.Println("Error writing response", err)
		}
		return
	}
	encounters, numInPage, pagesInCache, totInCache, cacheStatus, err := ef.GetFhirEncounterPage()

	if err != nil {
		resp.StatusCode = 200
		resp.Message = err.Error()
		WriteEncounterResponse(w, &resp)
		return
	}
	if ef.ResultFormat == "ca" {
		ca.FhirEncountersToCA(w, totInCache, pagesInCache, int64(numInPage), ef.Page, cacheStatus, encounters)
		return
	}

	resp.CacheStatus = cacheStatus
	resp.NumberInPage = numInPage
	resp.PagesInCache = pagesInCache
	resp.TotalInCache = totInCache
	resp.StatusCode = 200
	resp.Message = "Ok"
	resp.Encounters = encounters
	WriteEncounterResponse(w, &resp)
	log.Infof("FillEncounterCache:244 -- ElapsedTime: %f seconds", time.Since(startTime).Seconds())
}

//statusCode = 200

// encounters, err := m.FindFhirEncounters(query)
// if err != nil {
// 	var statusCode int
// 	s := strings.Split(err.Error(), "|")

// 	// if statusCode, err = strconv.ParseInt(s[0], 10, 64); err == nil {
// 	if statusCode, err = strconv.Atoi(s[0]); err == nil {
// 		statusCode = 500
// 	}

// 	err = fmt.Errorf("%v", s[1])
// 	// fmt.Printf("  StatusCode: %d\n", statusCode)
// 	// fmt.Printf("error: %s\n", err)
// 	//statusMessage := fmt.Sprintf("FindEncounter Error: %s", err)
// 	err = WriteGenericResponse(w, statusCode, err.Error())
// 	if err != nil {
// 		log.Println("Error writing response", err)
// 	}
// 	return
// }
// if len(encounters) == 0 {
// 	err := fmt.Errorf("No Encounters found for %s", ef.)
// 	err = WriteGenericResponse(w, 400, err.Error())
// 	if err != nil {
// 		log.Println("Error writing response", err)
// 	}
// 	return
// }

//encounters = *encs
// fmt.Println("Encounters found:")
// spew.Dump(Encounters)

// getEncounter finds the Encounter by their patentID
// findEncounter finds Encounters by information in their record (mrn, ssn.,,,)
// createEncounter
// updateEncounter
// deleteEncounter
//Route Handlers
// func GetEncounter(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	fmt.Printf("EncHandler:GetEncounter Params: %v\n", params)
// 	fmt.Printf("id: %v\n", params["id"])
// 	id := params["id"]

// 	if id == "" {
// 		statusMessage := fmt.Sprintf("EncHandler:GetEncounter - no Encounter id provided")
// 		err := WriteGenericResponse(w, 400, statusMessage)
// 		log.Println(statusMessage)
// 		if err != nil {
// 			log.Println("Error writing response", err)
// 		}
// 		return
// 	}
// 	encounterFilter, err := SetupEncounterFilter(r)

// 	if err != nil {
// 		statusMsg := "401|Unauthorized"
// 		println(statusMsg)
// 		err := fmt.Errorf("%s", statusMsg)
// 		HandleFhirError("GetEncounter", w, err)
// 		return
// 	}

// 	encounterFilter.ID = id

// 	startTime := time.Now()

// 	// Check if cached already
// 	var encounter = new(m.Encounter)

// 	encounter.EncounterID = id
// 	err = encounter.ForPatID() //"true", encounterFilter.SessionId)

// 	//encounter, err = m.GetFhirEncounter(id)
// 	if err != nil {
// 		HandleFhirError("EncHandler:GetEncounter", w, err)
// 		// , w http.ResponseWriter, err error)
// 		// statusMessage := fmt.Sprintf("Encounter using id- %s was not found Code: 404", id)
// 		// err := WriteGenericResponse(w, 400, statusMessage)
// 		// log.Println(statusMessage)
// 		// if err != nil {
// 		// 	log.Println("Error writing response", err)
// 		// }
// 		return
// 	}
// 	// fmt.Printf("GetEncounter:  \n")
// 	// spew.Dump(encounter)

// 	resp := new(EncounterWithPatientResponse)
// 	// patient := new(m.Patient)
// 	// patient.EnterpriseID = encounter.PatientID
// 	patientFilter := m.PatientFilter{}
// 	patientFilter.PatientGPI = encounter.PatientID
// 	patientFilter.SessionId = encounterFilter.SessionId
// 	patientFilter.UseCache = encounterFilter.UseCache
// 	//patient, err := patientFilter.ForEnterpriseID()
// 	//elapsedTime := time.Since(startTime)
// 	//log.Printf("   @@@ Get Encounter took %s ms\n", elapsedTime)
// 	// if err == nil {
// 	// 	resp.Patient = *patient
// 	// }
// 	resp.Encounter = *encounter
// 	elapsedTime := time.Since(startTime)
// 	log.Printf("   @@@ Get Encounter took %s ms\n", elapsedTime)
// 	WriteEncounterPatientResponse(w, 200, resp)
// }

// func FindPatientEncounters(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	fmt.Printf("GetPatientEncounters Params: %v\n", params)
// 	fmt.Printf("id: %v\n", params["id"])
// 	id := params["id"]
// 	if id == "" {
// 		statusMessage := fmt.Sprintf("FindPatientEncounters - no PatientID provided")
// 		err := WriteGenericResponse(w, 400, statusMessage)
// 		log.Println(statusMessage)
// 		if err != nil {
// 			log.Println("Error writing response", err)
// 		}
// 		return
// 	}
// 	startTime := time.Now()
// 	// Check if cached already
// 	//var encounter = new(m.Encounter)

// 	encounters, err := m.EncountersForPatientID(id)

// 	//encounter, err = m.GetFhirEncounter(id)
// 	if err != nil {
// 		HandleFhirError("GetEncounter", w, err)
// 		return
// 	}
// 	// fmt.Printf("GetEncounter:  \n")
// 	// spew.Dump(encounter)

// 	elapsedTime := time.Since(startTime)
// 	log.Printf("   @@@ Get Encounter took %s ms\n", elapsedTime)
// 	WriteEncountersResponse(w, 200, encounters)
// }

func SetupEncounterFilter(r *http.Request) (*m.EncounterFilter, error) {
	fmt.Printf("Raw Query:390 -- %s\n", r.URL.RawQuery)
	config := m.ActiveConfig()
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var encounterFilter m.EncounterFilter
	err := decoder.Decode(&encounterFilter, r.URL.Query())
	if err != nil {
		log.Println("Error in GET parameters : ", err)
	} else {
		//log.Println("GET parameters : ", spew.Sprint(encounterFilter))
	}

	session, err := m.ValidateSession(r.Header.Get("SESSION"))
	if err != nil {
		log.Errorf("setupEncounterFilter:308 - NoSession - %s", err.Error())
		//err = fmt.Errorf("Authorization was not found")
		return nil, errors.New("please log in")
	}
	encounterFilter.Session = session
	//fmt.Printf("Session:327  session : %s\n", spew.Sdump(session))
	// if as == nil {
	// 	log.Errorf("setupEncounterFilter:308 - NoSession - %s", err.Error())
	// 	//err = fmt.Errorf("Authorization was not found")
	// 	return nil, err
	// }
	encounterFilter.SessionId = session.EncSessionId
	//UseCache := r.Header.Get("UseCache")
	count := r.Header.Get("Count")
	if count == "" {
		log.Debugf("Setting Default Count: %s", config.RecordLimit())
		encounterFilter.Count = config.RecordLimit()
	} else {
		encounterFilter.Count = count
	}
	resultFormat := r.Header.Get("RESULTFORMAT")
	if resultFormat == "" {
		resultFormat = "ca"
	}
	encounterFilter.ResultFormat = resultFormat
	//fmt.Printf("header mode: %s\n", encounterFilter.Mode)
	//fmt.Printf("Encounter HeaderMode: %s\n", resultFormat)

	if encounterFilter.Count == "" {
		encounterFilter.Count = count
	}
	// if encounterFilter.EncounterID != "" {
	// 	encounterFilter.EncounterID = encounterFilter.EncounterID
	// 	encounterFilter.EncounterID = ""
	// }
	return &encounterFilter, nil
}
