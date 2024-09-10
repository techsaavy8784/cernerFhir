package ca

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	m "github.com/dhf0820/cernerFhir/pkg/model"
)

type Visit struct {
	ID              uint
	PatientID       uint
	VisitNum        string     `json:"visit_num"`
	MRN             string     `json:"mrn"`
	AdmitDate       *time.Time `json:"admit_date"`
	DischargeDate   *time.Time `json:"discharge_date"`
	Facility        string     `json:"facility"`
	Clinic          string     `json:"clinic"`
	PatientType     string     `json:"patient_type"`
	HospitalService string     `json:"hospital_service"`
	PayorCode       string     `json:"payor_code"`
	FinancialClass  string     `json:"financial_class"`
	AdmitSource     string     `json:"admit_source"`
	Comment         string     `json:"comment"`
	Origin          string     `json:"origin"`
	EnterpriseId    string     `json:"enterprise_id"`
	PatientGPI      string     `json:"patient_gpi"`
	AccountNumber   string     `json:"account_number"`
	Text            string     `json:"text"`
}

type CaEncounterResponse struct {
	StatusCode   int      `json:"status_code"`
	Message      string   `json:"message"`
	CacheStatus  string   `json:"cache_status"`
	TotalInCache int64    `json:"total_visits"`
	PagesInCache int64    `json:"pages_in_cache"`
	NumberInPage int64    `json:"visits_in_page"`
	Page         int64    `json:"page"`
	SessionId    string   `json:"session_id"`
	Visits       []*Visit `json:"visits"`
	Visit        *Visit   `json:"visit"`
}

func WriteCaEncounterResponse(w http.ResponseWriter, resp *CaEncounterResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	return nil
}

func FhirEncountersToCA(w http.ResponseWriter, total, pages, inPage, page int64, cacheStatus string, encs []*fhir.Encounter) {
	//fmt.Printf("Convert from: %s\n\n", spew.Sdump(encs))
	caVisits := []*Visit{}
	//fmt.Printf("\n###visits:62 - number in array := %d\n", len(encs))
	for _, fenc := range encs {
		visit := FhirEncToCA(*fenc)
		caVisits = append(caVisits, &visit)
	}
	//fmt.Printf("CaVisits:67 -  %s\n", spew.Sdump(caVisits))

	resp := CaEncounterResponse{}
	resp.Message = "Ok"
	resp.StatusCode = 200
	resp.CacheStatus = cacheStatus
	resp.NumberInPage = inPage
	resp.PagesInCache = pages
	resp.TotalInCache = total
	resp.Page = page
	resp.Visits = caVisits
	WriteCaEncounterResponse(w, &resp)

}

func FhirEncToCA(enc fhir.Encounter) Visit {
	v := Visit{}

	v.VisitNum = *enc.Id
	v.AccountNumber = m.ExtractAccountNum(enc.Identifier)
	v.Text = enc.Text.Div
	v.Facility = enc.Location[0].Location.Display
	//v.EnterpriseId = *enc.Id
	// if enc.Subject.Display != "" {
	// 	patref := strings.Split(enc.Subject.Reference, "/")
	// 	v.PatientGPI = patref[1]
	// } else{
	patref := strings.Split(enc.Patient.Reference, "/")
	v.PatientGPI = patref[1]
	// }
	if enc.Period != nil {
		v.AdmitDate = &enc.Period.Start
		v.DischargeDate = &enc.Period.End
	}

	if len(enc.Location) > 0 {
		v.Clinic = enc.Location[0].Location.Display
	}
	v.PatientType = enc.Class //.Display

	return v
}

func GetFhirReference(ref fhir.Reference) string {
	return ref.Display
}
