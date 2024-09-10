package model

import (
	//"bytes"
	"context"
	//"encoding/json"
	"errors"
	"fmt"

	//"os"
	//"strconv"

	//"net/http"
	//"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	"github.com/dhf0820/cernerFhir/pkg/storage"
	log "github.com/sirupsen/logrus"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

/////////////////////////////////////////////////////////////////////////////////////////
//                                  Handle CA                                       /
/////////////////////////////////////////////////////////////////////////////////////////
// SearchCAReports: Find and cache CA version of FhirDiagnosticReports and Find and
// cache CA version of FhirDocumentReferences returning requested page of the combined

// func (df *DocumentFilter) GetCaDocumentReferences() ([]*CADocument, error) {
// 	fhirDocRefs, err := df.GetFhirDocRefs()
// 	if err != nil {
// 		return nil, err
// 	}
// 	caDocRefs, err := df.FhirDocRefsToCADocuments(fhirDocRefs)
// 	return caDocRefs, err
// }

//Does not return anything. THis just fills the cache if necessary
func (df *DocumentFilter) GetFhirDocRefs() {
	fmt.Printf("\n////////////////////////  GetFhirDocRefs ////////////////////////////////\n")
	// c := config.Fhir()
	log.Debugf("GetFhirDocRefs is searching DocumentReferences for Patient: %s", df.PatientGPI)
	//returns cacheStatus, pagesInCache, totalInCache, error

	// _, _, totalInCache, _ := df.DocumentCacheStats()
	// cacheStatus := df.Session.GeReptRefStatus()
	//If the cacheStatus == done and there are some documents in cache we are done
	// if cacheStatus == done and there are documents in cache. we have the documents and do not restart caching
	// if cacheStatus == "done" && totalInCache > 0 {
	// 	fmt.Printf("FhirRef is Done and there are documents")
	// 	return
	// }

	// If cache is not "done" it is building and do not restart the caching
	// if cacheStatus != "done" {
	// 	fmt.Printf("FhirRef cacheStatus = %s we are Filling\n", cacheStatus)
	// 	return
	// 	// fhirDocs, _, _, _, _, err := df.GetFhirDocumentPage()
	// 	// return fhirDocs, err
	// }
	//if cacheStatus =="done" and nothing in cache, start caching

	df.CacheFhirDocRefs()
	fmt.Printf("GetFhirDocRefs:76 -- CacheFhirDocRef returned\n")

}

func (df *DocumentFilter) CacheFhirDocRefs() {
	c := config.Fhir()
	//start := time.Now()
	df.MakeFhirQuery()

	limit := RecordLimit()
	df.fhirQuery = fmt.Sprintf("?patient=%s&_count=%d", df.PatientGPI, limit)
	log.Debugf("\n\n\n###  c.FindDocumentReferences is called with %s", df.fhirQuery) //fmt.Sprintf("?patient=%s", df.PatientGPI))
	startTime := time.Now()
	results, err := c.FindDocumentReferences(df.fhirQuery)
	elapsedTime := time.Since(startTime).Seconds()
	// fmt.Printf("FindDocumentReferences returned\n")
	if err != nil {
		log.Errorf("c.FindDocumentReferences err: %s", err.Error())
		return
	}
	// fmt.Printf("CacheFhirDocRefs:94 -- FHIR returned results: %s\n", spew.Sdump(results))
	// fmt.Printf("\n\n\n###rept_ref:95 -- calling FollowNextRefLink: %s\n", spew.Sdump(results.Link))
	go df.FollowNextRefLinks(results.Link)
	// for _, l := range resp.Link {
	// 	if l.Relation == "next" {
	// 		nextLink := l.URL
	// 		log.Debugf("diagrepts:226 -- NextLink: %s\n", nextLink)
	// 		go df.ProcessRemainingDiagRepts(nextLink)
	// 		break
	// 	} // if no next link there are no more DiagRepts
	// }

	// log.Debugf("fhir returned no error")
	// refs := []*fhir.DocumentReference{}
	// entry := resp.Entry

	// log.Info("Processing the initial results of DocRefs")
	// for _, r := range entry {
	// 	ref := r.DocumentReference
	// 	refs = append(refs, &ref)
	// }
	log.Infof("CacheFhirDocRefs:118 -- for c.FindDocumentReferences elapse time: %f", elapsedTime)
	_, err = InsertFhirDocResults(results, df.Session.DocSessionId)
	if err != nil {
		log.Errorf("CacheFhirDocREfs:121 -- InsertFhirDocResults err: %s", err.Error())
		return
	}
}

// func (df *DocumentFilter) FhirDocRefsToCADocuments(fds []*fhir.DocumentReference) ([]*CADocument, error) {
// 	caDocuments := []*CADocument{}
// 	for _, d := range fds {
// 		doc := df.FhirDocRefToCADocument(d)
// 		caDocuments = append(caDocuments, doc)
// 	}
// 	//log.Debugf("FhirDiagReptsToCA returned %d documents", len(caDocuments))

// 	return caDocuments, nil
// }

// func (df *DocumentFilter) FhirDocRefToCADocument(fdr *fhir.DocumentReference) *CADocument {
// 	var caDoc CADocument
// 	var err error

// 	//fmt.Printf("Converting FhirDiag: %s\n", spew.Sdump(fd))

// 	caDoc.ID = fdr.ID
// 	caDoc.PatientGPI = GetFhirPerson(fdr.Subject, "ID")
// 	caDoc.VersionID, err = strconv.ParseUint(fdr.Meta.VersionID, 10, 64)
// 	if err != nil {
// 		log.Errorf("Invalid VersionID: [%s] error:%s\n", fdr.Meta.VersionID, err.Error())
// 	}
// 	enc := strings.Split(fdr.Context.Encounter.Reference, "/")
// 	if len(enc) > 1 {
// 		caDoc.Encounter = enc[1]
// 	} // Encounter is not available
// 	caDoc.Repository = "FHIR"
// 	//rpdt := fdr.EffectiveDateTime
// 	rpdt := fdr.Meta.LastUpdated
// 	caDoc.ReptDateTime = &rpdt
// 	caDoc.DocStatus = fdr.DocStatus.Text
// //	caDoc.Status = fdr.Status
// 	caDoc.ImageURL = GetRefImageURL(fdr, "application/pdf")
// 	caDoc.Pages = 0 // Unavailable
// 	//caDoc.Subtitle = fdr.Subtitle
// 	caDoc.Class = fdr.Type.Text
// 	caDoc.Text = fdr.Text.Div

// 	cfg := ActiveConfig()
// 	caDoc.Source = strings.ToLower(cfg.Source())
// 	if caDoc.Description == "" {
// 		caDoc.Description = fdr.Description
// 	}
// 	err = (caDoc).Insert(df.Session.DocSessionId)
// 	if err != nil {
// 		log.Errorf("%s", err.Error())
// 		return nil
// 	}
// 	return &caDoc
// }

////////////////////////////////////////////////////////////////////////////////////////////////////////
//                              Process Next Set of DiagRepts returned from  FHIR                     //
////////////////////////////////////////////////////////////////////////////////////////////////////////

func (df *DocumentFilter) FollowNextRefLinks(links []fhir.Link) {
	//fmt.Printf("\n\n\n### FollowNextRefLinks:186")
	// df.Session.Status.Reference = "filling"
	// df.Session.UpdateRefStatus( "filling")
	url := NextRefPageLink(links)
	i := 1
	//time.Sleep(10 * time.Second)
	for {
		startTime := time.Now()
		if url == "" {
			log.Info("CachePages for url is blank, done")
			df.Session.Status.Reference = "done"
			df.Session.UpdateRefStatus("done")
			break
		}
		//TODO: Get the next page and start its next page while processing current. Do in paraallel

		links, _ = df.ProcessRefPage(url)
		fmt.Printf("--------Reference Link Page: %d  added in %f seconds", i, time.Since(startTime).Seconds())
		i = i + 1
		url = NextRefPageLink(links)
	}
	df.Session.Status.Reference = "done"
	df.Session.UpdateRefStatus("done")
}

func NextRefPageLink(links []fhir.Link) string {
	for _, l := range links {
		//fmt.Printf("Looking at link: %s\n", spew.Sdump(l))
		if l.Relation == "next" {
			//fmt.Printf("##NextRefLink:215 %s \n\n", l.URL)
			return l.URL
		}
	}
	//fmt.Printf("NextRefPageLink:219 - No next link\n")
	return ""
}

func (df *DocumentFilter) ProcessRefPage(url string) ([]fhir.Link, error) {
	//fmt.Printf("\n\n\n###ProcessDiagPage:1998\n\n")
	//startTime := time.Now()

	results, err := fhirC.NextFhirDocRefs(url)
	if err != nil {
		log.Errorf("NextFHIRDocRefs returned err: %s\n", err.Error())
		return nil, err
	}

	// for _, entry := range results.Entry {

	// 	doc := entry.Document
	// 	doc.FullURL = entry.FullURL
	// 	err := InsertFhirDocument(&doc, df.Session.DocSessionId)
	// 	if err != nil {
	// 		msg := fmt.Sprintf("ProcessRefPage:217 --  failed: %s", err.Error())
	// 		log.Error(msg)
	// 		return nil, errors.New(msg)
	// 	}
	// }
	//fmt.Printf("\n\nNext Page of results: %s\n\n", spew.Sdump(fhirDocRefs))
	docs := []*fhir.Document{}
	entry := results.Entry

	for _, r := range entry {
		rpt := r.Document
		rpt.FullURL = r.FullURL

		InsertFhirDoc(&rpt, df.Session.DocSessionId)
		docs = append(docs, &rpt)
	}
	// refs, err := InsertFhirDocuments(fhirDocRefResults, df.Session.DocSessionId)
	// if err != nil {
	// 	return nil, err
	// }
	// refs := []*fhir.DocumentReference{}
	// entry := fhirDocRefResults.Entry

	// for _, r := range entry {
	// 	rpt := r.DocumentReference

	// 	refs = append(refs, &rpt)
	// }

	//_, err = df.FhirDocRefsToCADocuments(refs)
	//_, err = FhirDocRefsToCADocuments(fhirDiagRepts, df) //Documents are cached
	// if err != nil {
	// 	return nil, err
	// }
	//patients := parsePatientResults(fhirPatients, f.SessionId)
	//spew.Dump()

	if len(docs) == 0 {
		return nil, fmt.Errorf("404|no more Documents found")
	}
	return results.Link, nil
}

// func InsertFhirDocRefs(results *fhir.DocumentResults, sessionId string) ([]*fhir.DocumentReference, error) {
// 	//entry := results.Entry
// 	refs := []*fhir.DocumentReference{}
// 	for _, entry := range results.Entry {
// 		docRef := entry.DocumentReference
// 		refs = append(refs, &docRef)
// 		err := InsertFhirDocRef(&docRef, sessionId)
// 		if err != nil {
// 			msg := fmt.Sprintf("InsertFhirDocRef using sessionId: %s failed: %s", sessionId, err.Error())
// 			log.Error(msg)
// 			return nil, errors.New(msg)
// 		}
// 	}
// 	return refs, nil
// }

// func InsertFhirDocRefs(refs []*fhir.DocumentReference, sessionId string) error {
// 	for _, ref:= range refs {
// 		err := InsertFhirDocRef(ref, sessionId)
// 		if err != nil {
// 			msg := fmt.Sprintf("InsertFhirPat using sessionId: %s failed: %s", sessionId, err.Error())
// 			log.Error(msg)
// 			return errors.New(msg)
// 		}
// 	}
// 	return nil
// }

func InsertFhirDocRef(docRef *fhir.DocumentReference, sessionId string) error {
	//fmt.Printf("   Inserting InsertFhirDocRef:311 -- ID:%s\n", docRef.ID)
	docRef.SessionId = sessionId
	collection, _ := storage.GetCollection("documents")
	docRef.CacheID = primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), docRef)
	if err != nil {
		if !storage.IsDup(err) {
			msg := fmt.Sprintf("InsertDocRefCache error: %s", err.Error())
			log.Error(msg)
			return errors.New(msg)
		} else {
			err = nil
		}
	}
	return err
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                 CA DocumentReference Handlers                                    //
//////////////////////////////////////////////////////////////////////////////////////////////////////
// //  FhirDocRefsToCA: Accepts an slice of fhir.DocumentReferences and returns a slice of *CADocuments
// func (df DocumentFilter) FhirDocRefsToCA(docs []*fhir.DocumentReference) []*CADocument {
// 	startTime := time.Now()
// 	var caDocuments []*CADocument
// 	log.Debugf("Converting %d documents to CA documents\n", len(docs))
// 	// for _, d := range docs {
// 	// 	doc := d.ToCA()
// 	// 	caDocuments = append(caDocuments, doc)
// 	// }

// 	log.Infof("Convert %d documents to CA took %s\n", len(caDocuments), time.Since(startTime))
// 	return caDocuments
// }

/////////////////////////////////////////////////////////////////////////////////////////
//                                 FHIR Getters                                         /
/////////////////////////////////////////////////////////////////////////////////////////

func GetRefImageURL(fdr *fhir.DocumentReference, imageType string) string {
	for _, attachment := range fdr.PresentedForm {
		//attachment := cnt.Attachment
		if attachment.ContentType == "application/pdf" {
			return attachment.URL
		} else {
			log.Warnf("Other attachment types for %s : %s", fdr.ID, spew.Sdump(attachment))
		}
	}
	return ""
}

// type DocumentReference struct {
// 	ResourceType      string    `json:"resourceType"`
// 	SessionId	  string    `json:"-"`
// 	ID                string    `json:"id"`
// 	FullURL 		string 		`json:"fullUrl"`
// 	EffectiveDateTime time.Time `json:"effectiveDateTime"`
// 	Meta              MetaData  `json:"meta"`
// 	Text              TextData  `json:"text"`
// 	Subject           Person    `json:"subject"`
// 	Type              Concept   `json:"type"`
// 	Authenticator     Person    `json:"authenticator"`
// 	Created           time.Time `json:"created"`
// 	Indexed           time.Time `json:"indexed"`
// 	DocStatus         Concept   `json:"docStatus"`
// 	Description       string    `json:"description"`

// 	Content []struct {
// 		Attachment struct {
// 			ContentType string `json:"contentType"`
// 			URL         string `json:"url"`
// 		} `json:"attachment"`
// 	} `json:"content"`
// 	//} `json:"content"`
// 	//Content       []Attachment `json:"content"`
// 	Context struct {
// 		EncounterNum struct {
// 			Reference string `json:"reference"`
// 		} `json:"encounter"`
// 	} `json:"context"`
// }
// type Thing struct {
// 	Display   string `json:"display"`
// 	Reference string `json:"reference"`
// }

// Person is a human
// type Person Thing

// type CaResults struct {
// 	Documents  []*DocumentSummary	`json:"documents"`
// 	TotalPages uint					`json:"total_pages"`
// }

// type DocumentSummary struct {
// 	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
// 	SessionID    string             `json:"-"`
// 	PatientID    string             `json:"patient_id" bson:"patient_id"`
// 	PatientGPI   string             `json:"patient_gpi"`
// 	Encounter    string             `json:"visit_num" bson:"encounter"`
// 	Source       string             `json:"source" bson:"source"`
// 	FullLink     string             `json:"full_link" bson:"full_link"`
// 	Meta         fhir.MetaData      `json:"meta" bson:"meta"`
// 	Category     string             `json:"category" bson:"category"`
// 	Code         string             `json:"code" bson:"code"`
// 	Description  string             `json:"description" bson:"description"`
// 	Text         fhir.TextData      `json:"text" bson:"text"`
// 	ReptDatetime time.Time          `json:"rept_datetime" bson:"rept_datetime"`
// 	Issued       time.Time          `json:"issued" bson:"issued"`
// 	Subject      fhir.Person        `json:"subject" bson:"subject"`
// 	Performer    fhir.Person        `json:"performer" bson:"performer"`
// 	EnterpriseID string             `json:"documentid" bson:"enterprise_id"`
// 	CADocumentID uint64             `json:"document_id" bson:"document_id"`
// 	SourceType   string             `json:"sourceType" bson:"source_type"`
// 	CreatedAt    *time.Time          `json:"created_at" bson:"created_at"`
// 	UpdatedAt    *time.Time          `json:"updated_at" bson:"updated_at"`
// 	AccessedAt   *time.Time          `json:"accessed_at" bson:"accessed_at"`
// 	ImageURL     string             `json:"image_url" bson:"image_url"`
// 	FhirImageURL string             `json:"fhir_image_url" bson:fhir_&image_url`
// 	Images       []PresentedForm    `json:"images" bson:"images"`
// 	Image        string             `json:"image" bson:"image"` // base64
// 	Versionid    uint64             `json:"version_id" bson:"version_id"`
// 	TotalPages   int                `json:"total_pages"`
// 	//Image         primitive.Binary   `json:"image" bson:"image"`
// }

// type CADocument struct {
// 	ID           	string		`json:"document_id"`
// 	VersionID   	uint64		`json:"version_id"`
// 	Encounter    	string		`json:"visit_num"`
// 	Repository		string		`json:"repository"`
// 	Source			string 		`json:"source"`
// 	Description		string		`json:"description"`
// 	ImageURL 		string 		`json:"image_url"`
// 	Pages 			int 		`json:"pages"`
// 	ReptDateTime 	*time.Time	`json:"rept_datetime"`
// 	Subtitle		string		`json:"subtitle"`
// 	PatientGPI		string		`json:"patient_gpi"`
// 	Text			fhir.TextData `json:"text"`
// }

// type PresentedForm struct {
// 	ContentType string `json:"context_type" bson:"context_type"`
// 	URL         string `json:"url" bson:"url"`
// }

// type DocumentImage struct {
// 	ResourceType string        `json:"resourceType" bson:"resource_type"`
// 	ID           string        `json:"id" bson:"id"`
// 	Meta         fhir.MetaData `json:"meta" bson:"meta"`
// 	ContentType  string        `json:"contentType" bson:"content_type"`
// 	Content      string        `json:"content" bson:"content"`
// }

// type DocumentFilter struct {
// 	Page       int64    `schema:"page"`
// 	Skip       int64    `schema:"skip"`
// 	Limit      int64    `schema:"limit"`
// 	SortBy     []string `schema:"sort"`
// 	Column     string   `schema:"column"`
// 	Order      string   `schema:"order"`
// 	ResultFormat	string   `schema:"result_format"`//Header
// 	Count      string   `schema:"count"`								//Header
// 	UseCache   string   `schema:"useCache"`							//Header
// 	ClearCache string   `schema:"clearCache"`

// 	Session AuthSession
// 	Source  string
// 	// Specified filters.
// 	// PatientID/ID is the only one to query FHIR with
// 	// Once all documents are in the mongo cache we do the requested query using mongo query

// 	PatientID    string `schema:"patientid"`
// 	PatientGPI   string `schema:"patient_gpi"`
// 	Encounter    string `schema:"encounter"`
// 	MRN          string `schema:"mrn"`
// 	ID           string `schema:"id"`
// 	EnterpriseID string `schema:"enterpriseid"`
// 	ReptDatetime string `schema:"rept_datetime"`
// 	Category     string `schema:"category"`
// 	EncounterID  string `schema:"visit_num"`
// 	//Column        string `schema:"column"`
// 	SourceValues []string `schema:"source_values"`
// 	BeginDate 	 string		`schema:"begin_date"`
// 	EndDate			 string 	`schema:"end_date"`
// 	TabID        string   `schema:"tab_id"`

// 	queryString     string
// 	queryFhirString string
// 	queryMap        map[string]string
// 	QueryFilter     bson.M
// 	QueryFilterBase []bson.M
// }
/*
//var activeDocumentFilter DocumentFilter
////////////////////////////////////////////////////////////////////////////////////////////////
//                                           Search                                           //
////////////////////////////////////////////////////////////////////////////////////////////////
//
// Search searches for both DiagnosticDocuments and DocumentReferences converting them to CA  //
// and allowing sorting by date, document_type, encounterId. The CA format is agnostic as to  //
// the source.
////////////////////////////////////////////////////////////////////////////////////////////////
func (f *DocumentFilter) Search() ([]*fhir.DiagnosticReport, int64, error) {

	//log.Debugf("\n\nCurrent Filter:\n")

	err := f.makeQueryFilter()
	if err != nil {
		return nil, 0, err
	}
	log.Debugf("Document#Search Query Filters: %v\n", f.QueryFilter)
	startTime := time.Now()

	if f.ID != "" {
		f.UseCache = "true"
		//ds := new(DocumentSummary)
		summaries, cnt, err := f.FindDocuments()
		if err != nil {
			return nil, cnt, err
		}
		// err = ds.fromDocumentReference(ref, f.SessionId)
		// if err != nil {
		// 	return nil, int64(0), err
		// }
		// var summaries []*DocumentSummary
		//summaries = append(summaries, ds)
		return summaries, cnt, err
	}

	//caOnly, ok := os.LookupEnv("CA_ONLY")
	// Use GPI to determine if need to search ca only
	//f.PatientGPI = "12724066"
	if f.PatientGPI == "" {
		//if ok && strings.ToLower(caOnly) == "true" {
		log.Debugf("@  Handling CA Documents ONLY\n")
		page := f.Page
		f.Source = "ca"

		log.Debugf("Search CA Only: %s\n", spew.Sdump(f))
		err := f.GetCADocs()
		if err != nil {
			log.Errorf("GetCaDocs returned error %v\n", err)
		}
		f.Page = page
		f.PatientGPI = f.PatientID
		//f.PatientID = f.PatientGPI
		f.makeQueryFilter()
		//log.Debugf("@  New Query: %v\n", f.QueryFilter)

		// if ActiveConfig().Env("ca_only") != "true" {
		// 	_, _, err = f.FindDocuments()
		// }

		docs, _ := f.QueryCache()
		total, err := f.CountCachedDocuments()

		return docs, total, err
	} else {
		log.Infof("looking for both fhir and CA patients gpi: [%s]", f.PatientGPI)
	}

	if f.TabID != "" {
		log.Debugf("\n\n@  Handling CA Documents\n")
		page := f.Page
		f.Source = "ca"
		//f.Source = ""
		err := f.GetCADocs()
		f.Page = page
		f.PatientID = f.PatientGPI
		f.makeQueryFilter()
		log.Debugf("@  New Query: %v\n", f.QueryFilter)

		_, _, err = f.FindDocuments()

		docs, _ := f.QueryCache()
		total, err := f.CountCachedDocuments()

		return docs, total, err
	}

	if f.MRN != "" {
		// pf := PatientFilter{MRN: f.MRN, Session: f.Session}

		// pats, err := pf.Search()
		// if err != nil {
		// 	log.Errorf("DocumentFilter looking for MRN: %s failed with err: %v\n", f.MRN, err)
		// 	return nil, 0, fmt.Errorf("404|Patient was not found for MRN: %s", pf.MRN)
		// }
		//f.PatientID = pats[0].EnterpriseID
		f.PatientGPI = f.PatientID
		f.MRN = ""

		//log.Debugf("@   documentFilter before Rebuild (Patient should be blank) : %v\n", f.queryString)

		err = f.makeQueryFilter()
		if err != nil {
			log.Errorf("!   makeQueryFilter err: %v\n", err)
		}

		log.Debugf("\n     @@documentFilter update: %v\n", f.queryString)
	}

	documents, total, err := f.FindDocuments()
	log.Infof("Search found %d documents total %d in %s\n", len(documents), total, time.Since(startTime))
	return documents, total, err
}

func (f *DocumentFilter) FindDocuments() ([]*fhir.DiagnosticReport, int64, error) {
	var docs []*fhir.DiagnosticReport
	var err error
	//log.Debugf("\n\nFindDocuments: Use Cache: %s\n", f.UseCache)
	startTime := time.Now()
	if strings.ToLower(f.ClearCache) == "true" {
		// TODO: Implement f.ClearCache and be sure to set it to false once done
	}

	f.Source = ActiveConfig().Source()
	if f.Page == 0 {
		log.Debugf("\n@@    FindDocuments Check for updates in FHIR\n")
		//f.UseCache = "false"

		err := f.CheckForFhirUpdates()
		log.Debugf("@  CheckForFhirDiagRepts returned err: %v\n", err)
		f.UseCache = "true"
		f.Page = 1
	}
	//if f.UseCache == "true" {
	log.Debugf("@@@     FindDocuments-263 Use Cache with filter: %v\n", f.QueryFilter)
	count, _ := f.CountCachedDocumentsFromSource()
	if count > 0 { // We have cached documents, Return them
		log.Debugf("@@@     FindDocuments 266 found %d Cache documents,\n", count)
		docs, _ = f.QueryCache() // Filter how they user requested from cache
		log.Infof("      %d found in FindDocuments-232 in %s\n\n", len(docs), time.Since(startTime))
		return docs, count, nil
	}

	//} else {
	log.Debugf("@  FindDocuments-273 is not using the Cache. Requesting from FHIR")
	//}

	fhirTime := time.Now()
	docs, err = f.FhirDiagRepts()
	log.Infof("@    Get Fhir Documents took %s  returning %d documents\n", time.Since(fhirTime), len(docs))
	// need to add query for fhir document references
	if err != nil { // NOne were found in FHIR
		return nil, 0, err
	}

	log.Debugf("\n\n@@@      Count Total documents=284 in Cache with filter: %v\n", f.QueryFilter)
	total, err := f.CountCachedDocuments()
	if err != nil {
		return nil, 0, err
	}
	//total, _ := CountCachedDocuments(f.QueryFilter)
	log.Debugf("@@@      %d found in FindDocuments-290 Count returned %d \n\n", len(docs), total)

	docs, _ = f.QueryCache() // Filter how they user requested from cache
	log.Infof("@@@      %d found in FindDocuments-293  QueryCache in %s\n\n", len(docs), time.Since(startTime))
	return docs, total, nil
}

func (f *DocumentFilter) CheckForFhirUpdates() error {

	log.Debugf("In CheckForFhirUpdates\n")
	latest, err := f.FindLatestCachedFhir()
	if err != nil {
		log.Errorf("Latest returned error: %v\n", err)
		return err // there are none
	}
	lastDate := fmt.Sprintf("$gt%s", latest.ReptDatetime.Format("2006-01-02 15:04:05"))
	query := fmt.Sprintf("patient=%s&created=%s&created=$le2500-12-31", f.PatientID, lastDate)
	log.Debugf("\n\n@    CheckForFhirUpdates query: %s\n", query)

	err = f.QueryFhirDiagRepts(query) // if any found they are in the cache
	return err
}

func (f *DocumentFilter) FindDocument() ([]*fhir.DiagnosticReport, error) {
	var docs []*fhir.DiagnosticReport
	var err error
	startTime := time.Now()

	if f.UseCache == "true" {
		log.Debugf("\n\n@@@ FindDocument Use Cache with filter: %v\n", f.QueryFilter)
		count, err := CountCachedDocuments(f.QueryFilter)
		if err != nil {
			log.Errorf("FindDocument: CountCache returned err: %v\n", err)
			return nil, err
		}
		if count > 0 { // Return what we have in the cache
			docs, _ = f.QueryCache() // Filter how they user requested from cache
			log.Infof("      %d found in FindDocument-320 in %s\n\n", len(docs), time.Since(startTime))
			return docs, nil
		}
	}

	doc, err := f.FhirDocument()
	if err != nil {
		log.Errorf("    FindDocument: FhirDocument returned: %v\n", err)
		return nil, err
	}

	docs = append(docs, doc)
	//docs, _ = f.QueryCache() // Filter how they user requested from cache
	log.Infof("      %d found in FindDocument-334 in %s\n\n", len(docs), time.Since(startTime))
	return docs, nil

}

// FhirDiagRepts returns the document by ID
// It searches Cached first then DocumentReferences and DiagnosticReport returning the first it finds
// Need to figure out how to handle both diagnostic and references

func (f *DocumentFilter) FhirDiagRepts() ([]*fhir.DocumentReference, error) {
	//c := config.Fhir()
	var ds []*fhir.DocumentReference
	log.Debugf("\n@    FhirDocument is searching FhirDidagnostics for Patient: %s \n", f.PatientID)

	startTime := time.Now()

	fds, err := f.FhirDocumentReferences()
	if err != nil {
		log.Errorf("   FhirDiagRepts: FhirDocumentReferences error: %v\n", err)
		return nil, err
	}
	log.Infof("$   FHIR Request took %s\n", time.Since(startTime))
	startTime = time.Now()
	ds = fromDocumentReferences(fds, f)
	if err != nil {
		log.Errorf("   fromDocumentReferences returned err: %v\n", err)
		return nil, err
	}
	log.Debugf("$   Create DocumentReferences took %s to process %d documents\n", time.Since(startTime), len(ds))

	if ds == nil {
		errMsg := fmt.Sprintf("FhirDiagRepts Returned nothing for Session: %s and patient: %s", f.SessionId, f.PatientID)
		log.Errorf("%s\n", errMsg)
		err = fmt.Errorf(errMsg)
		return nil, err
	}
	return ds, err
}

// QueryFhirDiagRepts Querys for FHIR documents using the provided query
// It searches Cached first then DocumentReferences and DiagnosticReport returning the first it finds
// Need to figure out how to handle both diagnostic and references

func (f *DocumentFilter) QueryFhirDiagRepts(query string) error {
	c := config.Fhir()
	var ds []*DocumentSummary
	log.Debugf("\n@    QueryFhirDocument is searching FhirDiagRepts using query: %s \n", query)
	startTime := time.Now()
	fds, err := c.GetDocumentReferences(query)
	if err != nil {
		return err
	}
	log.Infof("@   FHIR Request took %s\n", time.Since(startTime))
	startTime = time.Now()
	ds = fromDocumentReferences(fds, f)
	log.Infof("@   Create DocumentReferences took %s to process %d documents\n", time.Since(startTime), len(ds))
	return nil
}

// FhirDocument returns the document by ID
// It searches Cached first then DocumentReferences and DiagnosticReport returning the first it finds
// Need to figure out how to handle both diagnostic and references
func (f *DocumentFilter) FhirDocument() (*DocumentSummary, error) {
	//c := config.Fhir()
	var ds = new(DocumentSummary)

	//log.Debugf("\nFhirDocument-297 is searching All Fhir documents for Patient: %s - DocumentID: %s\n", f.PatientID, f.EnterpriseID)

	ref, err := f.FhirDocumentReference()
	if err != nil {
		log.Errorf("FhirDocument error: %v\n", err)
		return nil, err
	}
	err = ds.fromDocumentReference(ref, f.SessionId)
	if err != nil {
		log.Errorf("!   fromDocumentReference returned err: %v\n", err)
		return nil, err
	}

	if ds == nil {
		errMsg := fmt.Sprintf("FhirDocument Returned nothing for Session: %s and EnterpriseID: %s", f.SessionId, f.EnterpriseID)
		log.Errorf("%s\n", errMsg)
		err = fmt.Errorf(errMsg)
		return nil, err
	}

	return ds, err
}

//: Figure out how to passes what goroutine needs

func fromDocumentReferences(docR *fhir.DocumentReferences, f *DocumentFilter) []*DocumentSummary {
	session := f.Session
	cacheName := session.CacheName

	var docs []*DocumentSummary
	var nextLink string
	startTime := time.Now()
	for _, l := range docR.Link {
		if l.Relation == "next" {
			nextLink = l.URL
			go ProcessRemainingDocuments(nextLink, session)
			break
		}
	}
	for _, entry := range docR.Entry {
		var doc DocumentSummary
		for _, item := range entry.Resource.Content {
			//log.Debugf("Attachment: %v\n", item.Attachment.URL)
			if item.Attachment.ContentType == "application/pdf" {
				doc.FhirImageURL = item.Attachment.URL
				doc.makeImageURL()
			}
		}
		// 	doc.EffectiveDate = entry.Resource.ResourcePartial.EffectiveDateTime


		doc.Source = ActiveConfig().Source()
		doc.SourceType = "Reference"
		doc.Text = entry.Resource.Text
		//doc.Code = entry.Resource.Code.Text
		// 	doc.Category = entry.Resource.Category.Text
		doc.FullLink = entry.FullURL
		doc.EnterpriseID = entry.Resource.ID
		doc.Subject = entry.Resource.Subject
		//TODO: check both patient and subject for the patient information
		doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
		doc.PatientGPI = doc.PatientID
		doc.Performer = entry.Resource.Authenticator
		doc.ReptDatetime = entry.Resource.Created
		enc := strings.Split(entry.Resource.Context.EncounterNum.Reference, "/")
		//log.Debugf("\nEncounter: %v\n", enc)
		if len(enc) > 1 {
			doc.Encounter = enc[1]
		}
		doc.Description = entry.Resource.Description
		doc.Category = entry.Resource.Type.Text
		// if item.Attachment.ContentType == "application/pdf" {
		// 	doc.FhirImageURL = item.Attachment.URL
		// 	doc.makeImageURL()
		// }
		doc.SessionID = cacheName
		(&doc).Insert()
		docs = append(docs, &doc)
		//log.Debugf("\nFillEntry returing\n")
	}
	log.Infof("\n#    Cached %d documents in %s\n", len(docs), time.Since(startTime))
	return docs

}

func (f *DocumentFilter) FindFhirDiagnosticDocuments() ([]*DocumentSummary, error) {
	c := config.Fhir()
	log.Debugf("\nFHirDiagnosticDocuments is searching DiagnosticDocuments fo Patientr: %s\n", f.PatientID)

	//start := time.Now()
	q := fmt.Sprintf("patient=%s", f.PatientID)
	diag, err := c.FindDiagnosticReports(q)
	if err != nil {
		return nil, err
	}
	ds := fillDiagnosticReports(diag, f.SessionId)
	log.Debugf("%d documents returned from Diagnostic:505\n", len(ds))

	// ref, err := f.FhirDocumentReferences()
	// if err != nil {
	// 	log.Errorf("   FhirDocumentReference error: %s\n", err)
	// 	return nil, err
	// }
	// ds = fromDocumentReferences(ref, f)
	// log.Debugf("%d documents returned from References\n", len(ds))
	return ds, err
}

func (f *DocumentFilter) GetFhirDiagnosticDocument() (*fhir.DiagnosticReportResponse, error) {
	c := config.Fhir()

	//start := time.Now()
	//q := fmt.Sprintf("patient=%s", f.PatientID)
	qry := f.MakeFhirQuery()
	//log.Debugf("GetFhirDiagnosticDocument is searching : %s\n", spew.Sdump(f))
	diag, err := c.FindDiagnosticReports(qry)
	if err != nil {
		msg := fmt.Sprintf("FindDiagnosticReports error: %s\n", err.Error())
		log.Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	// ds := fillDiagnosticReports(diag, f.SessionId)
	// log.Debugf("%d documents returned from  Diagnostic:530\n", len(ds))

	// ref, err := f.FhirDocumentReferences()
	// if err != nil {
	// 	return nil, err
	// }
	// ds = fromDocumentReferences(ref, f)
	// log.Debugf("%d documents returned from References\n", len(ds))
	// query the cache for the requested documents.
	// ds, err = f.QueryCache()
	// log.Debugf("%d documents returned from query\n", len(ds))

	return diag, err
}

func (f *DocumentFilter) FhirDocumentReferences() (*fhir.DocumentReferences, error) {
	c := config.Fhir()
	log.Debugf("\n@  FhirDocumentReferences is searching using queryString: %s  queryFilter: %v\n", f.queryString, f.QueryFilter)
	if f.queryMap["patient"] == "" {
		err := fmt.Errorf("Patient id required for Document search.  Patient and other criteria is valid")
		return nil, err
	}
	//start := time.Now()

	//q := fmt.Sprintf("patient=%s", f.PatientGPI)
	q := f.queryString
	log.Debugf("@  query: %v\n", q)
	//q := f.queryString

	dRef, err := c.GetDocumentReferences(q)

	//log.Infof("FhirDiagnosticDocuments took %s\n", time.Since(start))
	if err != nil {
		log.Errorf("       FhirDocumentReferences returned err: %v\n", err)
		return nil, err
	}
	return dRef, nil
}

func (f *DocumentFilter) FhirDocumentReference() (*fhir.DocumentReference, error) {
	c := config.Fhir()
	log.Debugf("\nFhirDocumentReference is searching for ID: %s using queryString: %s  queryFilter: %v\n", f.EnterpriseID, f.queryString, f.QueryFilter)

	dRef, err := c.GetDocumentReference(f.ID)

	//log.Infof("FhirDiagnosticDocuments took %s\n", time.Since(start))
	if err != nil {
		log.Errorf("      FhirDiagRepts returned err: %v\n", err)
		return nil, err
	}
	return dRef, nil
}

func (d *DocumentSummary) Insert() error {
	var insertResult *mongo.InsertOneResult
	var err error
	d.setDates()
	collection, _ := storage.GetCollection("documents")
	//ctx := context.Background()
	insertResult, err = collection.InsertOne(context.TODO(), d)

	if err == nil {
		d.ID = insertResult.InsertedID.(primitive.ObjectID)
	} else {

		//log.Debugf(" Insert error type: %T : Spew:  %s\n", err, spew.Sdump(err))
	}

	return err
}

func (d *DocumentSummary) insertWithSession(session *mongo.Session) error {
	var insertResult *mongo.InsertOneResult
	var err error
	//d.SetDates()
	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()
	if err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
		if insertResult, err = collection.InsertOne(ctx, d); err != nil {
			// Ignore errors for now
			//  Need to only ignore if exists report others
			return err
		}
		return nil
	}); err != nil {
		//log.Errorf("Insert with session failed: %v\n", err)
		return nil
	}

	d.ID = insertResult.InsertedID.(primitive.ObjectID)

	return nil
}

func (p *DocumentSummary) setDates() {
	t := time.Now()
	p.CreatedAt = &t
	p.UpdatedAt = &t
	p.AccessedAt = &t
}

func (doc *DocumentSummary) fromDocumentReference(docR *fhir.DocumentReference, sessionID string) error {
	doc.Source = ActiveConfig().Source()
	doc.SourceType = "Reference"
	doc.Text = docR.Text
	//doc.FullLink = docR.FullURL
	doc.EnterpriseID = docR.ID
	doc.Subject = docR.Subject
	doc.Meta = docR.Meta
	//TODO: check both patient and subject for the patient information
	doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
	doc.Performer = docR.Authenticator
	doc.ReptDatetime = docR.Created
	enc := strings.Split(docR.Context.EncounterNum.Reference, "/")

	if len(enc) > 1 {
		doc.Encounter = enc[1]
	}
	doc.Description = docR.Description
	doc.Category = docR.Type.Text
	for _, c := range docR.Content {
		if c.Attachment.ContentType == "application/pdf" {
			doc.FhirImageURL = c.Attachment.URL
			doc.makeImageURL()
		}
	}
	doc.SessionID = sessionID
	doc.Insert()
	return nil
}

func (f *DocumentSummary) makeImageURL() {
	config := ActiveConfig()
	fhirURL := config.ImageURL()

	f.ImageURL = fmt.Sprintf("%s%s", fhirURL, f.EnterpriseID)
	return

}

func fillDocumentSummary(diag *fhir.DiagnosticReport) *DocumentSummary {
	// *DocumentSummary {
	//fmt.Println("\n\n================diag===")

	var doc DocumentSummary
	return &doc
}

func fillDiagnosticReports(diags *fhir.DiagnosticReportResponse, cacheName string) []*DocumentSummary {
	var docs []*DocumentSummary
	for _, entry := range diags.Entry {
		if entry.DiagnosticReport.FullURL == "" {
			entry.DiagnosticReport.FullURL = entry.FullURL
		}
		var doc DocumentSummary
		rpt := entry.DiagnosticReport
		doc.ReptDatetime = entry.DiagnosticReport.EffectiveDateTime
		doc.Meta = rpt.Meta
	//	fmt.Printf("\n\n\n\n\n\n###Meta : %s\n", spew.Sdump(rpt.Meta))
		doc.Source = ActiveConfig().Source()
		doc.SourceType = "Diagnostic"
		doc.Text = rpt.Text
		doc.Code = rpt.Code.Text
		doc.Category = rpt.Category.Text
		doc.FullLink = entry.FullURL
		doc.EnterpriseID = rpt.ID
		doc.Subject = rpt.Subject
		//TODO: check both patient and subject for the patient information
		doc.PatientID = rpt.ID
		doc.Performer = rpt.Performer
		enc := strings.Split(rpt.Encounter.Reference, "/")
		if len(enc) > 1 {
			doc.Encounter = enc[1]
		}
		var pf PresentedForm
		for _, v := range rpt.PresentedForm {
			pf.ContentType = v.ContentType
			pf.URL = v.URL
			doc.Images = append(doc.Images, pf)
		}
		for _, attachment := range rpt.PresentedForm {
			if attachment.ContentType == "application/pdf" {
				doc.FhirImageURL = attachment.URL
				doc.makeImageURL()
			}
		}
		//doc.Images = entry.Resource.PresentedForm
		doc.SessionID = cacheName
		(&doc).Insert()
		docs = append(docs, &doc)
	}
	return docs
}

func (ds *DocumentSummary) fillDiagnosticReport(fdoc *fhir.DiagnosticReport, cacheName string) error {
	var doc DocumentSummary
	// for _, attachment := range entry.Resource.PresentedForm {
	// 	if attachment.ContentType == "application/pdf" {
	// 		doc.FhirImageURL = attachment.URL
	// 		doc.makeImageURL()
	// 	}
	// }

	doc.ReptDatetime = fdoc.EffectiveDateTime

	doc.Source = ActiveConfig().Source()
	doc.SourceType = "Diagnostic"
	doc.Text = fdoc.Text
	doc.Code = fdoc.Code.Text
	doc.Category = fdoc.Category.Text
	doc.FullLink = fdoc.FullURL
	//doc.EnterpriseID = fdoc.ID
	//doc.Subject = fdoc.Subject
	//TODO: check both patient and subject for the patient information
	doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
	//doc.Performer = fdoc.Performer
	enc := strings.Split(fdoc.Encounter.Reference, "/")
	if len(enc) > 1 {
		doc.Encounter = enc[1]
	}
	var pf PresentedForm
	for _, v := range fdoc.PresentedForm {
		pf.ContentType = v.ContentType
		pf.URL = v.URL
		doc.Images = append(doc.Images, pf)
	}
	for _, attachment := range fdoc.PresentedForm {
		if attachment.ContentType == "application/pdf" {
			doc.FhirImageURL = attachment.URL
			doc.makeImageURL()
		}
	}
	//doc.Images = entry.Resource.PresentedForm
	doc.SessionID = cacheName
	(&doc).Insert()
	ds = &doc
	return nil
}

func mapToQueryString(q map[string]string) string {
	var query string
	for k, v := range q {
		//log.Debugf("%s=%s\n", k, v)
		s := fmt.Sprintf("%s=%s", k, v)
		if query == "" {
			query = query + s

		} else {
			query = query + fmt.Sprintf("&%s", s)
		}

	}
	//log.Debugf("    Result Query: %s\n", query)
	return query
}

//func GetDocumentImage(docID string) (*DocumentSummary, error) { //(*[]byte, error)
func (d *DocumentSummary) GetDocumentImage() error { //(*[]byte, error) {

	log.Debugf("GetDocumentImage DocID: %s\n", d.EnterpriseID)
	// ss := strings.Split(url, "/")
	// docID := ss[len(ss)-1]
	// log.Debugf("ss: %v\n", ss)
	filter := bson.M{"enterprise_id": d.EnterpriseID}
	err := d.CachedDocument(filter)
	if err != nil {

		log.Errorf("GetDocumentImage: [%s] GetCache error %v\n", d.EnterpriseID, err)
		return err
	}

	if d.Image != "" {
		log.Debugf("   Using cached image\n")
		return nil // We have an image do not get from cerner
	}
	url := d.FhirImageURL
	log.Debugf("Full Url: %s\n", url)
	docImage, err := readImageFromURL(url)
	if err != nil {
		log.Errorf("GetDocumentImage Error: %v\n", err)
		return err
	}
	//doc.Image = image
	d.Image = docImage.Content
	d.UpdateImage()

	return nil
}

func readImageFromURL(url string) (*DocumentImage, error) {
	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json+fhir")
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	log.Debugf("Status Code: %d - %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
		return nil, err
	}
	var image DocumentImage

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err := json.Unmarshal(respByte, &image); err != nil {
		return nil, err
	}
	return &image, nil
	//return &buf, err
}

func (rec *DocumentSummary) UpdateImage() {
	filter := bson.M{"_id": rec.ID}
	//log.Debugf("@@ Updater: %v\n", filter)

	update := bson.M{"$set": bson.M{"image": rec.Image}}
	collection, _ := storage.GetCollection("documents")
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Errorf(" Update error ignored: %s\n\n", err)
	}
	log.Debugf("Matched: %d  -- modified: %d\n", res.MatchedCount, res.ModifiedCount)

	return
}

//Get Chached entries first. If none, get them from FHIR
func CountCachedDocuments(filter bson.M) (int64, error) {
	log.Debugf("Document CountCachedDocuments matching: %v\n", filter)
	c, err := storage.GetCollection("documents")
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

//Get Chached entries first. If none, get them from FHIR
func (f *DocumentFilter) CountCachedDocuments() (int64, error) {
	filter := f.QueryFilter
	log.Debugf("Document CountCachedDocuments matching: %v\n", filter)
	c, err := storage.GetCollection("documents")
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("  Count returned error: %v\n", err)
		return 0, err
	}
	return count, nil
}

//Get Chached entries first. If none, get them from FHIR
func (f *DocumentFilter) CountCachedDocumentsFromSource() (int64, error) {

	mq := append(f.QueryFilterBase, bson.M{"source": f.Source})
	filter := bson.M{"$and": mq}
	log.Debugf("@   Document CountCachedDocumentsForSource matching: %v\n", filter)
	c, err := storage.GetCollection("documents")
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("  Count returned error: %v\n", err)
		return 0, err
	}
	return count, nil
}

//Get Chached entries first. If none, get them from FHIR
func GetCachedDocuments(filter bson.M) ([]*DocumentSummary, error) {
	log.Debugf("Document GetCachedDocuments using %v\n", filter)
	var documents []*DocumentSummary
	collection, err := storage.GetCollection("documents")
	if err != nil {
		//log.Debugf(" Error getting Collection: %s\n", err)
		return nil, err
	}
	//var encounter = new(Encounter)
	//log.Debugf("\nGetCacheDocuments using Filter: %v\n", filter)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		//log.Debugf("GetCachedDocuments search for %s returned error: %v\n", filter, err)
		//cursor.Close(context.TODO())
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var document DocumentSummary
		err = cursor.Decode(&document)
		if err != nil {
			cursor.Close(context.TODO())
			return nil, err
		}
		//.Dump(documents)
		documents = append(documents, &document)
	}
	if documents == nil {
		err = fmt.Errorf("404|no documents found for %s", filter)
	}
	return documents, err
}

//Get Chached entries first. If none, get them from FHIR
func GetCachedDocumentSummary(filter bson.M) (*DocumentSummary, error) {
	log.Debugf("DocumentGetCachedDocument using: %v\n", filter)
	var document DocumentSummary
	collection, err := storage.GetCollection("documents")
	if err != nil {

		return nil, err
	}
	log.Debugf("\nGetCachedDocument Checking for DocumentSummary with query: %v\n", filter)
	err = collection.FindOne(context.TODO(), filter).Decode(&document)
	if err != nil {
		log.Errorf("GetCachedDocumentSummary-FindOne returned err: %v\n", err)
		return nil, err
	}
	return &document, nil
}

//Get Chached entries first. If none, get them from FHIR
func (d *DocumentSummary) CachedDocument(filter bson.M) error {
	log.Debugf("DocumentGetCachedDocument using: %v\n", filter)

	collection, err := storage.GetCollection("documents")
	if err != nil {

		return err
	}
	log.Debugf("\nGetCachedDocument Checking for DocumentSummary with query: %v\n", filter)
	err = collection.FindOne(context.TODO(), filter).Decode(d)
	if err != nil {
		log.Errorf("CachedDocument-FindOne returned err: %v\n", err)
		return err
	}
	return nil
}

func (f *DocumentFilter) QueryCacheByEncounter() (bson.D, error) {
	if f.Encounter == "" {
		return nil, fmt.Errorf("query_by_encounter has no encounter")
	}
	//q := bson.D{{"patient", f.PatientID}, {"encounter", f.Encounter}}
	q := bson.D{{"sessionid", f.SessionId}, {"patient", f.PatientGPI}, {"encounter", f.Encounter}}
	return q, nil
}

func (f *DocumentFilter) QueryCache() ([]*DocumentSummary, error) {
	var limit int64 = 20
	var skip int64 = 0

	if f.Limit > 0 {
		limit = f.Limit
	}
	if f.Page > 0 {
		skip = (f.Page - 1) * limit
	}
	if f.Skip > 0 {
		skip = f.Skip
	}
	// q, err := f.QueryCacheByEncounter()

	log.Debugf("@@@   995-Document QueryCache Using filter %s \n", f.QueryFilter)
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	// Multi sort fields are separated by [, ].
	// Order is based upon the order of the field names
	// sortFields := strings.Split(f.SortBy,", ")
	// for i, f := range sortFields {

	const DESC = -1
	const ASC = 1
	var sortFields bson.D
	var documents []*DocumentSummary
	order := ASC // Default Assending

	if strings.ToLower(f.Order) == "desc" {
		order = DESC
	}

	sort := bson.E{}
	if f.Column == "" {
		f.Column = "rept_datetime"
	}
	sort = bson.E{f.Column, order}
	sortFields = append(sortFields, sort)
	if len(f.SortBy) > 0 {
		for _, s := range f.SortBy {
			if s == "visit_num" {
				sort = bson.E{"encounter", order}
			} else {
				sort = bson.E{s, order}
			}
			sortFields = append(sortFields, sort)
		}
	}
	findOptions.SetSort(sortFields)
	log.Debugf("@   1045 --  sort: %v", sortFields)

	// }

	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, f.QueryFilter, findOptions)
	if err != nil {
		log.Debugf("QueryCache for %s returned error: %v\n", f.QueryFilter, err)
		//cursor.Close(ctx)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			//log.Debugf("document queryCache error while closing cursor: %v\n", err)
			log.WithError(err).Warn("Got error while closing cursor")
		}
	}()
	//log.Debugf("\n    No Error on QueryCache\n\n")
	for cursor.Next(ctx) {
		var document DocumentSummary
		err = cursor.Decode(&document)
		if err != nil {
			//log.Debugf("   Next error: %v\n", err)
			//cursor.Close(context.TODO())
			log.WithError(err).Warn("Got error while closing cursor")
			return nil, err
		}
		//log.Debugf("  Added one to documents\n")
		documents = append(documents, &document)
	}
	log.Debugf("   Finished fetching documents from cursor\n\n")
	//cursor.Close(context.TODO())
	if len(documents) == 0 {
		err = fmt.Errorf("404|no documents found for %s", f.QueryFilter)
		log.Error(err)
	} else {
		log.Debugf("QueryCache found %d documents \n", len(documents))
	}
	return documents, err
}

func (rec *DocumentSummary) UpdateAccess() time.Time {
	filter := bson.M{"_id": rec.ID}
	//log.Debugf("@@ Updater: %v\n", filter)
	loc, _ := time.LoadLocation("UTC")
	accessed := time.Now().In(loc)
	update := bson.M{"$set": bson.M{"accessedat": accessed}}
	collection, _ := storage.GetCollection("documents")
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Errorf(" Update error ignored: %s\n\n", err)
	}
	log.Debugf("Matched: %d  -- modified: %d\n", res.MatchedCount, res.ModifiedCount)
	(*rec).AccessedAt = &accessed
	return accessed
}

func (f *DocumentFilter) makeQueryMap() error {

	log.Debugf("@@@Document Using SessionID: %s \n\n", f.SessionId)

	patientID := strings.Trim(f.PatientID, " ")
	patientGPI := strings.Trim(f.PatientGPI, " ")
	mrn := strings.Trim(f.MRN, " ")
	encounterID := strings.Trim(f.Encounter, " ")
	id := strings.Trim(f.ID, " ")
	enterpriseID := strings.Trim(f.EnterpriseID, " ")
	category := strings.Trim(f.Category, " ")
	reptDatetime := strings.Trim(f.ReptDatetime, " ")
	count := strings.Trim(f.Count, " ")
	UseCache := strings.Trim(f.UseCache, "")

	sessionid := strings.Trim(f.SessionId, " ")
	config := ActiveConfig()
	m := make(map[string]string)
	if patientID == "" {
		log.Error("PatientID is a required filter on all DiagnosticReport searches")
		msg := "400|PatientID is a required filter on all DiagnosticReport searches"
		log.Errorf(msg)
		err := fmt.Errorf(msg)
		return err
	}

	if sessionid != "" {
		m["sessionid"] = sessionid
	}
	if patientID != "" {
		m["patient"] = patientID
	}
	if patientGPI != "" {
		m["patientgpi"] = patientGPI
	} else {
		//m["patientgpi"] = m["patient"]
	}
	if count != "" {
		m["count"] = count
	} else {
		count = config.RecordLimit()
		log.Debugf("     @@@ Setting Count from config: %s\n", count)
		f.Count = count
		m["count"] = count
	}
	if mrn != "" {

		m["mrn"] = mrn
	}
	if encounterID != "" {
		m["encounter"] = encounterID
	}
	if enterpriseID != "" {
		m["id"] = enterpriseID
	}
	if id != "" {
		m["id"] = id
	}
	if category != "" {
		m["category"] = category
	}
	if reptDatetime != "" {
		m["rept_datetime"] = reptDatetime
	}
	if UseCache == "" {
		f.UseCache = "true"
	} else {
		f.UseCache = "false"
	}
	// log.Debugf("\n\nNew map: \n")

	f.queryMap = m
	log.Debugf("Using SessionID: %s \n\n", f.SessionId)

	return nil
}

func (f *DocumentFilter) makeQueryFilter() error {
	log.Debugf("Document makeQueryFilter:1162 \n\n")
	err := f.makeQueryMap()
	if err != nil {
		fmt.Printf("!@!@!@!    makeQueryMap Failed: %s\n", err.Error())
		return err
	}
	f.makeQueryString()
	//f.QueryFilter, _ = com.FilterFromMap(f.queryMap)
	log.Debugf("\n\n@#@#@#   makeQueryFilter:1169 working with: %v\n\n", f.queryMap)
	layout := "2006-01-02"
	mq := []bson.M{}
	for k := range f.queryMap {
		val := f.queryMap[k]
		//log.Debugf("k: %s,  v: %s\n", k, q[k])
		if k == "count" {
			continue
		}
		// log.Debugf("\n    Document makeQueryFilter  Current Session: \n\n")

		if k == "sessionid" {
			log.Debugf("      Adding seleted sessionID: %s\n", f.SessionId)
			mq = append(mq, bson.M{"sessionid": f.SessionId})
		} else if k == "id" {
			//log.Debugf("Converting search for id %s to search for enterprise_id\n", val)
			mq = append(mq, bson.M{"enterprise_id": val})
		} else if k == "category" {
			mq = append(mq, bson.M{"category": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
			// } else if k == "patient" {
			// 	mq = append(mq, bson.M{"patient": val})
		} else if k == "effectivedate" {
			s := strings.Split(val, "|")
			condition := s[0]
			input := s[1]
			useDate, err := time.Parse(layout, input)
			if err != nil {
				log.Errorf("makeQueryFilter-Effective Date Error: %v\n", err)
				continue
			}
			q := bson.M{"effectivedate": bson.M{condition: useDate}}
			mq = append(mq, q)
		} else if k == "patient" {
			continue
		} else {
			log.Debugf("    Adding %s to QueryFilter\n", k)
			mq = append(mq, bson.M{k: val})
		}
	}
	//f.QueryFilter = bson.M{}
	if len(mq) > 0 {
		f.QueryFilterBase = mq
		f.QueryFilter = bson.M{"$and": mq}
	}
	log.Debugf("    Final Document QueryFilter:1214 %s\n", spew.Sdump(f.QueryFilter))
	return nil
}

func (f *DocumentFilter) makeQueryString() {
	fmt.Printf("\n\n\n       ### queryMap:1219 %v\n", f.queryMap)
	f.queryString = fmt.Sprintf("?%s=%s", "patient", f.queryMap["patient"])
	fmt.Printf("\n\n#### MakeQueryString:1221 f.QueryString: %s\n", f.queryString)
	if f.Count != "" {
		f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, "_count", f.queryMap["count"])
	}
	//f.queryString = ""
	// for k := range f.queryMap {
	// 	f.queryString = fmt.Sprintf("%s=%s", "patient", f.queryMap["patient"])
	// if f.queryString == "" {
	// 	f.queryString = fmt.Sprintf("%s=%s", k, f.queryMap[k])
	// } else {
	// 	f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, k, f.queryMap[k])
	// }
	// }
}

func ConvertDocumentsToCA(docs []*DocumentSummary) []*CADocument {
	startTime := time.Now()
	var caDocuments []*CADocument
	log.Debugf("Converting %d documents to CA documents\n", len(docs))
	// for _, d := range docs {
	// 	doc := d.ToCA()
	// 	caDocuments = append(caDocuments, doc)
	// }

	log.Infof("Convert %d documents to CA took %s\n", len(caDocuments), time.Since(startTime))
	return caDocuments
}

func (d *DocumentSummary) ByEnterpriseID() error {
	filter := bson.M{"enterprise_id": d.EnterpriseID}
	d, err := GetCachedDocumentSummary(filter)
	if err != nil { //Cache not found for ths Document
		return err
	}
	return nil
}


// func (f *fhir.DocumentReference) ToCA() *CADocument {
// 	var caDoc CADocument
// 	caDoc.PatientGPI = f.ID

// 	caDoc.PatientGPI = d.PatientGPI
// 	caDoc.Encounter = d.Encounter
// 	rpdt := d.ReptDatetime
// 	caDoc.ReptDateTime = &rpdt
// 	caDoc.Description = d.Description
// 	caDoc.Subtitle = d.Category
// 	caDoc.Text = d.Text
// 	caDoc.VersionID = d.Versionid

// 	caDoc.Source = d.SourceType

// 	if caDoc.Description == "" {
// 		caDoc.Description = d.Category
// 	}
// 	config := ActiveConfig()

// 	switch d.Source {
// 	case "cerner":
// 		caDoc.ImageURL = fmt.Sprintf("%s%s", config.ImageURL(), caDoc.ID)
// 	case "ca":
// 		caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
// 	case "QC":
// 		caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
// 	case "HPF":
// 		caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
// 	}

// 	//log.Debugf("\n@  ImageURL- %: %s\n", caDoc.Source, caDoc.ImageURL)
// 	// if caDoc.Source == "cerner" {
// 	// 	caDoc.ImageURL = fmt.Sprintf("%s%s", config.ImageURL(), caDoc.ID)
// 	// }
// 	return &caDoc
// }

func ConvertDocumentsToVS(docs []*DocumentSummary) {
	startTime := time.Now()
	for _, d := range docs {
		d.ImageURL = fmt.Sprintf("%s%s", config.ImageURL(), d.EnterpriseID)
		// doc := d.ToCA()
		// caDocuments = append(caDocuments, doc)
	}
	log.Infof("Convert %d documents to VS took %s\n", len(docs), time.Since(startTime))

}

func DeleteDocuments(sessionid string) {

	startTime := time.Now()
	log.Debugf("@@@!!!  Deleting Documents for session %s\n", sessionid)
	collection, _ := storage.GetCollection("documents")
	filter := bson.D{{"sessionid", sessionid}}
	log.Debugf("    bson filter delete: %v\n", filter)
	deleteResult, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Errorf("!     DeleteDocuments for session %s failed: %vn", sessionid, err)
		return
	}
	log.Infof("@@@!!!     Deleted %v Documents for session: %v in %s\n", deleteResult.DeletedCount, sessionid, time.Since(startTime))
}

// ProcessRemainingDocuments reads the next set of documents and if there is another link spins up
// another gorouting to process it while this one is filling the cache.

func ProcessRemainingDocuments(link string, session AuthSession) {

	err := ReadNextSet(link, &session)
	if err != nil {
		log.Errorf("!!!   ReadNextSet returned error: %v", err)
	}
	//ds := BuildDocumentReferences(ref, &session)
}

func ReadNextSet(url string, session *AuthSession) error {
	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"

	startTime := time.Now()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json+fhir")
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	log.Debugf("Status Code: %d - %s", resp.StatusCode, resp.Status)
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
		return err
	}
	var ref fhir.DocumentReferences

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err := json.Unmarshal(respByte, &ref); err != nil {
		return err
	}
	log.Infof("$  FHIR Request set of documents took %s\n", time.Since(startTime))
	err = BuildDocumentReferences(&ref, session)
	return err
	//return &buf, err
}

/////////////////////////////////////////////////////////////////////////////////////////
//                                 FHIR Getters                                         /
/////////////////////////////////////////////////////////////////////////////////////////

func GetImage(fd *fhir.DiagnosticReport, imageType string) string {
	for _, form := range fd.PresentedForm {
		switch form.ContentType {
		case imageType:
			return form.URL
		}
	}
	return ""
}

/////////////////////////////////////////////////////////////////////////////////////////
//                                  Convert to CA                                       /
/////////////////////////////////////////////////////////////////////////////////////////

func FhirDiagDocsToCA(fds []*fhir.DiagnosticReport) []*CADocument{
	caDocuments := []*CADocument{}
	for _, d := range fds {
		doc := FhirDiagDocToCA(d)
		caDocuments = append(caDocuments, doc)
	}
	return caDocuments
}

func FhirDiagDocToCA(fd* fhir.DiagnosticReport) *CADocument {
	var caDoc CADocument
	var err error
	caDoc.PatientGPI = fd.ID

	caDoc.PatientGPI = fd.ID
	caDoc.VersionID, err = strconv.ParseUint(fd.Meta.VersionID, 10, 64)
	if err != nil {
		log.Errorf("Invalid VersionID: [%s] error:%s\n", fd.Meta.VersionID, err.Error())
	}
	enc := strings.Split(fd.Encounter.Reference, "/")
	if len(enc) > 1 {
		caDoc.Encounter =enc[1]
	}			// Encounter is not available
	caDoc.Repository = "FHIR"
	rpdt := fd.EffectiveDateTime
	caDoc.ReptDateTime = &rpdt
	caDoc.ImageURL = GetImage(fd, "application/pdf")
	caDoc.Pages = 0            // Unavailable
	caDoc.Subtitle = fd.Code.Text
	caDoc.Text = fd.Text
	cfg := ActiveConfig()
	caDoc.Source = strings.ToLower(cfg.Source())

	if caDoc.Description == "" {
		caDoc.Description = fd.Code.Text
	}
	// config := ActiveConfig()


	// switch source {
	// case "cerner":
	// 	caDoc.ImageURL = fmt.Sprintf("%s%s", config.ImageURL(), caDoc.ID)
	// case "ca":
	// 	caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
	// case "QC":
	// 	caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
	// case "HPF":
	// 	caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
	// }
	return &caDoc
}

func BuildDocumentReferences(docR *fhir.DocumentReferences, sessionp *AuthSession) error {
	// *DocumentSummary {
	// fmt.Println("\n\n================docR===")

	session := *sessionp
	cacheName := session.CacheName
	dbSession, err := storage.GetSession()
	if err != nil {
		return err
	}
	//var docs []*DocumentSummary
	var nextLink string
	startTime := time.Now()
	for _, l := range docR.Link {
		if l.Relation == "next" {
			nextLink = l.URL
			//log.Debugf("NextLink: %s\n", nextLink)
			go ProcessRemainingDocuments(nextLink, session)
			break
		}
	}
	numDocs := 0
	for _, entry := range docR.Entry {
		var doc DocumentSummary

		doc.Source = ActiveConfig().Source()
		doc.SourceType = "Reference"
		doc.Text = entry.Resource.Text
		//doc.Code = entry.Resource.Code.Text
		// 	doc.Category = entry.Resource.Category.Text
		doc.ReptDatetime = entry.Resource.Created
		doc.FullLink = entry.FullURL
		doc.EnterpriseID = entry.Resource.ID
		doc.Subject = entry.Resource.Subject
		//TODO: check both patient and subject for the patient information
		doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
		doc.PatientGPI = doc.PatientID
		doc.Performer = entry.Resource.Authenticator
		enc := strings.Split(entry.Resource.Context.EncounterNum.Reference, "/")
		if len(enc) > 1 {
			doc.Encounter = enc[1]
		}
		doc.Description = entry.Resource.Description
		doc.Category = entry.Resource.Type.Text
		for _, item := range entry.Resource.Content {
			//log.Debugf("Attachment: %v\n", item.Attachment.URL)
			if item.Attachment.ContentType == "application/pdf" {
				doc.FhirImageURL = item.Attachment.URL
				doc.makeImageURL()
			}
		}
		doc.SessionID = cacheName
		doc.setDates()
		(&doc).insertWithSession(&dbSession)
		//(&doc).Insert(cacheName, &dbSession)
		numDocs++
	}
	log.Infof("#    Cached %d documents in %s\n", numDocs, time.Since(startTime))
	return nil
}

func (f *DocumentFilter) FindLatestCachedFhir() (*DocumentSummary, error) {
	var document *DocumentSummary
	var limit int64 = 1
	var skip int64 = 0
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	// Multi sort fields are separated by [, ].
	// Order is based upon the order of the field names
	// sortFields := strings.Split(f.SortBy,", ")
	// for i, f := range sortFields {

	const DESC = -1
	const ASC = 1
	var sortFields bson.D
	var documents []*DocumentSummary
	order := DESC

	mq := append(f.QueryFilterBase, bson.M{"source": f.Source})
	filter := bson.M{"$and": mq}

	sort := bson.E{"rept_datetime", order}
	sortFields = append(sortFields, sort)
	findOptions.SetSort(sortFields)
	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Debugf("QueryCache for %s returned error: %v\n", f.QueryFilter, err)
		//cursor.Close(ctx)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Errorf("Got error while closing cursor")
		}
	}()
	for cursor.Next(ctx) {
		var document DocumentSummary
		err = cursor.Decode(&document)
		if err != nil {
			log.Warnf("Got error while closing cursor")
			return nil, err
		}
		documents = append(documents, &document)
	}
	//cursor.Close(context.TODO())
	if len(documents) == 0 {
		err = fmt.Errorf("404|no documents found for %s", f.QueryFilter)
		log.Error(err)
	} else {
		log.Debugf("QueryLatest found %d documents \n", len(documents))
		document = documents[0]
	}
	return document, err
}

func (f *DocumentFilter) GetCADocs() error {
	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"
	// config := ActiveConfig()
	// caURL := config.Env("caServerURL")
	config := ActiveConfig()
	startTime := time.Now()
	caURL := config.Env("caServerURL")
	//page := f.Page
	patID := f.PatientID
	tabID := f.TabID
	if config.Mode() == "dev" {
		log.Warn("DEV Mode forcing Tab to 1000")
		tabID = "1000"
	}
	column := f.Column
	//reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&page=%d&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, page, tabID)

	reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, tabID)
	f.Source = "ca"

	log.Debugf("\n\n\n#####Requesting from CA: %s#####\n\n", reqURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("Accept", "application/json")
	authID := f.SessionId
	log.Debugf("Auth: %s\n", authID)
	req.Header.Set("AUTHORIZATION", authID)
	resp, err := client.Do(req)
	//resp, err := http.Get(url)
	if err != nil {
		log.Errorf("CA Query error: %v\n", err)
		return err
	}

	defer resp.Body.Close()
	log.Debugf("Status Code: %d - %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
		return err
	}
	log.Debugf("Ready to unmarashal\n")
	//var ref fhir.DocumentReferences
	var docs []*DocumentSummary
	//var docSet CaResults

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()
	if err := json.Unmarshal(respByte, &docs); err != nil {
		return err
	}
	log.Infof("GetCADocs request took %s\n", time.Since(startTime))
	procTime := time.Now()
	//err = BuildDocumentReferences(&ref, session)
	//log.Debugf("GetCADocs: %s\n", spew.Sdump(docs))
	err = f.ProcessCADocuments(docs)
	log.Infof("GetCADocs process time: %s  total time: %s\n", time.Since(procTime), time.Since(startTime))
	//err = fmt.Errorf("Testing force err")
	//log.Infof("$  CA Request set of documents took %s\n", time.Since(startTime))
	return err
	//return &buf, err
}

func (f *DocumentFilter) ProcessCADocuments(docs []*DocumentSummary) error {
	//caURL := ActiveConfig().Env("caServerURL")NO
	var err error

	for _, doc := range docs {
		doc.SessionID = f.SessionId
		doc.EnterpriseID = strconv.FormatUint(doc.CADocumentID, 10) //strconv.Itoa(doc.CADocumentID)
		//doc.ImageURL =
		doc.SourceType = doc.Source
		doc.Source = "ca"
		//log.Debugf("@@       Received source %s and sourceType: %s from ca\n", doc.Source, doc.SourceType)
		//doc.Source = f.Source  // use the value returned from query
		doc.PatientID = f.PatientID
		//TODO: doc.PatientGPI = f.PatientGPI
		doc.PatientGPI = f.PatientID
		doc.Issued = doc.ReptDatetime
		doc.ImageURL = fmt.Sprintf("%s/%d", ActiveConfig().Env("caImageURL"), doc.Versionid)
		//doc.Versionid = doc.Versionid
		log.Debugf("Received Version: %d : %T\n", doc.Versionid, doc.Versionid)
		//log.Debugf("ProcessCADocs: %s\n",spew.Dump(doc))
		doc.Insert()
	}
	return err
}

func (df *DocumentFilter) MakeFhirQuery() string {
	qry := "patient=" + df.PatientGPI
	if df.BeginDate != "" {
		mdyDate, _ := common.MDYToFhir(df.BeginDate)
		bDate  := common.FhirDateToString(mdyDate, "full")
		qry = fmt.Sprintf("%s&date=ge%s", qry, bDate)
	}
	if df.EndDate != "" {
		mdyDate, _ := common.MDYToFhir(df.EndDate)
		eDate  := common.FhirDateToString(mdyDate, "full")
		qry = fmt.Sprintf("%s&date=lt%s", qry, eDate)
	}
	if df.Count != "" {
		qry = fmt.Sprintf("%s&Count=%s", qry, df.Count)
	}
	//fmt.Printf("\n\n####  Query: %s\n", qry)
	df.queryFhirString = qry
	return qry
}

// func (f *DocumentFilter)CachePages(links []fhir.Link) {
// 	url := NextPageLink(links)
// 	i := 1
// 	for {
// 		if url == "" {
// 			break
// 		}
// 		links, _ = f.ProcessPage(url)
// 		//fmt.Printf("Page: %d - %s\n", i, spew.Sdump(links))
// 		i = i+1
// 		url = NextPageLink(links)
// 	}
// }

// func (f *DocumentFilter)ProcessPage(url string) ([]fhir.Link, error) {
// 	//startTime := time.Now()
// 	fhirPatients, err := fhirC.NextFhirPatients(url)
// 	if err != nil {
// 		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
// 		return nil, err
// 	}
// 	//fmt.Printf("\n\nNext Page of results: %s\n\n", spew.Sdump(fhirPatients))
// 	patients : PatientsFromResults(fhirPatients)
// 	//patients := parsePatientResults(fhirPatients, f.SessionId)
// 	//spew.Dump()

// 	if len(patients) == 0 {
// 		return nil, fmt.Errorf("404|no more Documents found for %s", f.queryString)
// 	}
// return fhirPatients.Link, nil
// }


// func (f *PatientFilter) makeDocCacheQueryFilter() {
// 	f.makeQueryMap()
// 	f.makeQueryString()
// 	f.queryFilter, _ = com.FilterFromMap(f.queryMap)
// 	//fmt.Printf("\n\ntoFilter returning: %v\n\n", filter)
// 	layout := "2006-01-02"
// 	mq := []bson.M{}
// 	for k := range f.queryMap {
// 		val := f.queryMap[k]
// 		//fmt.Printf("k: %s,  v: %s\n", k, q[k])
// 		if k == "id" {
// 			//fmt.Printf("Converting search for id %s to search for enterpriseid\n", val)
// 			mq = append(mq, bson.M{"id": val})
// 		} else if k == "given" {
// 			mq = append(mq, bson.M{"given": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
// 		} else if k == "family" {
// 			mq = append(mq, bson.M{"family": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
// 		} else if k == "email" {
// 			mq = append(mq, bson.M{"email": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
// 		} else if k == "birthdate" {
// 			condition := "eq"
// 			input := ""
// 			s := strings.Split(val, "|")
// 			if len(s) > 1 {
// 				condition = s[0]
// 				input = s[1]
// 			} else {
// 				condition = "$eq"
// 				input = s[0]
// 			}
// 			useDate, err := time.Parse(layout, input)
// 			if err != nil {
// 				fmt.Printf("Error: %v\n", err)
// 				continue
// 			}

// 			q := bson.M{"birthdate": bson.M{condition: useDate}}
// 			mq = append(mq, q)

// 		} else {
// 			mq = append(mq, bson.M{k: val})
// 		}
// 	}

// 	if len(mq) > 0 {
// 		f.queryFilter = bson.M{"$and": mq}
// 	}
// 	fmt.Printf("@    queryFilter: %v\n", f.queryFilter)
// 	return
// }

// func (f *DocumentFilter) makeDocQueryMap() error {

// 	m := make(map[string]string)
// 	mrn := strings.Trim(f.MRN, " ")
// 	given := strings.Trim(f.Given, " ")
// 	family := strings.Trim(f.Family, " ")
// 	encounter := strings.Trim(f.EncounterID, " ")
// 	enterpriseID := strings.Trim(f.EnterpriseID, " ")
// 	email := strings.Trim(f.Email, " ")
// 	id := strings.Trim(f.ID, " ")
// 	birthdate := strings.Trim(f.BirthDate, " ")

// 	if family != "" {
// 		m["family"] = family
// 		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
// 		f.SortBy = append(f.SortBy, "given")
// 	}
// 	if given != "" {
// 		if family != "" {
// 			m["given"] = given

// 		} else {
// 			return fmt.Errorf("400|Invalid search: given alone is invalid")
// 		}
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
// 		m["id"] = id
// 	}
// 	if email != "" {
// 		m["email"] = email
// 	}
// 	if birthdate != "" {
// 		m["birthdate"] = birthdate
// 	}
// 	f.queryMap = m

// 	return nil
// }

// func (f *DocumentFilter) makeDocQueryString() {

// 	f.queryString = ""
// 	for k := range f.queryMap {

// 		if f.queryString == "" {
// 			f.queryString = fmt.Sprintf("%s=%s", k, f.queryMap[k])
// 		} else {
// 			f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, k, f.queryMap[k])
// 		}
// 	}
// 	if f.Count != "" {
// 		f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, "_count", f.queryMap["count"])
// 	}
// 	fmt.Printf("QueryString:758 %s\n", f.queryString)
//}
*/
