package model

import (
	"fmt"
	"testing"
	//"time"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)


// func TestDocumentSearch(t *testing.T) {
// 	as := setupTest("")
// 	//as := setupTest("")
// 	// session, _ := ValidateAuth("test")
// 	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
// 	Convey("Subject: Document Search", t, func() {
// 		//DeleteDocuments(as.DocSessionId)
// 		df := DocumentFilter{PatientGPI: "12724066", Session: as, SessionId: as.DocSessionId, Cache: "reset"}
// 		fmt.Printf("\n\n   Document Search for Patient PatientGPI: %s\n", df.PatientGPI)
// 		_, cacheStatus, inPage, pagesInCache, totalInCache, err := df.SearchReports()
// 		So(err, ShouldBeNil)
// 		//So(cacheStatus, ShouldEqual, "done")
// 		//So(len(caDocs), ShouldBeGreaterThan, 1)
// 		fmt.Printf("InPage: %d, Pages: %d, Total: %d, CacheStatus: %s\n", inPage, pagesInCache, totalInCache, cacheStatus)
// 		//fmt.Printf("Retrieved document : %s\n", spew.Sdump(caDocs[0]))
// 		fmt.Printf("Giving a chance to collect more\n\n\n")
// 		time.Sleep(120 * time.Second)
// 		// Convey("Convverts To CA",func(){
// 		// 	cadoc := 
// 		// })
// 	})
// }
// func TestDocumentByID(t *testing.T) {
// 	as := setupTest("")
// 	//as := setupTest("")
// 	// session, _ := ValidateAuth("test")
// 	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
// 	Convey("Subject: Document By ID", t, func() {
// 		fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", as.DocSessionId)
// 		DeleteDocuments(as.DocSessionId)
// 		df := DocumentFilter{PatientGPI: "12746484", SessionId: as.DocSessionId, Cache: "reset"}
// 		fmt.Printf("\n\n   Document Search for Patient PatientGPI: %s\n", df.PatientGPI)
// 		caDocs, cacheStatus, inPage, pagesInCache, totalInCache, err := df.CaSearchReports()
// 		So(err, ShouldBeNil)
// 		So(cacheStatus, ShouldEqual, "done")
// 		So(len(caDocs), ShouldBeGreaterThan, 1)
// 		fmt.Printf("InPage: %d, Pages: %d, Total: %d\n", inPage, pagesInCache, totalInCache)
// 		fmt.Printf("Retrieved document : %s\n", spew.Sdump(caDocs[0]))
	
// 		// Convey("Convverts To CA",func(){
// 		// 	cadoc := 
// 		// })
// 	})
// }
//https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/DiagnosticReport?patient=12724066
// func TestCaDiagnosticByPatID(t *testing.T) {
// 	//as := setupTest("")
// 	as := setupTest("")

// 	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
// 	Convey("Subject: Document By ID", t, func() {
// 		fmt.Printf("AD: %s\n", spew.Sdump(as))
// 		//session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", as.DocSessionId)
// 		DeleteDocuments(as.DocSessionId)
// 		Convey("Given a document by PatientGPI: 11031299 ", func() {
// 			// df := DocumentFilter{PatientGPI: "12724066" ,BeginDate: "01-03-2020", 
// 			// 		EndDate: "08-01-2020", Session: *session}

// 			df := DocumentFilter{PatientGPI: "12724066", SessionId: as.DocSessionId, Cache: "reset"}
// 			fmt.Printf("\n\n   CaDiagnosticReport Search for Patient PatientGPI: %s\n", df.PatientGPI)
// 			//docs, err := df.GetFhirDiagnosticReport()
// 			docs, cacheStatus, _, _, _, err := df.CaSearchReports()
// 			fmt.Printf("Docs: %s\n", spew.Sdump(docs))
// 			//docs, err := df.FindFhirDiagnosticReports()
// 			So(err, ShouldBeNil)
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
// func TestDocumentByMRN(t *testing.T) {
// 	as := setupTest("")

// 	DeleteDocuments(as.DocSessionId)
// 	fmt.Printf("\n\n     #### Setup DocumentByMRN session: [%s]\n\n", as.DocSessionId)
// 	Convey("Subject: Document By MRN", t, func() {
// 		// session, _ := ValidateAuth("test")
// 		// DeleteDocuments(session.CacheName)
// 		// fmt.Printf("\n\n     #### Setup DocumentByMRN session: [%s]\n\n", session.CacheName)
// 		Convey("Given a query by MedicalRecordNumber", func() {
		
// 			DeleteDocuments(as.SessionID)
// 			fmt.Printf("\n\n     #### Setup DocumentByMRN session: [%s]\n\n", spew.Sdump(as))
// 			Convey("The request is received for Patient MRN: 10002701 ", func() {
// 				df := DocumentFilter{MRN: "10002701", Session: as, SessionId: as.DocSessionId}
// 				fmt.Printf("\n\n  Starting searc for mrn\n")
// 				//err := fmt.Errorf("Error")
// 				//fmt.Printf("   Document Search for mrn: %s\n", df.MRN)
// 				docs, _, err := df.Search()
// 				// fmt.Printf("     REturned %d documents, err: %v\n", len(docs), err)

// 				//docs := results

// 				Convey("There is no error", func() {
// 					fmt.Printf("There is no error\n")
// 					So(err, ShouldBeNil)
// 				})
// 				Convey("20 documents are found", func() {
// 					So(len(docs), ShouldEqual, 20)
// 					fmt.Printf("Test returned %d documents\n", len(docs))
// 					//spew.Dump(docs)
// 				})
// 			})
// 		})
// 	})
// }


// Need Mocks before we can do this

// func TestDocumentByMRNBeforDate(t *testing.T) {
// 	godotenv.Load("../../.env.test")
// 	InitializeAll("")
// 	sessionName := "Test"
// 	// session := new(AuthSession)
// 	// session.SessionID = sessionName

// 	session := ValidateAuth(sessionName)
// 	DeleteDocuments(session.CacheName)
// 	fmt.Printf("\n\n     #### Setup DocumentByMRNBeforDate session: [%s]\n\n", session.CacheName)
// 	Convey("Subject: DocumentSearch responds using PatientFilter", t, func() {
// 		//Can test when we have our own mocks
// 		// Need to pre query to get cache
// 		Convey("Given a query by MedicalRecordNumber before a date", func() {
// 			query := DocumentFilter{MRN: "10002701", EffectiveDate: "$lt|2018-01-01", Count: "10", Session: *session}
// 			Convey("The request is received for Patient MRN: 10002701 ", func() {
// 				//err := fmt.Errorf("Error")
// 				fmt.Printf("   Test 73 Search for mrn: %s  using cache: %s\n", query.MRN, session.CacheName)
// 				results, _, err := query.Search()
// 				docs := results

// 				Convey("There is no error", func() {
// 					So(err, ShouldEqual, nil)
// 				})
// 				Convey("1 documents are found", func() {
// 					So(len(docs), ShouldEqual, 1)
// 				})
// 			})
// 		})
// 	})
// }

// func TestDocumentByMrnToCA(t *testing.T) {
// 	godotenv.Load("../../.env.test")
// 	InitializeAll("")
// 	// session, _ := ValidateAuth("test")
// 	// fmt.Printf("\n\n     #### Setup DocumentByMRNToCA session: [%s]\n\n", session.CacheName)
// 	Convey("Subject: Document By MRN convert to CA", t, func() {
// 		session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup DocumentByMRNToCA session: [%s]\n\n", session.CacheName)
// 		DeleteDocuments(session.CacheName)
// 		Convey("Given a query by MedicalRecordNumber converting to ca with limit", func() {
// 			Convey("The request is received for Patient MRN: 10002701", func() {
// 				query := DocumentFilter{MRN: "10002701", ResultFormat: "ca", Limit: 5, Page: 1, Session: *session}
// 				fmt.Printf("\n\n\n +++++++++  Search for mrn: %s\n", query.MRN)
// 				results, _, _ := query.Search()
// 				docs := results
// 				caDocs := ConvertDocumentsToCA(docs)
// 				// Convey("There is no error", func() {
// 				// 	So(err("5 documents should be found", func() {
// 				So(len(caDocs), ShouldEqual, 5)
// 				// })
// 			})
// 		})
// 	})
// }

func TestDocumentByPatID(t *testing.T) {
	as := setupTest("")

	Convey("Subject: Document By PatientID andEncounter", t, func() {
		fmt.Printf("\n\n     #### Setup Document By Pat ID and Encounter session: [%s]\n\n",as.DocSessionId )
		DeleteDocuments(as.DocSessionId)
		Convey("Given a query by PatientID 4342009 with encounter 4027912", func() {
			df := DocumentFilter{PatientGPI: "12724066",  SessionId: as.DocSessionId}
			Convey("The request is received for Patient 4342009 and encounter 4027912", func() {
				//err := fmt.Errorf("Error")
				fmt.Printf("   Search for ID: %s\n", df.PatientID)
				fmt.Printf("Calling d.Search\n")
				// results, pages, total  err := df.SearchCAReports()
				// docs := results

				// Convey("There is no error", func() {
				// 	So(err, ShouldEqual, nil)
				// })
				// Convey("5 documents Should be found", func() {
				// 	So(len(docs), ShouldEqual, 5)
				// 	//spew.Dump(docs)
				// })
			})
		})
	})
}

// func TestFhirDiagByPatID(t *testing.T) {
// 	as := setupTest("")

// 	Convey("Subject: FhirDocument By PatientID andEncounter", t, func() {
// 		//fmt.Printf("\n\n     #### Setup Document By Pat ID and Encounter session: [%s]\n\n", as.Session)
// 		DeleteDocuments(as.DocSessionId)
// 		Convey("Given a query by PatientID 4342009 with encounter 4027912", func() {
// 			//df := DocumentFilter{PatientID: "12724066",  Session: as}
// 			df := DocumentFilter{PatientGPI: "12765407",  Session: as}
// 			Convey("The request is received for Patient 4342009 and encounter 4027912", func() {
// 				//err := fmt.Errorf("Error")
// 				fmt.Printf("   Search for ID: %s\n", df.PatientID)
// 				fmt.Printf("Calling FindFhirDiagRepts()\n")
// 				diagRepts, err := df.FindFhirDiagRepts()
// 				So(err, ShouldBeNil)
// 				So(diagRepts, ShouldNotBeNil)

// 				// results, pages, total  err := df.SearchCAReports()
// 				// docs := results

// 				// Convey("There is no error", func() {
// 				// 	So(err, ShouldEqual, nil)
// 				// })
// 				// Convey("5 documents Should be found", func() {
// 				// 	So(len(docs), ShouldEqual, 5)
// 				// 	//spew.Dump(docs)
// 				// })
// 			})
// 		})
// 	})
// }

// func TestDocumentForCategory(t *testing.T) {
// 	as := setupTest("")
// 	fmt.Printf("\n\n     #### Setup For Category session: [%s]\n\n", as.DocSessionId)
// 	Convey("Subject: Document for Category", t, func() {
// 		Convey("Given a query for Category Enhanced Note", func() {
// 			DeleteDocuments(as.DocSessionId)
// 			query := DocumentFilter{PatientGPI: "4342009", Category: "Enhanced Note", Count: "100", SessionId: as.DocSessionId}
// 			fmt.Printf("   Search for Category: %s\n", query.Category)
// 			Convey("The request is received for Patient 4342009 and Category Enhanced Note ", func() {
// 				//err := fmt.Errorf("Error")
// 				results, _, err := query.Search()
// 				docs := results

// 				Convey("There is no error", func() {
// 					So(err, ShouldEqual, nil)
// 				})
// 				Convey("3 documents are found", func() {
// 					So(len(docs), ShouldEqual, 0)
// 				})
// 			})
// 		})
// 	})
// }

// func TestDocumentForSkip(t *testing.T) {
// 	as := setupTest("")

// 	//fmt.Printf("\n\n     #### Setup testing Skip session: [%s]\n\n", spew.Sdump(as))
// 	Convey("Subject: TEst Skip on request", t, func() {
// 		DeleteDocuments(as.DocSessionId)
// 		Convey("Given a query of anything with a Skip parameter", func() {
// 			query := DocumentFilter{PatientGPI: "4342009", Skip: 1, Limit: 1, Count: "5", SessionId: as.DocSessionId}
// 			fmt.Printf("   Search documents skiping th first 1\n")
// 			Convey("The request is received for 1 record ", func() {
// 				//err := fmt.Errorf("Error")
// 				results, _, err := query.Search()
// 				docs := results

// 				Convey("There is no error", func() {
// 					So(err, ShouldEqual, nil)
// 				})
// 				Convey("1 documents are found", func() {
// 					So(len(docs), ShouldEqual, 1)
// 				})
// 			})
// 		})
// 	})
// }
