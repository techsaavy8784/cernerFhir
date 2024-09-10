package model

import (
	"fmt"
	"testing"
	//"time"

	//"github.com/davecgh/go-spew/spew"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReptRefByID(t *testing.T) {
	as := setupTest("")
	//as := setupTest("")
	// session, _ := ValidateAuth("test")
	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
	Convey("Subject: Document By ID", t, func() {
		fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", as.DocSessionId)
		DeleteDocuments(as.DocSessionId)
		df := DocumentFilter{PatientGPI: "12746484",Session: as, SessionId: as.DocSessionId, Cache: "reset"}
		fmt.Printf("\n\n   Document Search for Patient PatientGPI: %s\n", df.PatientGPI)
		caDocs, cacheStatus, inPage, pagesInCache, totalInCache, err := df.SearchReports()
		So(err, ShouldBeNil)
		So(cacheStatus, ShouldEqual, "done")
		So(len(caDocs), ShouldBeGreaterThan, 1)
		fmt.Printf("InPage: %d, Pages: %d, Total: %d\n", inPage, pagesInCache, totalInCache)
		fmt.Printf("Retrieved document : %s\n", spew.Sdump(caDocs[0]))
	
		// Convey("Convverts To CA",func(){
		// 	cadoc := 
		// })
	})
}
//https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/DiagnosticReport?patient=12724066
// func TestReptRefByPatID(t *testing.T) {
// 	//as := setupTest("")
// 	as := setupTest("")

// 	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
// 	Convey("Subject: Document By ID", t, func() {
// 		fmt.Printf("as: %s\n", spew.Sdump(as))
// 		//session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", as.DocSessionId)
// 		DeleteDocuments(as.DocSessionId)
// 		Convey("Given a document by PatientGPI: 11031299 ", func() {
// 			// df := DocumentFilter{PatientGPI: "12724066" ,BeginDate: "01-03-2020", 
// 			// 		EndDate: "08-01-2020", Session: *session}

// 			df := DocumentFilter{PatientGPI: "12724066", Session: as, SessionId: as.DocSessionId, Cache: "reset"}
// 			fmt.Printf("\n\n   ReptRef Search for Patient PatientGPI: %s\n", df.PatientGPI)
// 			//docs, err := df.GetFhirDiagnosticReport()
// 			fmt.Printf("TestReptRefByPatII:52 - ef=%s\n", spew.Sdump(df))
// 			docs, cacheStatus, _, _, _, err := df.CaSearchReports()
// 			fmt.Printf("Docs: %s\n", spew.Sdump(docs))
// 			//docs, err := df.FindFhirDiagnosticReports()
// 			So(err, ShouldBeNil)
// 			//So(cacheStatus, ShouldEqual, "done")
// 			time.Sleep(time.Second *10)
// 			cacheStatus, pagesInCache, totalInCache, err := df.DocumentCacheStats()
// 			So(pagesInCache, ShouldBeGreaterThan, 1)
// 			So(totalInCache, ShouldBeGreaterThan, 20)
// 			So(cacheStatus, ShouldEqual, "done")

// 			//So(len(docs), ShouldBeGreaterThan, 0)
// 			// fmt.Printf("Documents Returned: %d\n", len(docs))
// 			// fmt.Printf("Diag: %s\n", spew.Sdump(docs[0]))
// 			// fmt.Printf("\n\n\nDiag: %s\n", spew.Sdump(docs[1]))
// 			So(len(docs), ShouldEqual, 10)  // 10 total documents we filtered out the one in August
// 			//caDocs, err := df.FindCaRepts()
// 			// caDoc := FhirDiagDocToCA(&resp.Entry[4].DiagnosticReport)
// 			// So(caDoc, ShouldNotBeNil)
// 			// fmt.Printf("caDoc: %s\n", spew.Sdump(caDoc))
		
// 		})
// 	})
// }