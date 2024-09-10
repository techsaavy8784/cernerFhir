package handler

import (
	//http "net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	//h "github.com/dhf0820/cernerFhir/pkg/handler"

	//"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
	ca "github.com/dhf0820/cernerFhir/pkg/ca"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFHIRSearch(t *testing.T) {
	as := setupTest("")
	Convey("Search for encounters for a patient", t, func() {
		//fmt.Printf("Session: %s\n", spew.Sdump(as))
		//m.DeleteDocuments(session)
		Convey("Authorized with a valid session", func() {
			//fmt.Printf("\n\n@@@ Check for valid session: %s\n\n", spew.Sdump(ad))
			//req := httptest.NewRequest("GET", "/api/rest/v1/encounters?patientGPI=12724066&page=1&count=20&cache=&page=1", nil)
			req := httptest.NewRequest("GET", "/api/rest/v1/encounters?patientGPI=12724068&page=1&format=ca&page=2", nil)
			req.Header.Set("SESSION", as.SessionID)
			req.Header.Set("RESULTFORMAT", "ca")
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {
				NewRouter().ServeHTTP(resp, req)
				So(resp.Code, ShouldEqual, 200)
				results, err := CAEncounterResults(resp)
				So(err, ShouldBeNil)
				So(len(results.Visits), ShouldBeGreaterThan, 10)

				// %s\n", spew.Sdump(results))

			})
		})
	})
}

func CAEncounterResults(resp *httptest.ResponseRecorder) (*ca.CaEncounterResponse, error) {
	b, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("\n\n### b = %s\n\n\n", string(b))
	respData := ca.CaEncounterResponse{}

	if err := json.Unmarshal(b, &respData); err != nil {
		fmt.Printf("@      CAEncounterResults Error: %v\n", err)
		return nil, err
	}
	//visits := respData.Visits
	//fmt.Printf("Test returns %d Visits\n", len(visits))
	return &respData, nil
}
