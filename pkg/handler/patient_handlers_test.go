package handler

import (
	//http "net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	//fhir "github.com/dhf0820/cernerFhir/fhirongo"
	//m "github.com/dhf0820/cernerFhir/pkg/model"
	//h "github.com/dhf0820/cernerFhir/pkg/handler"
	//m "github.com/dhf0820/cernerFhir/pkg/model"
	//"github.com/davecgh/go-spew/spew"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

// func TestGetPatient(t *testing.T) {
// 	as := setupTest("")
// 	Convey("Subject: GetPatient returns the Specific Patient", t, func() {
// 		Convey("Given an ID: 12742397", func() {
// 			req := httptest.NewRequest("GET", "/api/rest/v1/patient/12742397", nil)
// 			req.Header.Set("ResultFormat", "ca")
// 			req.Header.Set("AUTHORIZATION", fmt.Sprintf("Bearer %s", ad.Token.Raw))
// 			resp := httptest.NewRecorder()
// 			NewRouter().ServeHTTP(resp, req)
// 			So(resp.Code, ShouldEqual, 200)

// 			//err: json.Unmarshal()
// 		})
// 	})

// }

func TestCaSearchPatient(t *testing.T) {
	//session := setupTest("handler-test")
	//fmt.Printf("\n\n#### Test token : [%s]\n", session.Token)
	var sessionId string
	//var query string
	as := setupTest("")
	sessionId = as.DocSessionId
	Convey("Subject: PatientSearch responds to queries properly", t, func() {
		// var sessionId string
		// //var query string
		// session := setupTest("handler-test")
		So(sessionId, ShouldNotBeEmpty)
		//fmt.Printf("\n\n#### Test token : [%s]\n", session)
		Convey("Given a query string of family=sma&given=baby&page=1", func() {
			//qry := "/api/v1/patient?family=sm&given=baby"
			qry := "/api/rest/v1/patient?family=smart&page=1&cache=reset"
			//query = qry
			req := httptest.NewRequest("GET", qry, nil)

			req.Header.Set("AUTHORIZATION", sessionId)
			fmt.Printf("\n\n### Execute GET for token: %s - qry: %s\n", sessionId, qry)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {
				startTime := time.Now()
				NewRouter().ServeHTTP(resp, req)
				fmt.Printf("ResponseTime: %f seconds\n", time.Since(startTime).Seconds())
				//defer resp.Body.Close()
				//b, _ := ioutil.ReadAll(resp.Body)

				//pats, _ := PatientResults(resp)

				caResp, err := CaPatientResults(resp)
				So(err, ShouldBeNil)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 200)
					fmt.Printf("Results: %s\n", spew.Sdump(resp))
					// So(caResp.Patients[0].Name, ShouldEqual, "SMART, BABY BOY")
					// So(caResp.Patients[0].MRN, ShouldEqual, "6946")
					// Total        int             `json:"total_documents"`
					// Pages        int             `json:"pages"`
					// NumberInPage int             `json:"page_documents"`
					// Page         int           	 `json:"page"`
					sessionId = caResp.SessionId
					fmt.Printf("Total: %d, Pages: %d, NumInPage: %d, Page: %d\n",
						caResp.Total, caResp.PagesInCache, caResp.NumberInPage, caResp.Page)
					fmt.Printf("SessionId: %s,  \n", sessionId)
				})
				Convey("Should be able to request page 2", func() {
					time.Sleep(5 * time.Second) // wait for a couple of pages
					pg2qry := "/api/rest/v1/patient?family=smart&page=3"
					req2 := httptest.NewRequest("GET", pg2qry, nil)
					req2.Header.Set("AUTHORIZATION", sessionId)
					resp2 := httptest.NewRecorder()
					startTime2 := time.Now()
					NewRouter().ServeHTTP(resp, req)
					fmt.Printf("ResponseTime: %f seconds\n", time.Since(startTime2).Seconds())
					caResp2, err := CaPatientResults(resp2)
					So(err, ShouldBeNil)
					So(caResp2, ShouldNotBeNil)
					sessionId = caResp.SessionId
					fmt.Printf("Total: %d, Pages: %d, NumInPage: %d, Page: %d\n",
						caResp.Total, caResp.PagesInCache, caResp.NumberInPage, caResp.Page)
					fmt.Printf(" sessionID Returned: %s\n", sessionId)

				})
			})
		})
		fmt.Printf("--Waiting for all pages to finish\n")
		time.Sleep(20 * time.Second)

		// Convey("RetrievingCacheInfo", func() {

		// })

		// SkipConvey("Given a query string of family=sma&given=baby use cache to find", func() {
		// 	fmt.Printf("\n\n\n\n Test query Cache Patients just found\n\n")
		// 	//qry := "/api/v1/patient?family=sm&given=baby"
		// 	qry := "/api/v1/patient?family=sm&limit=5"
		// 	req := httptest.NewRequest("GET", qry, nil)
		// 	req.Header.Set("AUTHORIZATION", session.SessionID)
		// 	req.Header.Set("UseCache", "true")
		// 	req.Header.Set("ResultFormat", "Full")  // Full, Summary, CA
		// 	req.Header.Set("Count", "3")
		// 	resp := httptest.NewRecorder()
		// 	Convey("When the request is handled by the router", func() {

		// 		NewRouter().ServeHTTP(resp, req)
		// 		//defer resp.Body.Close()
		// 		//b, _ := ioutil.ReadAll(resp.Body)

		// 		pats, _ := CaPatientResults(resp)

		// 		Convey("Then the response should be a 200", func() {
		// 			So(resp.Code, ShouldEqual, 200)

		// 			So(pats[0].Name, ShouldEqual, "SMART, BABY BOY")
		// 			So(pats[0].MRN, ShouldEqual, "6946")
		// 			//fmt.Printf("Found: %s\n", spew.Sdump(pats[0]))

		// 		})
		// 	})
		// })
		// SkipConvey("Given a query string of mrn=10002091", func() {

		// 	req := httptest.NewRequest("GET", "/api/v1/patient?mrn=10002091", nil)
		// 	resp := httptest.NewRecorder()
		// 	Convey("When the request is handled by the router", func() {

		// 		NewRouter().ServeHTTP(resp, req)
		// 		//defer resp.Body.Close()
		// 		//b, _ := ioutil.ReadAll(resp.Body)

		// 		caResp _ := CaPatientResults(resp)

		// 		Convey("Then the response should be a 200", func() {
		// 			So(resp.Code, ShouldEqual, 200)
		// 			So(caResp.Patients[0].Name, ShouldEqual, "Creevey, Colin Carl")
		// 			So(caResp.Patients[0].MRN, ShouldEqual, "10002091")
		// 		})
		// 	})
		// })
		// SkipConvey("Given a query string of Encounter=4027915", func() {
		// 	fmt.Printf("Find Patient for Encounter: 4027915\n")

		// 	req := httptest.NewRequest("GET", "/api/v1/patient?encounter=4027915", nil)
		// 	resp := httptest.NewRecorder()
		// 	Convey("When the request is handled by the router", func() {

		// 		NewRouter().ServeHTTP(resp, req)
		// 		//defer resp.Body.Close()
		// 		//b, _ := ioutil.ReadAll(resp.Body)

		// 		caResp, _ := CaPatientResults(resp)

		// 		Convey("Then the response should be a 200", func() {
		// 			So(resp.Code, ShouldEqual, 200)
		// 			So(pats[0].Name, ShouldEqual, "SMART, NANCY")
		// 			So(pats[0].MRN, ShouldEqual, "10002701")
		// 		})
		// 	})
		// })
	})
}

func CaPatientResults(resp *httptest.ResponseRecorder) (*PatientCaResponse, error) {
	b, _ := ioutil.ReadAll(resp.Body)
	//s := string(b)
	//fmt.Printf("\n\n\n\nraw response: \n%s\n", s)
	data := &PatientCaResponse{}
	fmt.Printf("StartUnmarshal starting\n")
	if err := json.Unmarshal(b, data); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return data, nil
}

func CaPatientResult(resp *httptest.ResponseRecorder) (*PatientCaResponse, error) {
	b, _ := ioutil.ReadAll(resp.Body)
	data := &PatientCaResponse{}
	if err := json.Unmarshal(b, &data); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return data, nil
}

// func PatientCaResult(resp *httptest.ResponseRecorder) (*m.CAPatient, error) {
// 	b, _ := ioutil.ReadAll(resp.Body)
// 	data := m.CAPatient{}
// 	if err := json.Unmarshal(b, &data); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 		return nil, err
// 	}
// 	return &data, nil
// }
