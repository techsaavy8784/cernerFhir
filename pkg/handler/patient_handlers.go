package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	//"strconv"
	//"time"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"

	//"github.com/dhf0820/cernerFhir/pkg/common"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	//fhir "github.com/dhf0820/cernerFhir/fhirongo"
)

//"gopkg.in/mgo.v2/bson"

//"strconv"
//"strings"
//"time"
//"github.com/davecgh/go-spew/spew"
//"github.com/oleiade/reflections"

type PatientResponse struct {
	Code     int             `json:"code"`
	Message  string          `json:"status"`
	Patients []*fhir.Patient `json:"patients"`
	Patient  *fhir.Patient   `json:"patient"`
}

type PatientEmrAddResponse struct {
	StatusCode int    `json:"code"`
	Message    string `json:"status"`
}
type EmrPatient struct {
	MRN       string `json:"mrn"`        // callers medical record number
	PatientId string `json:"patient_id"` // Fhir Patient Id
	UserId    string `json:"user_id"`    // User ID Associated with the EMR
}

type PatientRawCaResponse struct {
	Patients []*m.CAPatient `json:"patients"`
	Patient  *m.CAPatient   `json:"patient"`
}
type PatientCaResponse struct {
	StatusCode   int            `json:"status_code"`
	Message      string         `json:"message"`
	CacheStatus  string         `json:"cache_status"`
	Total        int64          `json:"total_documents"`
	PagesInCache int64          `json:"pages_in_cache"`
	NumberInPage int            `json:"num_in_page"`
	Page         int            `json:"page"`
	SessionId    string         `json:"session_id"`
	Token        string         `json:"token"`
	Patients     []*m.CAPatient `json:"patients"`
	Patient      *m.CAPatient   `json:"patient"`
}
type PatientSummaryResponse struct {
	Code     int          `json:"code"`
	Message  string       `json:"status"`
	Patients []*m.Patient `json:"patients"`
	Patient  *m.Patient   `json:"patient"`
}

type PatientEncounterResponse struct {
	// StatusCode uint64            `json:"status_code"`
	// Status     string            `json:"status"`
	Data []*m.Encounter `json:"data"`
}

func WritePatientResponse(w http.ResponseWriter, statusCode int, resp *PatientResponse) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err

	}
	return nil
}

func WritePatientCaResponse(w http.ResponseWriter, resp *PatientCaResponse) error {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("WritePatientCaResponse:80 -- statusCode: %d\n", resp.StatusCode)
	//w.WriteHeader()
	w.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err

	}
	return nil
}

func WritePatientAddEmrResponse(w http.ResponseWriter, resp *PatientEmrAddResponse) error {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("WritePatientCaResponse:80 -- statusCode: %d\n", resp.StatusCode)
	//w.WriteHeader()
	w.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err

	}
	return nil
}

func WritePatientSummaryResponse(w http.ResponseWriter, statusCode int, resp *PatientSummaryResponse) error {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("WriteResponse statusCode: %d\n", statusCode)

	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err

	}
	return nil
}

func WritePatientRawCaResponse(w http.ResponseWriter, statusCode int, rawResp *PatientRawCaResponse) error {
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("WriterResponse statusCode: %d\n", statusCode)

	var err error
	w.WriteHeader(statusCode)
	if rawResp.Patient != nil {
		err = json.NewEncoder(w).Encode(rawResp.Patient)
	} else {
		err = json.NewEncoder(w).Encode(rawResp.Patients)
	}
	if err != nil {
		fmt.Println("Error marshaling Raw Patient JSON:", err)
		return err

	}
	return nil
}

func WritePatientEncountersResponse(w http.ResponseWriter, statusCode int, data []*m.Encounter) error {
	w.Header().Set("Content-Type", "application/json")
	var resp PatientEncounterResponse
	//resp.StatusCode = uint64(status_code)
	resp.Data = data
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err

	}
	return nil
}

// func logFatal(err error) {
// 	if err != nil {
// 		fmt.Printf("LogFatal!!!   ")
// 		log.Fatal(err)
// 	}
// }

//############################################ Route Handlers ##########################################
// getPatient finds the Patient by their patentID
// findPatient finds Patients by information in their record (mrn, ssn.,,,)
// createPatient
// updatePatient
// deletePatient
//Route Handlers

func SessionPatients(w http.ResponseWriter, r *http.Request) {
}

func SearchCaPatient(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nPatientHandler:SearchCaPatient:162\n\n")
	fmt.Printf("\n\n\n###Raw patient raw query: %s\n\n", r.URL.RawQuery)
	caResp := PatientCaResponse{}
	pf, err := SetupPatientFilter(r)
	if err != nil {
		msg := fmt.Sprintf("SearchCaPatient:166 Error: %s", err.Error())
		log.Errorf("%s", msg)
		caResp.StatusCode = 400
		caResp.Message = msg
		WritePatientCaResponse(w, &caResp)
		return
	}
	pf.ResultFormat = r.Header.Get("ResultFormat")
	if pf.ResultFormat == "" {
		pf.ResultFormat = r.Header.Get("RESULTFORMAT")
	}
	if pf.ResultFormat == "" {
		pf.ResultFormat = "ca"
	}
	fmt.Printf("SearchCaPatient:195 -- filter: %s\n", spew.Sdump(pf))

	//TODO: Handle the Authorization
	fmt.Printf("Mode: %s,  ResultFormat: %s, RESULTFORMAT: %s\n",
		pf.Mode, r.Header.Get("ResultFormat"), r.Header.Get("RESULTFORMAT"))
	if pf.Mode == "ca" || r.Header.Get("ResultFormat") == "ca" || r.Header.Get("ResultFormat") == "" {
		if pf.Cache == "stats" { // Retrieve from cache only
			log.Debugf("PatientHandler:202 -- querying stats: Page: %d", pf.Page)
			pagesInCache, totalInCache, _ := pf.PatientCacheStats()
			//SetTokenCookie(w, pf.JWTokenStr)
			//tc, err := UpdateTokenCookie(pf.TokenCookie)

			//caResp.SessionId = pf.SessionId
			//caResp.NumberInPage = numberInPage // len(patients)
			total := totalInCache //pf.CountCachedCaPatients()
			caResp.Total = total
			caResp.PagesInCache = pagesInCache
			caResp.Page = int(pf.Page)
			caResp.StatusCode = 200
			caResp.Message = "Ok"
			WritePatientCaResponse(w, &caResp)
			return
		}
		if pf.Page == 0 {
			pf.Page = 1
		}
		if pf.Cache != "reset" { // Reset the cache based upon the cullent filters
			fmt.Printf("\n\n\n   ### Requesting cached page %d\n", pf.Page)
			caPats, numberInPage, pagesInCache, totalInCache, err := pf.GetPatientPage()
			if err != nil {
				caResp.Message = err.Error()
				caResp.StatusCode = 400
			} else {
				caResp.Message = "Ok"
				caResp.StatusCode = 200
			}
			caResp.Patients = caPats
			caResp.Page = int(pf.Page)
			caResp.NumberInPage = numberInPage
			caResp.PagesInCache = pagesInCache
			caResp.Total = totalInCache
			WritePatientCaResponse(w, &caResp)
			return
		}
		log.Debugf("SearchCaPatients:223 - calling CaPatientSearch to fill the cache")
		pfs, err := pf.Session.UpdatePatSessionId()
		if err != nil {
			log.Errorf("CaPatientSearch: UpdatePatSessionId: 61 - %s", err.Error())
			caResp.StatusCode = 400
			caResp.Message = err.Error()
			WritePatientCaResponse(w, &caResp)
			return
		}
		pf.Session = pfs
		caPats, numberInPage, pagesInCache, numberInCache, err := pf.CaPatientSearch()
		//patients, err := pf.CaSearch()
		if err != nil {
			log.Errorf("SearchCaPatients:225: - returned err: %s", err.Error())
			caResp.StatusCode = 400
			caResp.Message = err.Error()
			err = WritePatientCaResponse(w, &caResp)
			if err != nil {
				log.Println("Error writing Patientresponse:230", err)
			}
			return
		}
		caResp.SessionId = pf.Session.PatSessionId
		caResp.NumberInPage = numberInPage // len(patients)
		total := numberInCache             //pf.CountCachedCaPatients()
		caResp.Total = total
		caResp.PagesInCache = pagesInCache //int(m.LinesPerPage())  //TODO: Lines in page should be lines in current page
		caResp.Page = int(pf.Page)
		if caResp.Page == 0 {
			caResp.Page = 1
		}
		fmt.Printf("SearchCaPatient:254 partial returning: %s\n", spew.Sdump(caResp))
		// pgSize, err := strconv.Atoi(cfg.Env("page_size"))
		// if err!= nil {
		// 	pgSize = 20
		// }

		// pages, _ := common.CalcPages(int(total),pgSize )
		// caResp.Pages = pages
		// caResp.Page = int(pf.Page)
		caResp.Patients = caPats
		caResp.StatusCode = 200
		caResp.Message = "Ok"
		WritePatientCaResponse(w, &caResp)
		return
	}

	fhirPatients, err := pf.Search()
	if err != nil {
		caResp.StatusCode = 400
		caResp.Message = err.Error()
		err = WritePatientCaResponse(w, &caResp)
		if err != nil {
			log.Println("Error writing Patientresponse:265", err)
		}
		return
	}

	switch r.Header.Get("ResultFormat") {
	// case "ca-3":
	// 	caResp := PatientCaResponse{}
	// 	caResp.Patients = patients
	// 	caResp.Code = 200
	// 	WritePatientCaResponse(w, caResp.Code, &caResp)
	case "summary":
		sumResp := PatientSummaryResponse{}
		sumResp.Patients = m.FhirPatientsToSum(fhirPatients, pf.Session.PatSessionId)
		sumResp.Code = 200
		WritePatientSummaryResponse(w, sumResp.Code, &sumResp)
	case "fhir":
		resp := PatientResponse{}
		resp.Patients = fhirPatients
		resp.Code = 200
		resp.Message = "Ok"
		WritePatientResponse(w, resp.Code, &resp)
	default:
		resp := PatientResponse{}
		resp.Patients = fhirPatients
		resp.Code = 200
		resp.Message = "Ok"
		WritePatientResponse(w, resp.Code, &resp)
	}
	// resp.Code = 200
	// resp.Message = "Ok"
	// resp.Patients = patients
	// WritePatientResponse(w, resp.Code, &resp )
}

// getPatient finds the Patient by their patentID
// findPatient finds Patients by information in their record (mrn, ssn.,,,)
// createPatient
// updatePatient
// deletePatient
//Route Handlers

func GetPatient(w http.ResponseWriter, r *http.Request) {
	resp := PatientResponse{}
	// patientFilter, err := SetupPatientFilter(r)

	// if err != nil {fmt.Printf("")
	// 	caResp.Code = 400
	// 	caResp.Message = err.Error()
	// 	WriteCaPatientResponse(w, caResp.Code, &caResp)
	// 	return
	// }

	pf, err := SetupPatientFilter(r)
	if err != nil {
		msg := fmt.Sprintf("SetupPatientFilter Error: %s", err.Error())
		log.Errorf("%s", msg)
		resp.Code = 400
		resp.Message = msg
		WritePatientResponse(w, resp.Code, &resp)
		return
	}
	params := mux.Vars(r)
	fmt.Printf("\n\n\n\n\n#####HandleGetPatient Params: %v\n", params)
	fmt.Printf("id: %v\n", params["id"])
	gpi := params["id"]
	if gpi == "" {
		resp.Message = "HandleGetPatient - no patient patient_gpi provided"
		fmt.Printf("%s\n", resp.Message)
		resp.Code = 400
		WritePatientResponse(w, resp.Code, &resp)
		return
	}

	//patientFilter.EnterpriseID = id

	patient, err := m.ForPatientGPI(gpi)
	if err != nil {
		resp := PatientResponse{}
		resp.Code = 400
		resp.Message = err.Error()
		log.Errorf("GetPatient: rforEnterpriseID returned %s", resp.Message)
		WritePatientResponse(w, resp.Code, &resp)
		return
	}
	fmt.Printf("GetPatient:34 -- Results: %s\n", spew.Sdump(patient))
	fmt.Printf("Mode: %s,  ResultFormat: %s, RESULTFORMAT: %s\n",
		pf.Mode, r.Header.Get("ResultFormat"), r.Header.Get("RESULTFORMAT"))
	// if pf.Mode == "ca" || r.Header.Get("ResultFormat") == "ca" || r.Header.Get("ResultFormat") == "" {

	// 	caResp := PatientRawCaResponse{}
	// 	caResp.Patient = nil //model.FhirPatientToCA(patient)
	// 	WritePatientRawCaResponse(w, 200, &caResp)
	// 	return
	// }
	switch r.Header.Get("ResultFormat") {
	case "ca":
		caResp := PatientCaResponse{}
		caPat, err := m.FhirPatientToCA(patient, pf)
		if err != nil {
			caResp.StatusCode = 400
			caResp.Message = err.Error()
		} else {
			caResp.StatusCode = 200
			caResp.Message = "Ok"
			caResp.Patient = caPat
		}
		WritePatientCaResponse(w, &caResp)
	case "summary":
		sumResp := PatientSummaryResponse{}
		sumResp.Patient = m.FhirPatientToSum(patient, pf.Session.PatSessionId)
		sumResp.Code = 200
		WritePatientSummaryResponse(w, sumResp.Code, &sumResp)
	default:
		resp := PatientResponse{}
		resp.Patient = patient
		resp.Code = 200
		resp.Message = "Ok"
		WritePatientResponse(w, resp.Code, &resp)
	}
}

//func HandleCa()
// func PatientFhirResult(resp *resp) (*fhir.Patient, error) {
// 	b, _ := ioutil.ReadAll(resp.Body)
// 	data := []*m.Patient{}
// 	if err := json.Unmarshal(b, &data); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return nil, err
// 	}
// 	return data, nil
// }
// func GetPatientEncounters(w http.ResponseWriter, r *http.Request) {
// 	var encounters []m.Encounter
// 	queryMap := r.URL.Query()
// 	q := queryMap

// 	if val, ok := q["patient"]; ok {
// 		patientID := val[0]
// 		//fmt.Printf("Look for mrn: %s\n", mrn)
// 		encs, err := m.FindFhirPatientEncountersByPatientId(patientID)
// 		if err != nil {
// 			err := WriteGenericResponse(w, 404, fmt.Sprintf("FindPatientEncounters : %s", err))
// 			if err != nil {
// 				log.Println("Error writing response", err)
// 			}
// 			return
// 		}
// 		encounters = *encs
// 		// fmt.Println("patients found:")
// 		// spew.Dump(patients)
// 	} else {
// 		//fmt.Printf("Look for patient Encounters: type: %T  val: %v\n", q["id"], q)

// 		// query := ""
// 		// for k := range q {
// 		// 	fmt.Printf("FindPatientEncounters k: %s  v: %s\n", k, q[k][0])
// 		// 	if query == "" {
// 		// 		query = fmt.Sprintf("%s=%s", k, q[k][0])
// 		// 	} else {
// 		// 		query = query + "&" + fmt.Sprintf("%s=%s", k, q[k][0])
// 		// 	}
// 		// }
// 		query := makeQueryString(q)
// 		fmt.Printf("Looking for %s\n", query)
// 		encs, err := m.FindFhirPatientEncountersByPatientId(query)
// 		statusMessage := fmt.Sprintf("FindPatient Error: %s", err)
// 		err = WriteGenericResponse(w, 404, statusMessage)
// 		//log.Println(statusMessage)
// 		if err != nil {
// 			log.Println("Error writing response", err)
// 		}
// 		fmt.Printf("Found %d Encounters\n", len(*encs))
// 		//spew.Dump(encs)
// 		encounters = *encs
// 		// fmt.Println("patients found:")
// 		// spew.Dump(patients)
// 	}

// 	WritePatientEncountersResponse(w, 200, encounters)
// }

func AddPatientEMR(w http.ResponseWriter, r *http.Request) {
	resp := PatientEmrAddResponse{}
	resp.StatusCode = 200
	resp.Message = "OK"
	WritePatientAddEmrResponse(w, &resp)
	return

}

// func UpdateTokenCookie(cookie *http.Cookie) (*http.Cookie, error) {
// 	token, err := m.VerifyTokenString(cookie.Value)
// 	if err != nil {
// 		return nil, fmt.Errorf("TokenCookie:430 - %s", err.Error())
// 	}

// 	updToken := m.UpdateTokenExpire(token)
// 	tokenStr, err := m.TokenSignedString(updToken)
// 	tokenCookie := &http.Cookie {
// 		Name:"token",
// 		Value:tokenStr,
// 		MaxAge: int(60*m.LoginExpiresAfter()),
// 		HttpOnly: true,
// 	}
// 	return tokenCookie, nil
// }

func SetupPatientFilter(r *http.Request) (*m.PatientFilter, error) {
	config := m.ActiveConfig()
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var patientFilter m.PatientFilter
	err := decoder.Decode(&patientFilter, r.URL.Query())
	if err != nil {
		log.Errorf("Decode Params:469 -- error: %s ", err.Error())
		return nil, fmt.Errorf("DecodeParams:469 -- error: %s", err.Error())
	}

	sessionId := r.Header.Get("SESSION")
	if sessionId == "" {
		log.Errorf("No SESSION Header")
		return nil, errors.New("no SESSION Header")
	}
	log.Debugf("Got Session Header: %s", sessionId)
	as, err := m.ValidateSession(sessionId)
	if err != nil {
		log.Errorf("NoSession:483 - %s", err.Error())
		return nil, errors.New("please log in")
	}
	log.Debugf("SetupPatientFilter:486 - session %s is valid", sessionId)
	patientFilter.Session = as
	// tokenCookie, err := r.Cookie("token")
	// if err != nil{
	// 	msg := "no cookie exists:473, please log in"
	// 	log.Error(msg)
	// 	return nil, errors.New(msg)
	// }
	// log.Debugf("444: Token MaxAge: %d",tokenCookie.MaxAge)
	// log.Debugf("445: tokenCookie expires at: %s", tokenCookie.Expires.String())
	// // if tokenCookie.Expires.Before(time.Now()) {   // cookie has expired. Request login
	// // 	log.Debugf("Expired token: %s", spew.Sdump(tokenCookie))
	// // 	//return nil, errors.New("session Expired, please login again")
	// // }
	// tokenStr := tokenCookie.Value
	// token, err :=m.VerifyTokenString(tokenStr)
	// if err != nil {
	// 	return nil, err
	// }
	// // auth := r.Header.Get("AUTHORIZATION")
	// // log.Infof("Authorization is: [%s]", auth)
	// // tokenStr := m.ExtractToken(r)
	// // token, err := m.VerifyTokenString(tokenStr)
	// ad, err := m.GetTokenMetaData(token)
	// if err != nil {
	// 	return nil, err
	// }
	fmt.Printf("SetupPatientFilter:520 -- as: %s\n", spew.Sdump(as))
	//patientFilter.AccessDetails = ad
	patientFilter.SessionId = as.SessionID
	patientFilter.UserId = as.UserID.Hex()
	// patientFilter.JWTokenStr = tokenStr
	// patientFilter.JWToken = token
	// patientFilter.TokenCookie = tokenCookie
	count := r.Header.Get("Count")
	if count == "" {
		//fmt.Printf("   @@@ Setting default Count: %s\n\n", config.RecordLimit())
		patientFilter.Count = config.RecordLimit()
	} else {
		patientFilter.Count = count
	}

	// session, err := m.ValidateAuth(r.Header.Get("AUTHORIZATION"))
	// if err != nil {
	// 	session = &m.AuthSession{}
	// 	session.Token = r.Header.Get("AUTHORIZATION")
	// 	session.CreateSession()
	// 	//return nil, err
	// }
	// if session == nil {
	// 	err = fmt.Errorf("authorization was not found")
	// 	return nil, err
	// }
	//session.UpdateExpire()
	//patientFilter.Session = *session
	//patientFilter.SessionId = patientFilter.Session.PatSessionID
	//UseCache := r.Header.Get("UseCache")
	// count := r.Header.Get("Count")
	// if count == "" {
	// 	//fmt.Printf("   @@@ Setting default Count: %s\n\n", config.RecordLimit())
	// 	patientFilter.Count = config.RecordLimit()
	// }

	if r.Header.Get("ResultFormat") != "" {
		patientFilter.ResultFormat = r.Header.Get("ResultFormat")
	} else if r.Header.Get("RESULTFORMAT") != "" {
		patientFilter.ResultFormat = r.Header.Get("RESULTFORMAT")
	} else {
		patientFilter.ResultFormat = "ca" //Default to ca-ChartArchive
	}

	return &patientFilter, nil
}
