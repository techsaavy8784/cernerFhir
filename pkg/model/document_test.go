package model

import (
	"fmt"
	"testing"
	"time"
	"github.com/davecgh/go-spew/spew"

	. "github.com/smartystreets/goconvey/convey"
)


func TestFhirDocumentById(t *testing.T) {
	setupTest("")
	Convey("FhirDocumentById", t, func() {
		doc, err := FhirDocumentById("197588922")
		So(err, ShouldBeNil)
		So(doc, ShouldNotBeNil)
		fmt.Printf("Doc: %s\n", spew.Sdump(doc))
	})
}


func TestSearchReports(t *testing.T) {
	as := setupTest("")
	//as := setupTest("")
	// session, _ := ValidateAuth("test")
	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
	Convey("Subject: Document Search", t, func() {
		df := DocumentFilter{PatientGPI: "12724068", Session: as, SessionId: as.DocSessionId, Page: 1}
		df.Session.Status.Diagnostic = "done"
		df.Session.UpdateDiagStatus("done")
		df.Session.Status.Reference = "done"
		df.Session.UpdateRefStatus("done")
		inCache, _ := df.DocumentsInCache()
		fmt.Printf("###TestSearchReports:24 -- Before delete %d in cache\n\n", inCache)
		DeleteDocuments("12724068")
		inCache, _ = df.DocumentsInCache()
		fmt.Printf("###TestSearchReports:27 -- After delete %d in cache\n\n", inCache)

		fmt.Printf("\n\n   Document Search for Patient PatientGPI: %s\n", df.PatientGPI)
		startTime := time.Now()
		//docs, cacheStatus, inPage, pagesInCache, totalInCache, err := df.SearchReports()
		docs, _, _, _, totalInCache, err := df.SearchReports()
		elapsedTime := time.Since(startTime).Seconds()

		fmt.Printf("\n#### Returns with %d in cache in %f seconds\n", totalInCache, elapsedTime)
		So(err, ShouldBeNil)
		So(len(docs), ShouldNotEqual, 0)
		//fmt.Printf("TestSearchReports: 28 -- status: %s,inPage: %d, pagesInCache: %d, totalInCache: %d\n", cacheStatus, inPage, pagesInCache, totalInCache)
		//So(cacheStatus, ShouldEqual, "done")
		//So(len(caDocs), ShouldBeGreaterThan, 1)
		//fmt.Printf("InPage: %d, Pages: %d, Total: %d, CacheStatus: %s\n", inPage, pagesInCache, totalInCache, cacheStatus)
		//fmt.Printf("Retrieved document : %s\n", spew.Sdump(caDocs[0]))
		//fmt.Printf("Giving a chance to collect more\n\n\n")
		//cacheStatus, pagesInCache, totalInCache, err := df.DocumentCacheStats()
		//So(err, ShouldBeNil)
		//So(cacheStatus, ShouldEqual, "filling")
		//So(pagesInCache, ShouldEqual, 1)
		//So(totalInCache, ShouldEqual, 6)
		// fmt.Printf("\n\n\n### Waiting for background to finish\n")
		time.Sleep(3 * time.Second)
		incache, _ := df.CountCachedFhirDocuments()
		fmt.Printf("Actual inserted %d in cache in %f seconds\n", incache, elapsedTime)
		// cacheStatus, pagesInCache, totalInCache, _ = df.DocumentCacheStats()
		// fmt.Printf("status: %s, pagesInCache: %d, totalInCache: %d\n", cacheStatus, pagesInCache, totalInCache)

	})
}

func TestDocsForEncounter(t *testing.T) {
	as := setupTest("")
	df := DocumentFilter{PatientGPI: "12724066", EncounterID: "97953477", Session: as, SessionId: as.DocSessionId, Page: 1}

	Convey("Subject: Document Search", t, func() {
		docs, _, _, _, _, err := df.GetFhirDocumentPage()
		So(err, ShouldBeNil)
		So(len(docs), ShouldEqual, 10)
	})
}

func TestDeleteReports(t *testing.T) {
	//as := setupTest("")
	//as := setupTest("")
	// session, _ := ValidateAuth("test")
	// fmt.Printf("\n\n     #### Setup DocumentByID session: [%s]\n\n", session.SessionID)
	Convey("Subject: Document Search", t, func() {
		//df := DocumentFilter{PatientGPI: "12724066", Session: as, SessionId: as.DocSessionId, Page: 1}
	})
}
