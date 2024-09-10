package handler

import (
	//http "net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	fhir "github.com/dhf0820/cernerFhir/fhirongo"

	//h "github.com/dhf0820/cernerFhir/pkg/handler"

	//"github.com/davecgh/go-spew/spew"
	"github.com/dhf0820/cernerFhir/pkg/ca"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSearchDocuments(t *testing.T) {
	// godotenv.Load("env.test")
	// m.InitializeAll("")
	// sessionName := "test"
	// session := m.AuthSession{SessionID: sessionName}
	as := setupTest("")
	Convey("Subject: Cerner Documents", t, func() {
		//m.DeleteDocuments("12724068")
		Convey("Given a patient number get documents and return CA ", func() {
			//req := httptest.NewRequest("GET", "/api/v1/patient/12724066/documents?mode=ca&limit=10&page=1&visit_num=4027912", nil)
			req := httptest.NewRequest("GET", "/api/rest/v1/documents?patient_gpi=12724068&mode=ca&page=1&cache=clear", nil)
			req.Header.Set("SESSION", as.SessionID)
			//req.Header.Set("RESULTFORMAT")
			resp := httptest.NewRecorder()
			NewRouter().ServeHTTP(resp, req)
			//defer resp.Body.Close()
			//b, _ := ioutil.ReadAll(resp.Body)
			//docs, err := DocumentResults(resp)
			//_, err := DocumentResults(resp)
			//So(err, ShouldBeNil)
			//spew.Dump(fDocs)
			caDocs, err := CADocumentResults(resp)
			So(err, ShouldBeNil)
			//spew.Dump(caDocs)
			So(resp.Code, ShouldEqual, 200)
			So(len(caDocs), ShouldEqual, 20)
		})
	})
}

func TestDocumentsForPage(t *testing.T) {
	// godotenv.Load("env.test")
	// m.InitializeAll("")
	// sessionName := "test"
	// session := m.AuthSession{SessionID: sessionName}
	as := setupTest("")
	Convey("Subject: Cerner Documents", t, func() {
		//m.DeleteDocuments("12724068")
		Convey("Given a patient number get documents and return CA ", func() {
			//req := httptest.NewRequest("GET", "/api/v1/patient/12724066/documents?mode=ca&page=1&visit_num=4027912", nil)
			req := httptest.NewRequest("GET", "/api/rest/v1/documents?patient_gpi=12724066&visit_num=97953477&mode=ca&page=1&cache=clear", nil)
			req.Header.Set("SESSION", as.SessionID)
			//req.Header.Set("RESULTFORMAT")
			resp := httptest.NewRecorder()
			NewRouter().ServeHTTP(resp, req)
			//defer resp.Body.Close()
			//b, _ := ioutil.ReadAll(resp.Body)
			//docs, err := DocumentResults(resp)
			//_, err := DocumentResults(resp)
			//So(err, ShouldBeNil)
			//spew.Dump(fDocs)
			caDocs, err := CADocumentResults(resp)
			So(err, ShouldBeNil)
			//spew.Dump(caDocs)
			So(resp.Code, ShouldEqual, 200)
			So(len(caDocs), ShouldEqual, 10)
		})
	})
}
func TestDocHandlerAuthorized(t *testing.T) {
	as := setupTest("")
	Convey("Given the  User is authorized with AUTHORIZATION Header", t, func() {
		//m.DeleteDocuments(session)
		Convey("Authorized with a valid session", func() {
			//fmt.Printf("\n\n@@@ Check for valid session: %s\n\n", spew.Sdump(ad))
			req := httptest.NewRequest("GET", "/api/v1/patient/12746484/documents?mode=ca&page=1&visit_num=4027912", nil)
			// req.Header.Set("AUTHORIZATION", session.Token)
			req.Header.Set("ResultFormat", "ca")
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {
				m.DeleteDocuments(as.DocSessionId)
				NewRouter().ServeHTTP(resp, req)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 200)
				})
			})
		})
	})
}

func TestDocHandlerUnauthorized(t *testing.T) {
	// godotenv.Load("../../.env.test")
	// m.InitializeAll("")
	// sessionName := "test"
	// session := m.AuthSession{SessionID: sessionName}
	// session.Insert()
	//session := setupTest()
	Convey("Subject: Document requires an AUTHORIZARION header", t, func() {
		//m.DeleteDocuments(as.DocSessionID)
		Convey("Unauthorized with out a valid session", func() {
			fmt.Printf("\n\n@@@ Check for valid session\n\n")
			req := httptest.NewRequest("GET", "/api/v1/patient/4342009/documents?mode=ca&limit=20&page=1&visit_num=4027912", nil)
			//req.Header.Set("AUTHORIZATION", "sessionName")
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {

				NewRouter().ServeHTTP(resp, req)
				//defer resp.Body.Close()
				//b, _ := ioutil.ReadAll(resp.Body)

				//_, err := CADocumentResults(resp)
				//spew.Dump(docs)
				Convey("Then the response should be a 401", func() {
					//So(err, ShouldBeNil)
					So(resp.Code, ShouldEqual, 401)
				})
			})
		})
	})
}

// func TestSearchDocuments(t *testing.T) {
// 	// godotenv.Load("env.test")
// 	// m.InitializeAll("")
// 	// sessionName := "test"
// 	// session := m.AuthSession{SessionID: sessionName}
// 	as := setupTest("")
// 	Convey("Subject: Cerner Documents", t, func() {
// 		m.DeleteDocuments(as.DocSessionId)
// 		Convey("Given a patient number get documents and return CA ", func() {
// 			//req := httptest.NewRequest("GET", "/api/v1/patient/12724066/documents?mode=ca&limit=10&page=1&visit_num=4027912", nil)
// 			req := httptest.NewRequest("GET", "/api/rest/v1/documents?patient_gpi=12724068&mode=ca&page=1", nil)
// 			req.Header.Set("SESSION", as.SessionID)
// 			//req.Header.Set("RESULTFORMAT")
// 			resp := httptest.NewRecorder()
// 			NewRouter().ServeHTTP(resp, req)
// 			//defer resp.Body.Close()
// 			//b, _ := ioutil.ReadAll(resp.Body)

// 			docs, err := CADocumentResults(resp)
// 			fmt.Printf("\n!   CaDocumentResults error: %v\n", err)
// 			//spew.Dump(docs)
// 			So(resp.Code, ShouldEqual, 200)
// 			So(len(docs), ShouldEqual, 5)
// 		})
// 	})
// }
func TestDocHandlerByVisitNum(t *testing.T) {
	// godotenv.Load("../../.env.test")
	// m.InitializeAll("")
	// sessionName := "test"
	// session := m.AuthSession{SessionID: sessionName}
	// session.Insert()
	as := setupTest("")
	Convey("Subject: Document by Visit_number", t, func() {
		m.DeleteDocuments(as.DocSessionId)
		Convey("Given a visit_num get documents and return CA ", func() {
			req := httptest.NewRequest("GET", "/api/v1/patient/4342009/documents?mode=ca&limit=10&page=1&visit_num=4027912", nil)
			//req.Header.Set("AUTHORIZATION", session.Token)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {
				NewRouter().ServeHTTP(resp, req)
				//defer resp.Body.Close()
				//b, _ := ioutil.ReadAll(resp.Body)

				docs, err := CADocumentResults(resp)
				fmt.Printf("\n!   CaDocsResults error: %v\n", err)
				//spew.Dump(docs)
				Convey("Then the response should be a 200", func() {
					So(resp.Code, ShouldEqual, 200)
					So(len(docs), ShouldEqual, 5)
				})
			})
		})
	})
}

func TestDocHandlerByPatientID(t *testing.T) {
	as := setupTest("")
	Convey("Subject: Document by Visit_number", t, func() {
		m.DeleteDocuments(as.DocSessionId)
		Convey("Given a patientID get documents and return CA mode", func() {
			req := httptest.NewRequest("GET", "/api/v1/patient/12724066/documents?mode=ca&limit=10&page=1", nil)
			//req.Header.Set("AUTHORIZATION", session.Token)
			resp := httptest.NewRecorder()

			Convey("When the request is handled by the router", func() {
				NewRouter().ServeHTTP(resp, req)
				//defer resp.Body.Close()
				//b, _ := ioutil.ReadAll(resp.Body)

				docs, err := CADocumentResults(resp)

				Convey("Then the response should be a 200", func() {
					So(err, ShouldBeNil)
					So(resp.Code, ShouldEqual, 200)
					So(len(docs), ShouldEqual, 10)
				})
			})
		})
	})
}

func TestDocHandlerByMRNCA(t *testing.T) {
	//ad := setupAD()
	as := setupTest("")
	Convey("Subject: Document by MRN Returning CA", t, func() {
		//m.DeleteDocuments(ad)
		Convey("Given a query string of mrn=10002701 in CA mode", func() {
			req := httptest.NewRequest("GET", "/api/v1/documents?mrn=10002701&mode=ca&limit=20&page=1", nil)
			req.Header.Set("AUTHORIZATION", as.SessionID)
			resp := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {

				NewRouter().ServeHTTP(resp, req)
				//defer resp.Body.Close()
				//b, _ := ioutil.ReadAll(resp.Body)

				docs, err := CADocumentResults(resp)

				Convey("Then the response should be a 200", func() {
					So(err, ShouldBeNil)
					So(resp.Code, ShouldEqual, 200)
					So(len(docs), ShouldEqual, 20)
				})
			})
		})
	})
}

// func TestDocumentByMRNCA(t *testing.T) {
// 	godotenv.Load("../../.env.test")
// 	m.InitializeAll("")
// 	sessionName := "test"
// 	session := m.AuthSession{Token: sessionName}
// 	session.Insert()
// 	//m.DeleteDocuments(sessionName)
// 	Convey("Subject: Document by MRN Returning CA", t, func() {
// 		Convey("Given a query string of mrn=10002701", func() {
// 			req := httptest.NewRequest("GET", "/api/v1/documents?mrn=10002701&limit=20&page=1", nil)
// 			req.Header.Set("AUTHORIZATION", session.Token)
// 			resp := httptest.NewRecorder()
// 			Convey("When the request is handled by the router", func() {

// 				NewRouter().ServeHTTP(resp, req)
// 				//defer resp.Body.Close()
// 				//b, _ := ioutil.ReadAll(resp.Body)

// 				docs, _ := DocumentResults(resp)

// 				Convey("Then the response should be a 200", func() {
// 					//So(true, ShouldBeTrue)
// 					So(resp.Code, ShouldEqual, 200)
// 					//spew.Dump(docs)
// 					So(len(docs), ShouldEqual, 20)
// 					//spew.Dump(docs)
// 				})
// 			})
// 		})
// 	})
// }

// func DocumentResults(resp *httptest.ResponseRecorder) ([]*m.DocumentSummary, error) {
// 	//spew.Dump(resp.Body)
// 	b, _ := ioutil.ReadAll(resp.Body)
// 	//data := []*m.DocumentSummary{}
// 	respData := DocumentSummaryResponse{}

// 	if err := json.Unmarshal(b, &respData); err != nil {
// 		fmt.Printf("DocumentResults Error: %v\n", err)
// 		return nil, err
// 	}
// 	documents := respData.Documents
// 	//spew.Dump(respData)
// 	fmt.Printf("Test returns %d documents\n", len(documents))
// 	return documents, nil

// }

func CADocumentResults(resp *httptest.ResponseRecorder) ([]*ca.CADocument, error) {
	b, _ := ioutil.ReadAll(resp.Body)
	respData := ca.CaDocumentResponse{}

	if err := json.Unmarshal(b, &respData); err != nil {
		fmt.Printf("@      CADocumentResults Error: %v\n", err)
		return nil, err
	}
	documents := respData.Documents
	//fmt.Printf("CaDocumentResults:265 -- Test returns %d documents\n", len(documents))
	return documents, nil
}

func DocumentResults(resp *httptest.ResponseRecorder) ([]*fhir.Document, error) {
	b, _ := ioutil.ReadAll(resp.Body)
	respData := DocumentResponse{}

	if err := json.Unmarshal(b, &respData); err != nil {
		fmt.Printf("@      DocumentResults Error: %v\n", err)
		return nil, err
	}
	documents := respData.Documents
	//fmt.Printf("DocumentREsults:278 -- Test returns %d documents\n", len(documents))
	return documents, nil
}

// func setupTest() *m.AuthSession {
// 	godotenv.Load("../../.env.test")
// 	m.InitializeAll("")

// 	_, err := m.CreateSession("test")
// 	if err != nil {
// 		return nil
// 	}
// 	session := m.ValidateAuth("test")
// 	return session
// }
