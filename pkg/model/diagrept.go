package model

import (
	"bytes"
	//"context"
	"encoding/json"
	//"errors"
	"fmt"

	//"os"
	//"strconv"

	"net/http"
	"strings"
	"time"

	//"github.com/dgrijalva/jwt-go"
	"github.com/davecgh/go-spew/spew"

	log "github.com/sirupsen/logrus"

	//"github.com/dhf0820/cernerFhir/pkg/storage"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Thing struct {
	Display   string `json:"display"`
	Reference string `json:"reference"`
}

// Person is a human
type Person Thing

type CaResults struct {
	Documents  []*DocumentSummary `json:"documents"`
	TotalPages uint               `json:"total_pages"`
}

type DocumentSummary struct {
	CacheID      primitive.ObjectID `json:"cache_id" bson:"_id,omitempty"`
	SessionID    string             `json:"-"`
	PatientID    string             `json:"patient_id" bson:"patient_id"`
	PatientGPI   string             `json:"patient_gpi" bson:"patient_GPI"`
	Encounter    string             `json:"visit_num" bson:"encounter"`
	Source       string             `json:"source" bson:"source"`
	FullLink     string             `json:"full_link" bson:"full_link"`
	Meta         fhir.MetaData      `json:"meta" bson:"meta"`
	Category     string             `json:"category" bson:"category"`
	Code         string             `json:"code" bson:"code"`
	Description  string             `json:"description" bson:"description"`
	Text         fhir.TextData      `json:"text" bson:"text"`
	ReptDatetime time.Time          `json:"rept_datetime" bson:"rept_datetime"`
	Issued       time.Time          `json:"issued" bson:"issued"`
	Subject      fhir.Person        `json:"subject" bson:"subject"`
	Performer    fhir.Person        `json:"performer" bson:"performer"`
	EnterpriseID string             `json:"documentid" bson:"enterprise_id"`
	CADocumentID uint64             `json:"document_id" bson:"document_id"`
	SourceType   string             `json:"sourceType" bson:"source_type"`
	CreatedAt    *time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt    *time.Time         `json:"updated_at" bson:"updated_at"`
	AccessedAt   *time.Time         `json:"accessed_at" bson:"accessed_at"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	FhirImageURL string             `json:"fhir_image_url" bson:"fhir_&image_url"`
	Images       []PresentedForm    `json:"images" bson:"images"`
	Image        string             `json:"image" bson:"image"` // base64
	Versionid    uint64             `json:"version_id" bson:"version_id"`
	TotalPages   int                `json:"total_pages"`
	//Image         primitive.Binary   `json:"image" bson:"image"`
}

type CADocument struct {
	CacheID      primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	SessionId    string             `json:"-" bson:"session_id"`
	ID           string             `json:"document_id" bson:"document_id"`
	VersionID    uint64             `json:"version_id" bson:"version_id"`
	Encounter    string             `json:"visit_num" bson:"visit_num"`
	Repository   string             `json:"repository" bson:"repository"`
	Category     string             `json:"category" bson:"category"`
	Class        string             `json:"class" bson:"class"`
	Source       string             `json:"source" bson:"source"`
	Description  string             `json:"description" bson:"description"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	Pages        int                `json:"pages" bson:"pages"`
	ReptDateTime *time.Time         `json:"rept_datetime" bson:"rept_datetime"`
	Subtitle     string             `json:"subtitle" bson:"subtitle"`
	PatientGPI   string             `json:"patient_gpi" bson:"patient_gpi"`
	Text         string             `json:"text" bson:"text"`
	Type         string             `json:"type" bson:"type"`
	DocStatus    string             `json:"doc_status" bson:"doc_status"`
	CreatedAt    *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt    *time.Time         `json:"-" bson:"updated_at"`
	AccessedAt   *time.Time         `json:"-" bson:"accessed_at"`
	//Text         fhir.TextData `json:"text"`
}

type PresentedForm struct {
	ContentType string `json:"context_type" bson:"context_type"`
	URL         string `json:"url" bson:"url"`
}

type DocumentImage struct {
	ResourceType string        `json:"resourceType" bson:"resource_type"`
	ID           string        `json:"id" bson:"id"`
	Meta         fhir.MetaData `json:"meta" bson:"meta"`
	ContentType  string        `json:"contentType" bson:"content_type"`
	Content      string        `json:"content" bson:"content"`
}

//var activeDocumentFilter *DocumentFilter

func (df *DocumentFilter) FindFhirDiagRepts() {
	fmt.Printf("\n####################FindFhirDiagReports is searching DiagnosticReports for Patient: %s#############\n\n", df.PatientGPI)
	_, _, totalInCache, _ := df.DocumentCacheStats()
	cacheStatus := df.Session.GetDiagReptStatus()
	//If the cacheStatus == done and there are some documents in cache we are done
	// if cacheStatus == done and there are documents in cache. we have the documents and do not restart caching
	if cacheStatus == "done" && totalInCache > 0 {
		fmt.Printf("FhirDiag is Done and there are documents")
		return
	}

	// If cache is not "done" it is building and do not restart the caching
	if cacheStatus != "done" {
		fmt.Printf("FhirDiag cacheStatus = %s we are Filling\n", cacheStatus)
		return
		// fhirDocs, _, _, _, _, err := df.GetFhirDocumentPage()
		// return fhirDocs, err
	}
	//if cacheStatus =="done" and nothing in cache, start caching

	df.CacheFhirDiagRepts()
	fmt.Printf("FindFhirDiagRepts:135 -- CacheFhirDiagREpts returned")
	return
}

func (df *DocumentFilter) CacheFhirDiagRepts() {
	fmt.Printf("\n##################### CacheFhirDiagRepts ################################\n")
	c := config.Fhir()
	//caching needs to be started.
	df.MakeFhirQuery()
	//start := time.Now()

	//q := fmt.Sprintf("patient=%s", f.fhirQuery)
	log.Debugf("CacheDiagRepts:147 -- c.FindDiagnosticReports is called with %s", df.fhirQuery)
	resp, err := c.FindDiagnosticReports(df.fhirQuery)

	if err != nil {
		return
	}
	fmt.Printf("\n\n\n###CacheDiagRepts153 -- calling FollowDiagNextLink: %s\\nn", spew.Sdump(resp.Link))
	go df.FollowDiagNextLinks(resp.Link)

	docs := []*fhir.Document{}
	entry := resp.Entry

	for _, r := range entry {
		rpt := r.Document
		rpt.FullURL = r.FullURL

		InsertFhirDoc(&rpt, df.Session.DocSessionId)
		docs = append(docs, &rpt)
	}
	fmt.Printf("### CacheDiagRepts Initially found %d diagrpts and returning ###\n", len(docs))
	return
}

/////////////////////////////////Cache FHIR  Diagnostic Report Methods  ////////////////////////////

////////////////////////////////////////////////////////////////////////////////////////////////////////
//                              Process Next Set of DiagRepts returned from  FHIR                     //
////////////////////////////////////////////////////////////////////////////////////////////////////////

func (df *DocumentFilter) FollowDiagNextLinks(links []fhir.Link) {
	// df.Session.UpdateDiagStatus("filling")
	// df.Session.Status.Diagnostic = "filling"
	fmt.Printf("\n\n\n### FollowNextLinks:184 -- Links: %s\n", spew.Sdump(links))
	url := NextDiagPageLink(links)
	i := 1
	for {
		startTime := time.Now()
		if url == "" {
			log.Info("FollowDiagNextLinks: 342 for url is blank, done")
			//df.Session.Status.Diagnostic = "done"
			//df.Session.UpdateDiagStatus( "done")
			break
		}
		//TODO: Get the next page and start its next page while processing current. Do in paraallel
		fmt.Printf("\nFilling Diagnostic Page %d\n", i)
		links, _ = df.ProcessDiagPage(url)
		fmt.Printf("    Page: %d  added in %f seconds\n\n", i, time.Since(startTime).Seconds())
		i = i + 1
		url = NextDiagPageLink(links)
	}
	df.Session.Status.Diagnostic = "done"
	df.Session.UpdateDiagStatus("done")
}

func NextDiagPageLink(links []fhir.Link) string {
	for _, l := range links {
		fmt.Printf("NextDiagPageLink: %s\n", spew.Sdump(l))
		if l.Relation == "next" {
			fmt.Printf("##NextDiagLink:210 %s \n\n", l.URL)
			return l.URL
		}
	}
	fmt.Printf("NextDiagPageLink:214  - No next link\n")
	return ""
}

func (df *DocumentFilter) ProcessDiagPage(url string) ([]fhir.Link, error) {
	fhirDiagResult, err := fhirC.NextFhirDiagRepts(url)
	if err != nil {
		log.Errorf("NextFHIRDiagnosticReports returned err: %s\n", err.Error())
		return nil, err
	}
	docs, err := InsertFhirDocResults(fhirDiagResult, df.Session.DocSessionId)
	if err != nil {
		return nil, err
	}
	if docs == nil {
		return nil, fmt.Errorf("404|no more patients found")
	}
	return fhirDiagResult.Link, nil
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

func GetFhirPerson(per fhir.Person, item string) string {
	if item == "ID" {
		id := strings.Split(per.Reference, "/")
		if len(id) > 1 {
			return id[1]
		} else {
			return ""
		}
	} else {
		return per.Display
	}
}

func GetFhirReference(per fhir.Reference) string {

	id := strings.Split(per.Reference, "/")
	if len(id) > 1 {
		return id[1]
	} else {
		return ""
	}
}

// SplitReference: Accepts a string of ssss/dddd and returns the second part
func SplitReference(ref string) string {
	id := strings.Split(ref, "/")
	if len(id) > 1 {
		return id[1]
	} else {
		return ""
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                                   CA Handlers                                    //
//////////////////////////////////////////////////////////////////////////////////////////////////////

func (df *DocumentFilter) SearchCaDocRef() error {
	// config := ActiveConfig()
	// caURL := config.Env("caServerURL")
	config := ActiveConfig()
	startTime := time.Now()
	caURL := config.Env("caServerURL")
	//page := f.Page
	patID := df.PatientID
	tabID := df.TabID
	if config.Mode() == "dev" {
		log.Warn("DEV Mode forcing Tab to 1000")
		tabID = "1000"
	}
	column := df.Column
	//reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&page=%d&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, page, tabID)

	reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, tabID)
	df.Source = "ca"

	log.Debugf("\n\n\n#####Requesting from CA: %s#####\n\n", reqURL)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", reqURL, nil)
	req.Header.Set("Accept", "application/json")
	authID := df.SessionId
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
	//err = f.ProcessCADocuments(docs)
	log.Infof("GetCADocs process time: %s  total time: %s\n", time.Since(procTime), time.Since(startTime))
	//err = fmt.Errorf("Testing force err")
	//log.Infof("$  CA Request set of documents took %s\n", time.Since(startTime))
	return err
	//return &buf, err
}

/*
// func (df *DocumentFilter) ProcessCaDiagPage(url string) ([]fhir.Link, error) {
// 	//fmt.Printf("\n\n\n###ProcessCaDiagPage:1998\n\n")
// 	//startTime := time.Now()

// 	fhirDiagResult, err := fhirC.NextFhirDiagRepts(url)
// 	if err != nil {
// 		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
// 		return nil, err
// 	}
// 	//fmt.Printf("\n\nNext Page of results: %s\n\n", spew.Sdump(fhirPatients))
// 	err = InsertFhirDiagResults(fhirDiagResult, df.Session.DocSessionId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	diags := []*fhir.DiagnosticReport{}
// 	entry := fhirDiagResult.Entry

// 	for _, r := range entry {
// 		rpt := r.DiagnosticReport

// 		diags = append(diags, &rpt)
// 	}

// 	_, err = df.FhirDiagRptsToCADocuments(diags)
// 	//_, err = FhirDiagReptsToCADocuments(fhirDiagRepts, df) //Documents are cached
// 	if err != nil {
// 		return nil, err
// 	}
// 	//patients := parsePatientResults(fhirPatients, f.SessionId)
// 	//spew.Dump()

// 	if diags == nil {
// 		return nil, fmt.Errorf("404|no more patients found")
// 	}
// 	return fhirDiagResult.Link, nil
// }

func (df *DocumentFilter) FollowNextDiagLinks(links []fhir.Link) {
	df.Session.UpdateDiagStatus("filling")
	df.Session.Status.Diagnostic = "filling"
	fmt.Printf("\n\n\n### FollowNextDiagLinks:2025\n")
	url := NextDiagPageLink(links)
	i := 1
	for {
		startTime := time.Now()
		if url == "" {
			log.Info("CachePages for url is blank, done")
			//df.Session.Status.Diagnostic = "done"
			//df.Session.UpdateDiagStatus( "done")
			break
		}
		//TODO: Get the next page and start its next page while processing current. Do in paraallel
		fmt.Printf("\nFilling Patient Page %d\n", i)
		links, _ = df.ProcessDiagPage(url)
		fmt.Printf("    Page: %d  added in %f seconds\n\n", i, time.Since(startTime).Seconds())
		i = i + 1
		url = NextDiagPageLink(links)
	}
	df.Session.Status.Diagnostic = "done"
	df.Session.UpdateDiagStatus( "done")
}
*/
// func (df *DocumentFilter) ProcessDiagPage(url string) ([]fhir.Link, error) {
// 	fhirDiagResult, err := fhirC.NextFhirDiagRepts(url)
// 	if err != nil {
// 		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
// 		return nil, err
// 	}
// 	docs, err := InsertFhirDocResults(fhirDiagResult, df.Session.DocSessionId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if docs == nil {
// 		return nil, fmt.Errorf("404|no more patients found")
// 	}
// 	return fhirDiagResult.Link, nil
// }

// func (f *DocumentFilter) Search() ([]*fhir.DiagnosticReport, int64, error) {
// 	return nil, 0, nil
// }

// func (df *DocumentFilter) FhirDiagRptsToCADocuments(fds []*fhir.DiagnosticReport) ([]*CADocument, error) {
// 	caDocuments := []*CADocument{}
// 	for _, d := range fds {
// 		doc := df.FhirDiagReptToCADocument(d)
// 		caDocuments = append(caDocuments, doc)
// 	}
// 	//log.Debugf("FhirDiagReptsToCA returned %d documents", len(caDocuments))

// 	return caDocuments, nil
// }

// func (df *DocumentFilter) FhirDiagReptToCADocument(fd *fhir.DiagnosticReport) *CADocument {
// 	var caDoc CADocument
// 	var err error

// 	fmt.Printf("Converting FhirDiag:348 rept: %s\n", spew.Sdump(fd))

// 	caDoc.ID = fd.ID
// 	caDoc.PatientGPI = GetFhirPerson(fd.Subject, "ID")
// 	caDoc.VersionID, err = strconv.ParseUint(fd.Meta.VersionID, 10, 64)
// 	if err != nil {
// 		log.Errorf("Invalid VersionID: [%s] error:%s\n", fd.Meta.VersionID, err.Error())
// 	}
// 	enc := strings.Split(fd.Encounter.Reference, "/")
// 	if len(enc) > 1 {
// 		caDoc.Encounter = enc[1]
// 	} // Encounter is not available
// 	caDoc.Repository = "FHIR"
// 	rpdt := fd.EffectiveDateTime
// 	caDoc.ReptDateTime = &rpdt
// 	caDoc.ImageURL = GetImage(fd, "application/pdf")
// 	caDoc.Pages = 0 // Unavailable
// 	caDoc.Subtitle = fd.Code.Text
// 	caDoc.Text = fd.Text.Div
// 	cfg := ActiveConfig()
// 	caDoc.Source = strings.ToLower(cfg.Source())
// 	if caDoc.Description == "" {
// 		caDoc.Description = fd.Code.Text
// 	}
// 	err = (caDoc).Insert(df.Session.DocSessionId)
// 	if err != nil {
// 		log.Errorf("%s", err.Error())
// 		return nil
// 	}
// 	return &caDoc
// }

//////////////////////////////////////////////////////////////////////////////////////////////////////
//                                 CA DocumentReference Handlers                                    //
//////////////////////////////////////////////////////////////////////////////////////////////////////
//  FhirDocRefsToCA: Accepts an slice of fhir.DocumentReferences and returns a slice of *CADocuments
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

/* // GetCADiagRepts: query the ChartArchive servr for the requested documents and combine them with the fhir documents
// func (f *DocumentFilter) GetCADiagDocs() error {
// 	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"
// 	// config := ActiveConfig()
// 	//caURL := config.Env("caServerURL")
// 	config := ActiveConfig()
// 	startTime := time.Now()
// 	caURL := config.Env("caServerURL")
// 	//page := f.Page
// 	patID := f.PatientID
// 	tabID := f.TabID
// 	if config.Mode() == "dev" {
// 		log.Warn("DEV Mode forcing Tab to 1000")
// 		tabID = "1000"
// 	}
// 	column := f.Column
// 	//reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&page=%d&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, page, tabID)

// 	reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, tabID)
// 	f.Source = "ca"

// 	log.Debugf("\n\n\n#####Requesting from CA: %s#####\n\n", reqURL)

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", reqURL, nil)
// 	req.Header.Set("Accept", "application/json")
// 	authID := f.SessionId
// 	log.Debugf("Auth: %s\n", authID)
// 	req.Header.Set("AUTHORIZATION", authID)
// 	resp, err := client.Do(req)
// 	//resp, err := http.Get(url)
// 	if err != nil {
// 		log.Errorf("CA Query error: %v\n", err)
// 		return err
// 	}

// 	defer resp.Body.Close()
// 	log.Debugf("Status Code: %d - %s\n", resp.StatusCode, resp.Status)
// 	if resp.StatusCode > 299 || resp.StatusCode < 200 {
// 		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
// 		return err
// 	}
// 	log.Debugf("Ready to unmarashal\n")
// 	//var ref fhir.DocumentReferences
// 	var docs []*DocumentSummary
// 	//var docSet CaResults

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	respByte := buf.Bytes()
// 	if err := json.Unmarshal(respByte, &docs); err != nil {
// 		return err
// 	}
// 	log.Infof("GetCADocs request took %s\n", time.Since(startTime))
// 	procTime := time.Now()
// 	//err = BuildDocumentReferences(&ref, session)
// 	fmt.Printf("GetCADocs: %s\n", spew.Sdump(docs))
// 	err = f.ProcessCADocuments(docs)
// 	log.Infof("GetCADocs process time: %s  total time: %s\n", time.Since(procTime), time.Since(startTime))
// 	//err = fmt.Errorf("Testing force err")
// 	//log.Infof("$  CA Request set of documents took %s\n", time.Since(startTime))
// 	return err
// 	//return &buf, err
// }
*/
// func (df *DocumentFilter) SearchCaDocRef() error {
// 	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"
// 	// config := ActiveConfig()
// 	// caURL := config.Env("caServerURL")
// 	config := ActiveConfig()
// 	startTime := time.Now()
// 	caURL := config.Env("caServerURL")
// 	//page := f.Page
// 	patID := df.PatientID
// 	tabID := df.TabID
// 	if config.Mode() == "dev" {
// 		log.Warn("DEV Mode forcing Tab to 1000")
// 		tabID = "1000"
// 	}
// 	column := df.Column
// 	//reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&page=%d&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, page, tabID)

// 	reqURL := fmt.Sprintf("%spatient/%s/docs?column=%s&tab_id=%s&source_values[]=CA&source_values[]=QC&source_values[]=HPF&source_values[]=ca", caURL, patID, column, tabID)
// 	df.Source = "ca"

// 	log.Debugf("\n\n\n#####Requesting from CA: %s#####\n\n", reqURL)

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", reqURL, nil)
// 	req.Header.Set("Accept", "application/json")
// 	authID := df.SessionId
// 	log.Debugf("Auth: %s\n", authID)
// 	req.Header.Set("AUTHORIZATION", authID)
// 	resp, err := client.Do(req)
// 	//resp, err := http.Get(url)
// 	if err != nil {
// 		log.Errorf("CA Query error: %v\n", err)
// 		return err
// 	}

// 	defer resp.Body.Close()
// 	log.Debugf("Status Code: %d - %s\n", resp.StatusCode, resp.Status)
// 	if resp.StatusCode > 299 || resp.StatusCode < 200 {
// 		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
// 		return err
// 	}
// 	log.Debugf("Ready to unmarashal\n")
// 	//var ref fhir.DocumentReferences
// 	var docs []*DocumentSummary
// 	//var docSet CaResults

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	respByte := buf.Bytes()
// 	if err := json.Unmarshal(respByte, &docs); err != nil {
// 		return err
// 	}
// 	log.Infof("GetCADocs request took %s\n", time.Since(startTime))
// 	procTime := time.Now()
// 	//err = BuildDocumentReferences(&ref, session)
// 	//log.Debugf("GetCADocs: %s\n", spew.Sdump(docs))
// 	//err = f.ProcessCADocuments(docs)
// 	log.Infof("GetCADocs process time: %s  total time: %s\n", time.Since(procTime), time.Since(startTime))
// 	//err = fmt.Errorf("Testing force err")
// 	//log.Infof("$  CA Request set of documents took %s\n", time.Since(startTime))
// 	return err
// 	//return &buf, err
// }

/*
// func (f *DocumentFilter) ProcessCADocuments(docs []*DocumentSummary) error {
// 	//caURL := ActiveConfig().Env("caServerURL")NO
// 	var err error

// 	for _, doc := range docs {
// 		doc.SessionID = f.SessionId
// 		doc.EnterpriseID = strconv.FormatUint(doc.CADocumentID, 10) //strconv.Itoa(doc.CADocumentID)
// 		//doc.ImageURL =
// 		doc.SourceType = doc.Source
// 		doc.Source = "ca"
// 		//log.Debugf("@@       Received source %s and sourceType: %s from ca\n", doc.Source, doc.SourceType)
// 		//doc.Source = f.Source  // use the value returned from query
// 		doc.PatientID = f.PatientID
// 		//TODO: doc.PatientGPI = f.PatientGPI
// 		doc.PatientGPI = f.PatientID
// 		doc.Issued = doc.ReptDatetime
// 		doc.ImageURL = fmt.Sprintf("%s/%d", ActiveConfig().Env("caImageURL"), doc.Versionid)
// 		//doc.Versionid = doc.Versionid
// 		log.Debugf("Received Version: %d : %T\n", doc.Versionid, doc.Versionid)
// 		//log.Debugf("ProcessCADocs: %s\n",spew.Dump(doc))
// 		fmt.Printf("### Insert caDiag into cache\n")
// 		doc.Insert()
// 	}
// 	return err
// }

////////////////////////////////////////////////////////////////////////////////////////////
//                               CA Cache Methods                                          /
////////////////////////////////////////////////////////////////////////////////////////////

func (d *CADocument) Insert(sessionId string) error {
	//log.Debug("CADocument insert:506")
	//var insertResult *mongo.InsertOneResult
	var err error

	t := time.Now()
	d.CreatedAt = &t
	d.UpdatedAt = &t
	d.AccessedAt = &t
	d.SessionId = sessionId
	collection, _ := storage.GetCollection("documents")
	//ctx := context.Background()
	_, err = collection.InsertOne(context.TODO(), d)
	if err != nil {
		return fmt.Errorf("Cache ca_document failed: %s", err.Error())
	}
	// if err == nil {
	// 	d.ID = insertResult.InsertedID.(primitive.ObjectID)
	// } else {

	// 	//log.Debugf(" Insert error type: %T : Spew:  %s\n", err, spew.Sdump(err))
	// }

	return err
}

// func (df *DocumentFilter) DocumentCacheStats() (string, int64, int64, error) {
// 	totalInCache, err := df.DocumentsInCache()
// 	if err != nil {
// 		msg := fmt.Sprintf("DocumentCacheStats:534 -- err: %s", err.Error())
// 		return"", 0, 0, errors.New(msg)
// 	}
// 	pageSize := LinesPerPage()
// 	pagesInCache, _ := CalcPages(totalInCache, pageSize)
// 	//pages := inCache/pageSize
// 	log.Debugf("PatientPagesInCache:540 -- pageSize: %d  InCache: %d, pagesInCaches: %d", pageSize, totalInCache, pagesInCache)
// 	cacheStatus := df.Session.GetDocumentStatus()
// 	return cacheStatus, pagesInCache, totalInCache, nil
// }

// func (df *DocumentFilter) DocumentPagesInCache() (int64, error) {
// 	numInCache, err := df.DocumentsInCache()
// 	if err != nil {
// 		return 0, err
// 	}
// 	pageSize := LinesPerPage()
// 	pagesInCache, _ := common.CalcPages(numInCache, pageSize)
// 	//pages := inCache/pageSize
// 	log.Debugf("PatientPagesInCache:5553 -- pageSize: %d  InCache: %d, pagesInCaches: %d", pageSize, numInCache, pagesInCache)
// 	return int64(pagesInCache), nil
// }

func (df *DocumentFilter) CountCachedFhirDocument() (int64, error) {
	return df.DocumentsInCache()
}

func (df *DocumentFilter) DocumentsInCache() (int64, error) {
	//mq := append(f.CacheFilterBase, bson.M{"source": f.Source})
	//filter := bson.M{"$and": mq}

	c, err := storage.GetCollection("documents")
	if err != nil {
		log.Errorf("Settng Collection(%s) failed: %s", "documents", err.Error())
	}
	filter := bson.M{"session_id": df.Session.DocSessionId}
	log.Infof("DocumentsInCache:570 For session matching: [%v]\n", filter)
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("documentInCache:570  Count returned error: %v\n", err.Error())
		return 0, err
	}
	log.Debugf("DocumentsInCache:576 -   Counted %d documents", count)
	return count, nil
}

// func (df *DocumentFilter) CountCachedFhirDocuments() (int64, error) {

// 	filter := bson.M{"subject.reference": "Patient/"+df.PatientGPI} //f.CacheFilterBase //append(f.QueryFilterBase, bson.M{"session_id": f.SessionId})
// 	log.Debugf("CountFhirDocs filter: %v", filter)
// 	log.Debugf("@   Document CountCachedFhirdDocuments matching: %v\n", filter)
// 	collection, err := storage.GetCollection("documents")
// 	if err != nil {
// 		return -1, err
// 	}
// 	count, err := collection.CountDocuments(context.TODO(), filter)
// 	if err != nil {
// 		log.Errorf("CountchedFhirDocuments:642 returned error: %v\n", err)
// 		return 0, err
// 	}
// 	return count, nil
// }

func (df *DocumentFilter) CountCachedCaDocuments() (int64, error) {

	filter := bson.M{"session_id": df.Session.DocSessionId} //f.CacheFilterBase //append(f.QueryFilterBase, bson.M{"session_id": f.SessionId})
	//filter := bson.M{"$and": mq}
	log.Debugf("CountCaDocs filter: %v", filter)
	log.Debugf("@   Document CountCacheCadDocuments matching: %v\n", filter)
	collection, err := storage.GetCollection("documents")
	if err != nil {
		return -1, err
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("CountCAchedDocuments:660 returned error: %v\n", err)
		return 0, err
	}
	return count, nil
}

// func (f *DocumentFilter) CountCachedCaDocumentsFromSource() (int64, error) {

// 	mq := append(f.CacheFilterBase, bson.M{"source": f.Source})
// 	filter := bson.M{"$and": mq}
// 	log.Debugf("@   Document CountCachedDocumentsForSource matching: %v\n", filter)
// 	c, err := storage.GetCollection("ca_documents")
// 	if err != nil {
// 		return -1, err
// 	}
// 	count, err := c.CountDocuments(context.TODO(), filter)
// 	if err != nil {
// 		log.Errorf("  Count returned error: %v\n", err)
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (df *DocumentFilter) GetFhirDocumentPage() ([]*fhir.Document, string, int64, int64, int64, error) {
// 	var linesPerPage int64 = LinesPerPage()
// 	var skip int64 = 0
// 	var caDocs []*CADocument
// 	const DESC = -1
// 	const ASC = 1


// 	if df.Limit > 0 {
// 		//fmt.Printf("GetDocumentPage:553 - setting linesPerPage: %d\n", df.Limit)
// 		linesPerPage = df.Limit
// 	}
// 	if df.Page > 0 {
// 		skip = (df.Page - 1) * linesPerPage
// 		//fmt.Printf("GetDocumentPage:558 -- setting skip: %d\n", skip)
// 	}
// 	if df.Skip > 0 {
// 		skip = df.Skip
// 	}
// 	cacheFilter := bson.M{"session_id": df.Session.DocSessionId}
// 	//log.Debugf("GetDocumentPage:635 -- Debbie  Document filter %s \n", cacheFilter)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(linesPerPage)
// 	findOptions.SetSkip(skip)
// 	sortOrder := ASC // Default Assending
// 	var sortFields bson.D

// 	if strings.ToLower(df.Order) == "desc" {
// 		sortOrder = DESC
// 	}


// 	//sortFields = append(sortFields, bson.E{"rept_datetime", sortOrder})

// 	sort := bson.E{}
// 	if df.Column == "" {
// 		df.Column = "rept_datetime"
// 	}
// 	sort = bson.E{df.Column, sortOrder}
// 	sortFields = append(sortFields, sort)
// 	if len(df.SortBy) > 0 {
// 		for _, s := range df.SortBy {
// 			if s == "visit_num" {
// 				sort = bson.E{"encounter", sortOrder}
// 			} else {
// 				sort = bson.E{s, sortOrder}
// 			}
// 			sortFields = append(sortFields, sort)
// 		}
// 	}
// 	findOptions.SetSort(sortFields)

// 	log.Debugf("GetDocumentPage:667 - cacheFilter: %v", cacheFilter)

// 	findOptions.SetSort(sortFields)
// 	log.Debugf("GetDocumentPage:670 -- sort: %v", sortFields)
// 	collection, _ := storage.GetCollection("ca_documents")
// 	ctx := context.Background()
// 	cursor, err := collection.Find(ctx, cacheFilter, findOptions)
// 	if err != nil {
// 		log.Debugf("GetDovumentPage:675 for %s returned error: %s\n", cacheFilter, err.Error())
// 		return nil, "",0,0,0, err
// 	}
// 	defer func() {
// 		if err := cursor.Close(ctx); err != nil {
// 			log.WithError(err).Warn("Got error while closing Diag cursor")
// 		}
// 	}()
// 	for cursor.Next(ctx) {
// 		var caDoc CADocument
// 		err = cursor.Decode(&caDoc)
// 		if err != nil {
// 			log.WithError(err).Warn("Got error while closing cursor")
// 			return nil, "",0,0,0,err
// 		}
// 		caDocs = append(caDocs, &caDoc)
// 	}
// 	log.Debugf("GetDocumentPage:692 - Finished fetching %d documents from cursor", len(caDocs))
// 	//cursor.Close(context.TODO())
// 	if len(caDocs) == 0 {
// 		err = fmt.Errorf("404|no DiagnosticReports found for %s", cacheFilter)
// 		log.Error(err)
// 	} else {
// 		log.Debugf("QueryDiagCache found %d documents \n", len(caDocs))
// 	}
// 	cacheStatus, pagesInCache, totalInCache, err := df.DocumentCacheStats()
// 	return caDocs, cacheStatus, int64(len(caDocs)), pagesInCache, totalInCache, err

// }

// func (df *DocumentFilter) GetFhirDocumentPage() ([]*fhir.Document, string, int64, int64, int64, error) {
// 	var linesPerPage int64 = LinesPerPage()
// 	var skip int64 = 0
// 	//var caDocs []*CADocument
// 	var fDocs []*fhir.Document
// 	const DESC = -1
// 	const ASC = 1


// 	if df.Limit > 0 {
// 		//fmt.Printf("GetFhirDocumentPage:780 - setting linesPerPage: %d\n", df.Limit)
// 		linesPerPage = df.Limit
// 	}
// 	if df.Page > 0 {
// 		skip = (df.Page - 1) * linesPerPage
// 		//fmt.Printf("GetFhirDocumentPage:785 -- setting skip: %d\n", skip)
// 	}
// 	if df.Skip > 0 {
// 		skip = df.Skip
// 	}
// 	cacheFilter := bson.M{"subject.patient.reference": "Patient/"+df.PatientGPI}
// 	findOptions := options.Find()
// 	findOptions.SetLimit(linesPerPage)
// 	findOptions.SetSkip(skip)
// 	sortOrder := ASC // Default Assending
// 	var sortFields bson.D

// 	if strings.ToLower(df.Order) == "desc" {
// 		sortOrder = DESC
// 	}


// 	//sortFields = append(sortFields, bson.E{"rept_datetime", sortOrder})

// 	sort := bson.E{}
// 	if df.Column == "" {
// 		df.Column = "rept_datetime"
// 	}
// 	sort = bson.E{df.Column, sortOrder}
// 	sortFields = append(sortFields, sort)
// 	if len(df.SortBy) > 0 {
// 		for _, s := range df.SortBy {
// 			if s == "visit_num" {
// 				sort = bson.E{"encounter", sortOrder}
// 			} else {
// 				sort = bson.E{s, sortOrder}
// 			}
// 			sortFields = append(sortFields, sort)
// 		}
// 	}
// 	findOptions.SetSort(sortFields)

// 	log.Debugf("GeFhirtDocumentPage:822 - cacheFilter: %v", cacheFilter)

// 	findOptions.SetSort(sortFields)
// 	log.Debugf("GetFhirDocumentPage:825 -- sort: %v", sortFields)
// 	collection, _ := storage.GetCollection("documents")
// 	ctx := context.Background()
// 	cursor, err := collection.Find(ctx, cacheFilter, findOptions)
// 	if err != nil {
// 		log.Debugf("GetFhirDocumentPage:830 for %s returned error: %s\n", cacheFilter, err.Error())
// 		return nil, "",0,0,0, err
// 	}
// 	defer func() {
// 		if err := cursor.Close(ctx); err != nil {
// 			log.WithError(err).Warn("Got error while closing Diag cursor")
// 		}
// 	}()
// 	for cursor.Next(ctx) {
// 		var fDoc fhir.Document
// 		err = cursor.Decode(&fDoc)
// 		if err != nil {
// 			log.WithError(err).Warn("Got error while closing cursor")
// 			return nil, "",0,0,0,err
// 		}
// 		fDocs = append(fDocs, &fDoc)
// 	}
// 	log.Debugf("GetFhirDocumentPage:847 - Finished fetching %d documents from cursor", len(fDocs))
// 	//cursor.Close(context.TODO())
// 	if len(fDocs) == 0 {
// 		err = fmt.Errorf("GetFhirDocument:850 -- 404|no DiagnosticReports found for %s", cacheFilter)
// 		log.Error(err)
// 	} else {
// 		log.Debugf("QueryDiagCache found %d documents \n", len(fDocs))
// 	}
// 	cacheStatus, pagesInCache, totalInCache, err := df.DocumentCacheStats()
// 	return fDocs, cacheStatus, int64(len(fDocs)), pagesInCache, totalInCache, err

// }

func (f *DocumentFilter) QueryCachedDocuments() ([]*CADocument, error) {
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
	// q, err := f.QueryFhirCacheByEncounter()
	queryFilter := bson.M{"session_id": f.Session.DocSessionId}
	log.Debugf("@@@  720 -- Document QueryCachedDocuments Using filter %s \n", queryFilter)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	// Multi sort fields are separated by [, ].
	// Order is based upon the order of the field names
	// sortFields := strings.Split(f.SortBy,", ")
	// for i, f := range sortFields {

	const DESC = -1
	const ASC = 1
	var sortFields bson.D
	var caDocs []*CADocument
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
	log.Debugf("@    757 -- sort: %v", sortFields)
	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, queryFilter, findOptions)
	if err != nil {
		log.Debugf("QueryDiagCache for %s returned error: %v\n", queryFilter, err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.WithError(err).Warn("Got error while closing Diag cursor")
		}
	}()
	for cursor.Next(ctx) {
		var caDoc CADocument
		err = cursor.Decode(&caDoc)
		if err != nil {
			log.WithError(err).Warn("Got error while closing cursor")
			return nil, err
		}
		caDocs = append(caDocs, &caDoc)
	}
	log.Debugf("   Finished fetching documents from cursor")
	//cursor.Close(context.TODO())
	if len(caDocs) == 0 {
		err = fmt.Errorf("404|no DiagnosticReports found for %s", queryFilter)
		log.Error(err)
	} else {
		log.Debugf("QueryDiagCache found %d documents \n", len(caDocs))
	}
	return caDocs, err
}

// func (f *DocumentFilter) CheckForCaDocUpdates() error {

// 	log.Debugf("In CheckForCaDocUpdates\n")
// 	latest, err := f.FindLatestCaCachedDocs()
// 	if err != nil {
// 		log.Errorf("Latest returned error: %v\n", err)
// 		return err // there are none
// 	}
// 	lastDate := fmt.Sprintf("$gt%s", latest.ReptDatetime.Format("2006-01-02 15:04:05"))
// 	query := fmt.Sprintf("patient=%s&created=%s&created=$le2500-12-31", f.PatientID, lastDate)
// 	log.Debugf("\n\n@    CheckForCaDocsUpdates query: %s\n", query)

// 	err = f.QueryFhirDiagRepts(query) // if any found they are in the cache
// 	return err
// }

func (f *DocumentFilter) FindFhirDiagnosticReports() ([]*fhir.Document, error) {
	c := config.Fhir()
	log.Debugf("FindFhirDiagnosticReports:737 -  is searching DiagnosticReports fo Patient: %s\n", f.PatientGPI)

	//start := time.Now()
	f.MakeFhirQuery()

	//q := fmt.Sprintf("patient=%s", f.fhirQuery)
	log.Debugf("c.FindDiagnosticReports is called with %s", f.fhirQuery)
	resp, err := c.FindDiagnosticReports(f.fhirQuery)

	if err != nil {
		return nil, err
	}
	for _, l := range resp.Link {
		if l.Relation == "next" {
			nextLink := l.URL
			fmt.Printf("\n\n\n")
			log.Debugf("FindFhirDiagnosticReports:752 - NextLink: %s", nextLink)
			//go f.ProcessRemainingDiagRepts(nextLink)
			break
		} // if no next link there are no more DiagRepts
	}

	log.Debugf("fhir returned no error")
	docs := []*fhir.Document{}
	entry := resp.Entry

	for _, r := range entry {
		rpt := r.Document

		docs = append(docs, &rpt)
	}
	// ds := fillDiagnosticReports(diag, f.SessionId)
	// log.Debugf("%d documents returned from Diagnostic:505\n", len(ds))

	// ref, err := f.FhirDocumentReferences()
	// if err != nil {
	// 	log.Errorf("   FhirDocumentReference error: %s\n", err)
	// 	return nil, err
	// }
	// ds = fromDocumentReferences(ref, f)
	log.Debugf("%d DiagnosticReports returned\n", len(docs))
	//fmt.Printf("diag-1: %s\n", spew.Sdump(diags[0]))
	//fmt.Printf("diag-2: %s\n", spew.Sdump(diags[1]))
	return docs, err
}



// func (f *DocumentFilter) Search() ([]*fhir.DiagnosticReport, int64, error) {
// 	f.MakeFhirQuery()
// 	if f.fhirQuery == "" {
// 		return nil, 0, fmt.Errorf("No fhirQery was created")
// 	}

// 	log.Debugf("Document#Find fhirQuery: %v\n", f.fhirQuery)

// 	if f.ID != "" {
// 		f.UseCache = "true"
// 		diagDocs, cnt, err := f.FindDiagRepts()
// 		if err != nil {
// 			return nil, cnt, err
// 		}
// 		// err = ds.fromDocumentReference(ref, f.SessionId)
// 		// if err != nil {
// 		// 	return nil, int64(0), err
// 		// }
// 		// var summaries []*DocumentSummary
// 		//summaries = append(summaries, ds)
// 		return diagDocs, cnt, err
// 	}

// 	if f.PatientGPI == "" {
// 		//if ok && strings.ToLower(caOnly) == "true" {
// 		log.Debugf("@  Handling CA Documents ONLY\n")
// 		page := f.Page
// 		f.Source = "ca"

// 		log.Debugf("Search CA Only: %s\n", spew.Sdump(f))
// 		err := f.GetCADocs()
// 		if err != nil {
// 			log.Errorf("GetCaDocs returned error %v\n", err)
// 		}
// 		f.Page = page
// 		f.PatientGPI = f.PatientID
// 		//f.PatientID = f.PatientGPI
// 		//f.makeQueryFilter()
// 		//log.Debugf("@  New Query: %v\n", f.QueryFilter)

// 		// if ActiveConfig().Env("ca_only") != "true" {
// 		// 	_, _, err = f.FindDocuments()
// 		// }

// 		docs, _ := f.QueryFhirCache()
// 		total, err := f.CountCachedDocuments()

// 		return docs, total, err
// 	} else {
// 		log.Infof("looking for both fhir and CA patients gpi: [%s]", f.PatientGPI)
// 	}

// 	// if f.TabID != "" {
// 	// 	log.Debugf("\n\n@  Handling CA Documents\n")
// 	// 	page := f.Page
// 	// 	f.Source = "ca"
// 	// 	//f.Source = ""
// 	// 	err := f.GetCADocs()
// 	// 	f.Page = page
// 	// 	f.PatientID = f.PatientGPI
// 	// 	f.makeQueryFilter()
// 	// 	log.Debugf("@  New Query: %v\n", f.QueryFilter)

// 	// 	_, _, err = f.FindDiagReports()

// 	// 	docs, _ := f.QueryDiagCache()
// 	// 	total, err := f.CountCachedDocuments()

// 	// 	return docs, total, err
// 	// }

// 	// No mrn search available.  ONLY fhir patient id.
// 	if f.MRN != "" {
// 		// pf := PatientFilter{MRN: f.MRN, Session: f.Session}

// 		// pats, err := pf.Search()
// 		// if err != nil {
// 		// 	log.Errorf("DocumentFilter looking for MRN: %s failed with err: %v\n", f.MRN, err)
// 		// 	return nil, 0, fmt.Errorf("404|Patient was not found for MRN: %s", pf.MRN)
// 		// }
// 		//f.PatientID = pats[0].EnterpriseID
// 		//f.PatientGPI = f.PatientID
// 		//f.MRN = ""

// 		//log.Debugf("@   documentFilter before Rebuild (Patient should be blank) : %v\n", f.queryString)

// 		// err = f.makeQueryFilter()
// 		// if err != nil {
// 		// 	log.Errorf("!   makeQueryFilter err: %v\n", err)
// 		// }

// 		//log.Debugf("\n     @@documentFilter update: %v\n", f.queryString)
// 	}

// 	diagRepts, total, err := f.FindDiagRepts()
// 	log.Infof("Search found %d documents total %d in %s\n", len(diagRepts))
// 	return diagRepts, total, err
// }

// func (f *DocumentFilter) FindDiagRepts() ([]*fhir.DiagnosticReport, int64, error) {
// 	var docs []*fhir.DiagnosticReport
// 	var err error
// 	//log.Debugf("\n\nFindDiagRepts: Use Cache: %s\n", f.UseCache)
// 	startTime := time.Now()
// 	if strings.ToLower(f.ClearCache) == "true" {
// 		// TODO: Implement f.ClearCache and be sure to set it to false once done
// 	}

// 	f.Source = ActiveConfig().Source()
// 	if f.Page == 0 {
// 		log.Debugf("\n@@    FindDiagRports Check for updates in FHIR\n")
// 		//f.UseCache = "false"

// 		err := f.CheckForFhirDiagUpdates()
// 		log.Debugf("@  CheckForFhirDiagRepts returned err: %v\n", err)
// 		f.UseCache = "true"
// 		f.Page = 1
// 	}
// 	//if f.UseCache == "true" {
// 	log.Debugf("@@@     FindDocuments-263 Use Cache with filter: %v\n", f.QueryFilter)
// 	count, _ := f.CountCachedDocumentsFromSource()
// 	if count > 0 { // We have cached documents, Return them
// 		log.Debugf("@@@     FindDocuments 266 found %d Cache documents,\n", count)
// 		docs, _ = f.QueryDiagCache() // Filter how they user requested from cache
// 		log.Infof("      %d found in FindDocuments-232 in %s\n\n", len(docs), time.Since(startTime))
// 		return docs, count, nil
// 	}

// 	//} else {
// 	log.Debugf("@  FindDocuments-273 is not using the Cache. Requesting from FHIR")
// 	//}

// 	fhirTime := time.Now()
// 	docs, err = f.GetFhirDiagRepts()
// 	log.Infof("@    Get Fhir Documents took %s  returning %d documents\n", time.Since(fhirTime), len(docs))
// 	// need to add query for fhir document references
// 	if err != nil { // NOne were found in FHIR
// 		return nil, 0, err
// 	}

// 	log.Debugf("\n\n@@@      Count Total documents=284 in Cache with filter: %v\n", f.QueryFilter)
// 	total, err := f.CountCachedDocuments()
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	//total, _ := CountCachedDocuments(f.QueryFilter)
// 	log.Debugf("@@@      %d found in FindDocuments-290 Count returned %d \n\n", len(docs), total)

// 	docs, _ = f.QueryFhirCache() // Filter how they user requested from cache
// 	log.Infof("@@@      %d found in FindDocuments-293  QueryFhirCache in %s\n\n", len(docs), time.Since(startTime))
// 	return docs, total, nil
// }

// func (f *DocumentFilter) CheckForFhirDiagUpdates() error {

// 	log.Debugf("In CheckForFhirDiagUpdates\n")
// 	latest, err := f.FindLatestCachedFhir()
// 	if err != nil {
// 		log.Errorf("Latest returned error: %v\n", err)
// 		return err // there are none
// 	}
// 	lastDate := fmt.Sprintf("$gt%s", latest.ReptDatetime.Format("2006-01-02 15:04:05"))
// 	query := fmt.Sprintf("patient=%s&created=%s&created=$le2500-12-31", f.PatientID, lastDate)
// 	log.Debugf("\n\n@    CheckForFhirDiagUpdates query: %s\n", query)

// 	err = f.QueryFhirDiagRepts(query) // if any found they are in the cache
// 	return err
// }

// func (f *DocumentFilter) FindDocument() ([]*fhir.DiagnosticReport, error) {
// 	var docs []*fhir.DiagnosticReport
// 	var err error
// 	startTime := time.Now()

// 	if f.UseCache == "true" {
// 		log.Debugf("\n\n@@@ FindDocument Use Cache with filter: %v\n", f.QueryFilter)
// 		count, err := CountCachedDocuments(f.QueryFilter)
// 		if err != nil {
// 			log.Errorf("FindDocument: CountCache returned err: %v\n", err)
// 			return nil, err
// 		}
// 		if count > 0 { // Return what we have in the cache
// 			docs, _ = f.QueryFhirCache() // Filter how they user requested from cache
// 			log.Infof("      %d found in FindDocument-320 in %s\n\n", len(docs), time.Since(startTime))
// 			return docs, nil
// 		}
// 	}

// 	doc, err := f.FhirDocument()
// 	if err != nil {
// 		log.Errorf("    FindDocument: FhirDocument returned: %v\n", err)
// 		return nil, err
// 	}

// 	docs = append(docs, doc)
// 	//docs, _ = f.QueryFhirCache() // Filter how they user requested from cache
// 	log.Infof("      %d found in FindDocument-334 in %s\n\n", len(docs), time.Since(startTime))
// 	return docs, nil

// }

// FhirDiagRepts returns the document by ID
// It searches Cached first then DocumentReferences and DiagnosticReport returning the first it finds
// Need to figure out how to handle both diagnostic and references

// func (f *DocumentFilter) FhirDiagRepts() ([]*fhir.DiagnosticReport, error) {
// 	//c := config.Fhir()
// 	var diags []*fhir.DiagnosticReport
// 	log.Debugf("\n@    FhirDocument is searching FhirDidagnostics for Patient: %s \n", f.PatientID)

// 	startTime := time.Now()

// 	fds, err := f.FhirDocumentReferences()
// 	if err != nil {
// 		log.Errorf("   FhirDiagRepts: FhirDocumentReferences error: %v\n", err)
// 		return nil, err
// 	}
// 	//log.Infof("$   FHIR Request took %s\n", time.Since(startTime))
// 	//startTime = time.Now()
// 	//ds = fromDocumentReferences(fds, f)
// 	if err != nil {
// 		log.Errorf("   fromDocumentReferences returned err: %v\n", err)
// 		return nil, err
// 	}
// 	log.Debugf("$   Create DocumentReferences took %s to process %d documents\n", time.Since(startTime), len(ds))

// 	if ds == nil {
// 		errMsg := fmt.Sprintf("FhirDiagRepts Returned nothing for Session: %s and patient: %s", f.SessionId, f.PatientID)
// 		log.Errorf("%s\n", errMsg)
// 		err = fmt.Errorf(errMsg)
// 		return nil, err
// 	}
// 	return ds, err
// }

// QueryFhirDiagRepts Querys for FHIR documents using the provided query
// It searches Cached first then DocumentReferences and DiagnosticReport returning the first it finds
// Need to figure out how to handle both diagnostic and references

// func (f *DocumentFilter) QueryFhirDiagRepts(query string) error {
// 	c := config.Fhir()
// 	var ds []*DocumentSummary
// 	log.Debugf("\n@    QueryFhirDocument is searching FhirDiagRepts using query: %s \n", query)
// 	startTime := time.Now()
// 	fds, err := c.GetDocumentReferences(query)
// 	if err != nil {
// 		return err
// 	}
// 	log.Infof("@   FHIR Request took %s\n", time.Since(startTime))
// 	startTime = time.Now()
// 	ds = fromDocumentReferences(fds, f)
// 	log.Infof("@   Create DocumentReferences took %s to process %d documents\n", time.Since(startTime), len(ds))
// 	return nil
// }

// FhirDocument returns the document by ID
// It searches Cached first then DocumentReferences and DiagnosticReport returning the first it finds
// Need to figure out how to handle both diagnostic and references
// func (f *DocumentFilter) FhirDocument() (*DocumentSummary, error) {
// 	//c := config.Fhir()
// 	var ds = new(DocumentSummary)

// 	//log.Debugf("\nFhirDocument-297 is searching All Fhir documents for Patient: %s - DocumentID: %s\n", f.PatientID, f.EnterpriseID)

// 	ref, err := f.FhirDocumentReference()
// 	if err != nil {
// 		log.Errorf("FhirDocument error: %v\n", err)
// 		return nil, err
// 	}
// 	err = ds.fromDocumentReference(ref, f.SessionId)
// 	if err != nil {
// 		log.Errorf("!   fromDocumentReference returned err: %v\n", err)
// 		return nil, err
// 	}

// 	if ds == nil {
// 		errMsg := fmt.Sprintf("FhirDocument Returned nothing for Session: %s and EnterpriseID: %s", f.SessionId, f.EnterpriseID)
// 		log.Errorf("%s\n", errMsg)
// 		err = fmt.Errorf(errMsg)
// 		return nil, err
// 	}

// 	return ds, err
// }

//: Figure out how to passes what goroutine needs

// func fromDocumentReferences(docR *fhir.DocumentReferences, f *DocumentFilter) []*DocumentSummary {
// 	session := f.Session
// 	cacheName := session.SessionID

// 	var docs []*DocumentSummary
// 	var nextLink string
// 	startTime := time.Now()
// 	for _, l := range docR.Link {
// 		if l.Relation == "next" {
// 			nextLink = l.URL
// 			go ProcessRemainingDocuments(nextLink, session)
// 			break
// 		}
// 	}
// 	for _, entry := range docR.Entry {
// 		var doc DocumentSummary
// 		for _, item := range entry.DocumentReference.Content {
// 			//log.Debugf("Attachment: %v\n", item.Attachment.URL)
// 			if item.Attachment.ContentType == "application/pdf" {
// 				doc.FhirImageURL = item.Attachment.URL
// 				doc.makeImageURL()
// 			}
// 		}
// 		// 	doc.EffectiveDate = entry.Resource.ResourcePartial.EffectiveDateTime

// 		doc.Source = ActiveConfig().Source()
// 		doc.SourceType = "Reference"
// 		doc.Text = entry.DocumentReference.Text
// 		//doc.Code = entry.Resource.Code.Text
// 		// 	doc.Category = entry.Resource.Category.Text
// 		doc.FullLink = entry.FullURL
// 		doc.EnterpriseID = entry.DocumentReference.ID
// 		doc.Subject = entry.DocumentReference.Subject
// 		//TODO: check both patient and subject for the patient information
// 		doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
// 		doc.PatientGPI = doc.PatientID
// 		doc.Performer = entry.DocumentReference.Authenticator
// 		doc.ReptDatetime = entry.DocumentReference.Created
// 		enc := strings.Split(entry.DocumentReference.Context.EncounterNum.Reference, "/")
// 		//log.Debugf("\nEncounter: %v\n", enc)
// 		if len(enc) > 1 {
// 			doc.Encounter = enc[1]
// 		}
// 		doc.Description = entry.DocumentReference.Description
// 		doc.Category = entry.DocumentReference.Type.Text
// 		// if item.Attachment.ContentType == "application/pdf" {
// 		// 	doc.FhirImageURL = item.Attachment.URL
// 		// 	doc.makeImageURL()
// 		// }
// 		doc.SessionID = cacheName
// 		err := (&doc).Insert()
// 		if err != nil {
// 			log.Errorf("insert document into Summary Cache: %s\n", err.Error())
// 		}
// 		docs = append(docs, &doc)
// 		//log.Debugf("\nFillEntry returing\n")
// 	}
// 	log.Infof("\n#    Cached %d documents in %s\n", len(docs), time.Since(startTime))
// 	return docs

// }

// func (f *DocumentFilter) FindFhirDiagnosticReports() ([]*fhir.DiagnosticReport, error) {
// 	c := config.Fhir()
// 	log.Debugf("\nFHirDiagnosticDocuments is searching DiagnosticDocuments fo Patientr: %s\n", f.PatientID)

// 	//start := time.Now()
// 	f.MakeFhirQuery()

// 	//q := fmt.Sprintf("patient=%s", f.fhirQuery)
// 	resp, err := c.FindDiagnosticReports(f.fhirQuery)
// 	if err != nil {
// 		return nil, err
// 	}
// 	diags := []*fhir.DiagnosticReport{}
// 	entry := resp.Entry
// 	for _, rpt := range entry {
// 		diags = append(diags, &rpt.DiagnosticReport)
// 	}
// 	// ds := fillDiagnosticReports(diag, f.SessionId)
// 	// log.Debugf("%d documents returned from Diagnostic:505\n", len(ds))

// 	// ref, err := f.FhirDocumentReferences()
// 	// if err != nil {
// 	// 	log.Errorf("   FhirDocumentReference error: %s\n", err)
// 	// 	return nil, err
// 	// }
// 	// ds = fromDocumentReferences(ref, f)
// 	// log.Debugf("%d documents returned from References\n", len(ds))
// 	return diags, err
// }

func (f *DocumentFilter) GetFhirDiagnosticReport() (*fhir.DocumentResults, error) {
	c := config.Fhir()
	fmt.Printf("\n\n####GetFhirDiagnosticReport filter: %s\n", spew.Sdump(f))
	//start := time.Now()
	//q := fmt.Sprintf("patient=%s", f.PatientID)
	f.MakeFhirQuery()
	qry := f.fhirQuery
	//log.Debugf("GetFhirDiagnosticDocument is searching : %s\n", spew.Sdump(f))
	results, err := c.FindDiagnosticReports(qry)
	if err != nil {
		log.Errorf("Fhir:1400 returned error")
		msg := fmt.Sprintf("FindDiagnosticReports error: %s", err.Error())
		log.Errorf(msg)
		return nil, fmt.Errorf(msg)
	}
	//ds := fillDiagnosticReports(diag, f.SessionId)
	//log.Debugf("%d documents returned from  Diagnostic:530\n", len(ds))

	//ref, err := f.FhirDocumentReferences()
	if err != nil {
		return nil, err
	}
	//ds = fromDocumentReferences(ref, f)
	//log.Debugf("%d documents returned from References\n", len(ds))
	//query the cache for the requested documents.
	drs, err := f.QueryFhirCache()
	log.Debugf("%d documents returned from query\n", len(drs))

	return results, err
}

func (df *DocumentFilter) FhirDocumentReferences() (*fhir.DocumentResults, error) {
	c := config.Fhir()
	log.Debugf("\n@  FhirDocumentReferences is searching using fhirQuery: %s  queryFilter: %v\n", df.fhirQuery, df.CacheFilter)
	if df.queryMap["patient"] == "" {
		err := fmt.Errorf("Patient id required for Document search.  Patient and other criteria is valid")
		return nil, err
	}
	//start := time.Now()

	//q := fmt.Sprintf("patient=%s", f.PatientGPI)
	q := df.fhirQuery
	log.Debugf("@  query: %v\n", q)
	//q := f.queryString

	results, err := c.FindDocumentReferences(q)

	//log.Infof("FhirDiagnosticDocuments took %s\n", time.Since(start))
	if err != nil {
		log.Errorf("       FhirDocumentReferences returned err: %v\n", err)
		return nil, err
	}
	return results, nil
}

func (df *DocumentFilter) FhirDocumentReference() (*fhir.DocumentReference, error) {
	c := config.Fhir()
	log.Debugf("\nFhirDocumentReference is searching for ID: %s using fhirQuery: %s  queryFilter: %v\n",
			df.EnterpriseID, df.fhirQuery, df.CacheFilter)

	dRef, err := c.GetDocumentReference(df.ID)

	//log.Infof("FhirDiagnosticDocuments took %s\n", time.Since(start))
	if err != nil {
		log.Errorf("      FhirDiagRepts returned err: %v\n", err)
		return nil, err
	}
	return dRef, nil
}

// func (d *DocumentSummary) Insert() error {
// 	var insertResult *mongo.InsertOneResult
// 	var err error
// 	d.setDates()
// 	collection, _ := storage.GetCollection("doc_summary")
// 	//ctx := context.Background()
// 	insertResult, err = collection.InsertOne(context.TODO(), d)

// 	if err == nil {
// 		d.ID = insertResult.InsertedID.(primitive.ObjectID)
// 	} else {

// 		//log.Debugf(" Insert error type: %T : Spew:  %s\n", err, spew.Sdump(err))
// 	}

// 	return err
// }

// func (d *DocumentSummary) insertWithSession(session *mongo.Session) error {
// 	var insertResult *mongo.InsertOneResult
// 	var err error
// 	//d.SetDates()
// 	collection, _ := storage.GetCollection("documents")
// 	ctx := context.Background()
// 	if err = mongo.WithSession(ctx, *session, func(sc mongo.SessionContext) error {
// 		if insertResult, err = collection.InsertOne(ctx, d); err != nil {
// 			// Ignore errors for now
// 			//  Need to only ignore if exists report others
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		//log.Errorf("Insert with session failed: %v\n", err)
// 		return nil
// 	}

// 	d.ID = insertResult.InsertedID.(primitive.ObjectID)

// 	return nil
// }

// func (p *DocumentSummary) setDates() {
// 	t := time.Now()
// 	p.CreatedAt = &t
// 	p.UpdatedAt = &t
// 	p.AccessedAt = &t
// }

// func (doc *DocumentSummary) fromDocumentReference(docR *fhir.DocumentReference, sessionID string) error {
// 	doc.Source = ActiveConfig().Source()
// 	doc.SourceType = "Reference"
// 	doc.Text = docR.Text
// 	//doc.FullLink = docR.FullURL
// 	doc.EnterpriseID = docR.ID
// 	doc.Subject = docR.Subject
// 	doc.Meta = docR.Meta
// 	//TODO: check both patient and subject for the patient information
// 	doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
// 	doc.Performer = docR.Authenticator
// 	doc.ReptDatetime = docR.Created
// 	enc := strings.Split(docR.Context.EncounterNum.Reference, "/")

// 	if len(enc) > 1 {
// 		doc.Encounter = enc[1]
// 	}
// 	doc.Description = docR.Description
// 	doc.Category = docR.Type.Text
// 	for _, c := range docR.Content {
// 		if c.Attachment.ContentType == "application/pdf" {
// 			doc.FhirImageURL = c.Attachment.URL
// 			doc.makeImageURL()
// 		}
// 	}
// 	doc.SessionID = sessionID
// 	doc.Insert()
// 	return nil
// }

func (f *DocumentSummary) makeImageURL() {
	config := ActiveConfig()
	fhirURL := config.ImageURL()

	f.ImageURL = fmt.Sprintf("%s%s", fhirURL, f.EnterpriseID)
	return

}

// func fillDocumentSummary(diag *fhir.DiagnosticReport) *DocumentSummary {
// 	// *DocumentSummary {
// 	//fmt.Println("\n\n================diag===")

// 	var doc DocumentSummary
// 	return &doc
// }

// func fillDiagnosticReports(diags *fhir.DiagnosticReportResponse, cacheName string) []*DocumentSummary {
// 	var docs []*DocumentSummary

// 	for _, entry := range diags.Entry {
// 		if entry.DiagnosticReport.FullURL == "" {
// 			entry.DiagnosticReport.FullURL = entry.FullURL
// 		}
// 		var doc DocumentSummary
// 		rpt := entry.DiagnosticReport
// 		doc.ReptDatetime = entry.DiagnosticReport.EffectiveDateTime
// 		doc.Meta = rpt.Meta
// 		//	fmt.Printf("\n\n\n\n\n\n###Meta : %s\n", spew.Sdump(rpt.Meta))
// 		doc.Source = ActiveConfig().Source()
// 		doc.SourceType = "Diagnostic"
// 		doc.Text = rpt.Text
// 		doc.Code = rpt.Code.Text
// 		doc.Category = rpt.Category.Text
// 		doc.FullLink = entry.FullURL
// 		doc.EnterpriseID = rpt.ID
// 		doc.Subject = rpt.Subject
// 		//TODO: check both patient and subject for the patient information
// 		doc.PatientID = rpt.ID
// 		doc.Performer = rpt.Performer
// 		enc := strings.Split(rpt.Encounter.Reference, "/")
// 		if len(enc) > 1 {
// 			doc.Encounter = enc[1]
// 		}
// 		var pf PresentedForm
// 		for _, v := range rpt.PresentedForm {
// 			pf.ContentType = v.ContentType
// 			pf.URL = v.URL
// 			doc.Images = append(doc.Images, pf)
// 		}
// 		for _, attachment := range rpt.PresentedForm {
// 			if attachment.ContentType == "application/pdf" {
// 				doc.FhirImageURL = attachment.URL
// 				doc.makeImageURL()
// 			}
// 		}
// 		//doc.Images = entry.Resource.PresentedForm
// 		doc.SessionID = cacheName
// 		(&doc).Insert()
// 		docs = append(docs, &doc)
// 	}
// 	return docs
// }

// func (ds *DocumentSummary) fillDiagnosticReport(fdoc *fhir.DiagnosticReport, cacheName string) error {
// 	var doc DocumentSummary
// 	// for _, attachment := range entry.Resource.PresentedForm {
// 	// 	if attachment.ContentType == "application/pdf" {
// 	// 		doc.FhirImageURL = attachment.URL
// 	// 		doc.makeImageURL()
// 	// 	}
// 	// }

// 	doc.ReptDatetime = fdoc.EffectiveDateTime

// 	doc.Source = ActiveConfig().Source()
// 	doc.SourceType = "Diagnostic"
// 	doc.Text = fdoc.Text
// 	doc.Code = fdoc.Code.Text
// 	doc.Category = fdoc.Category.Text
// 	doc.FullLink = fdoc.FullURL
// 	//doc.EnterpriseID = fdoc.ID
// 	//doc.Subject = fdoc.Subject
// 	//TODO: check both patient and subject for the patient information
// 	doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
// 	//doc.Performer = fdoc.Performer
// 	enc := strings.Split(fdoc.Encounter.Reference, "/")
// 	if len(enc) > 1 {
// 		doc.Encounter = enc[1]
// 	}
// 	var pf PresentedForm
// 	for _, v := range fdoc.PresentedForm {
// 		pf.ContentType = v.ContentType
// 		pf.URL = v.URL
// 		doc.Images = append(doc.Images, pf)
// 	}
// 	for _, attachment := range fdoc.PresentedForm {
// 		if attachment.ContentType == "application/pdf" {
// 			doc.FhirImageURL = attachment.URL
// 			doc.makeImageURL()
// 		}
// 	}
// 	//doc.Images = entry.Resource.PresentedForm
// 	doc.SessionID = cacheName
// 	(&doc).Insert()
// 	ds = &doc
// 	return nil
// }

// func mapToQueryString(q map[string]string) string {
// 	var query string
// 	for k, v := range q {
// 		//log.Debugf("%s=%s\n", k, v)
// 		s := fmt.Sprintf("%s=%s", k, v)
// 		if query == "" {
// 			query = query + s

// 		} else {
// 			query = query + fmt.Sprintf("&%s", s)
// 		}

// 	}
// 	//log.Debugf("    Result Query: %s\n", query)
// 	return query
// }

//func GetDocumentImage(docID string) (*DocumentSummary, error) { //(*[]byte, error)
// func (d *DocumentSummary) GetDocumentImage() error { //(*[]byte, error) {

// 	log.Debugf("GetDocumentImage DocID: %s\n", d.EnterpriseID)
// 	// ss := strings.Split(url, "/")
// 	// docID := ss[len(ss)-1]
// 	// log.Debugf("ss: %v\n", ss)
// 	filter := bson.M{"enterprise_id": d.EnterpriseID}
// 	err := d.CachedDocument(filter)
// 	if err != nil {

// 		log.Errorf("GetDocumentImage: [%s] GetCache error %v\n", d.EnterpriseID, err)
// 		return err
// 	}

// 	if d.Image != "" {
// 		log.Debugf("   Using cached image\n")
// 		return nil // We have an image do not get from cerner
// 	}
// 	url := d.FhirImageURL
// 	log.Debugf("Full Url: %s\n", url)
// 	docImage, err := readImageFromURL(url)
// 	if err != nil {
// 		log.Errorf("GetDocumentImage Error: %v\n", err)
// 		return err
// 	}
// 	//doc.Image = image
// 	d.Image = docImage.Content
// 	d.UpdateImage()

// 	return nil
// }

// func readImageFromURL(url string) (*DocumentImage, error) {
// 	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"

// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Accept", "application/json+fhir")
// 	resp, err := client.Do(req)
// 	//resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()
// 	log.Debugf("Status Code: %d - %s\n", resp.StatusCode, resp.Status)
// 	if resp.StatusCode > 299 || resp.StatusCode < 200 {
// 		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
// 		return nil, err
// 	}
// 	var image DocumentImage

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	respByte := buf.Bytes()

// 	if err := json.Unmarshal(respByte, &image); err != nil {
// 		return nil, err
// 	}
// 	return &image, nil
// 	//return &buf, err
// }

// func (rec *DocumentSummary) UpdateImage() {
// 	filter := bson.M{"_id": rec.ID}
// 	//log.Debugf("@@ Updater: %v\n", filter)

// 	update := bson.M{"$set": bson.M{"image": rec.Image}}
// 	collection, _ := storage.GetCollection("documents")
// 	res, err := collection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		log.Errorf(" Update error ignored: %s\n\n", err)
// 	}
// 	log.Debugf("Matched: %d  -- modified: %d\n", res.MatchedCount, res.ModifiedCount)

// 	return
// }

//Get Chached entries first. If none, get them from FHIR
// func CountCachedDocuments(filter bson.M) (int, error) {
// 	log.Debugf("Document CountCachedDocuments:1192 matching: %v\n", filter)
// 	collection, err := storage.GetCollection("ca_documents")
// 	if err != nil {
// 		return -1, err
// 	}
// 	count, err := collection.CountDocuments(context.TODO(), filter)
// 	if err != nil {
// 		return 0, err
// 	}
// 	log.Debugf("CountCachedDocs:1198 returned a count of %d documents", count)
// 	return int(count), nil
// }

//Get Chached entries first. If none, get them from FHIR
// func (f *DocumentFilter) CountCachedDocuments() (int64, error) {
// 	filter := f.makeQueryFilter
// 	log.Debugf("Document CountCachedDocuments matching: %v\n", filter)
// 	c, err := storage.GetCollection("documents")
// 	count, err := c.CountDocuments(context.TODO(), filter)
// 	if err != nil {
// 		log.Errorf("  Count returned error: %v\n", err)
// 		return 0, err
// 	}
// 	return count, nil
// }

//Get Chached entries first. If none, get them from FHIR
func (f *DocumentFilter) CountCachedDocumentsFromSource() (int64, error) {

	mq := append(f.CacheFilterBase, bson.M{"source": f.Source})
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
// func GetCachedDocuments(filter bson.M) ([]*DocumentSummary, error) {
// 	log.Debugf("Document GetCachedDocuments using %v\n", filter)
// 	var documents []*DocumentSummary
// 	collection, err := storage.GetCollection("documents")
// 	if err != nil {
// 		//log.Debugf(" Error getting Collection: %s\n", err)
// 		return nil, err
// 	}
// 	//var encounter = new(Encounter)
// 	//log.Debugf("\nGetCacheDocuments using Filter: %v\n", filter)
// 	cursor, err := collection.Find(context.TODO(), filter)
// 	if err != nil {
// 		//log.Debugf("GetCachedDocuments search for %s returned error: %v\n", filter, err)
// 		//cursor.Close(context.TODO())
// 		return nil, err
// 	}
// 	for cursor.Next(context.TODO()) {
// 		var document DocumentSummary
// 		err = cursor.Decode(&document)
// 		if err != nil {
// 			cursor.Close(context.TODO())
// 			return nil, err
// 		}
// 		//.Dump(documents)
// 		documents = append(documents, &document)
// 	}
// 	if documents == nil {
// 		err = fmt.Errorf("404|no documents found for %s", filter)
// 	}
// 	return documents, err
// }

//Get Chached entries first. If none, get them from FHIR
// func GetCachedDocumentSummary(filter bson.M) (*DocumentSummary, error) {
// 	log.Debugf("DocumentGetCachedDocument using: %v\n", filter)
// 	var document DocumentSummary
// 	collection, err := storage.GetCollection("documents")
// 	if err != nil {

// 		return nil, err
// 	}
// 	log.Debugf("\nGetCachedDocument Checking for DocumentSummary with query: %v\n", filter)
// 	err = collection.FindOne(context.TODO(), filter).Decode(&document)
// 	if err != nil {
// 		log.Errorf("GetCachedDocumentSummary-FindOne returned err: %v\n", err)
// 		return nil, err
// 	}
// 	return &document, nil
// }

// //Get Chached entries first. If none, get them from FHIR
// func (d *DocumentSummary) CachedDocument(filter bson.M) error {
// 	log.Debugf("DocumentGetCachedDocument using: %v\n", filter)

// 	collection, err := storage.GetCollection("documents")
// 	if err != nil {

// 		return err
// 	}
// 	log.Debugf("\nGetCachedDocument Checking for DocumentSummary with query: %v\n", filter)
// 	err = collection.FindOne(context.TODO(), filter).Decode(d)
// 	if err != nil {
// 		log.Errorf("CachedDocument-FindOne returned err: %v\n", err)
// 		return err
// 	}
// 	return nil
// }

func (f *DocumentFilter) QueryDiagCacheByEncounter() (bson.D, error) {
	if f.EncounterID == "" {
		return nil, fmt.Errorf("query_by_encounter has no encounter")
	}
	//q := bson.D{{"patient", f.PatientID}, {"encounter", f.Encounter}}
	q := bson.D{{"sessionid", f.SessionId}, {"patient", f.PatientGPI}, {"encounter", f.EncounterID}}
	return q, nil
}

// func (f *DocumentFilter) QueryCaDocumentCache() ([]*CADocument, error) {
// 	var limit int = 20
// 	var skip int = 0

// 	if f.Limit > 0 {
// 		limit = f.Limit
// 	}
// 	if f.Page > 0 {
// 		skip = (f.Page - 1) * limit
// 	}
// 	if f.Skip > 0 {
// 		skip = f.Skip
// 	}
// 	// q, err := f.QueryFhirCacheByEncounter()
// 	queryFilter := bson.M{"session_id": f.SessionId}
// 	log.Debugf("@@@  1494 -- Document QueryDiagCache Using filter %s", queryFilter)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(int64(limit))
// 	findOptions.SetSkip(int64(skip))

// 	// Multi sort fields are separated by [, ].
// 	// Order is based upon the order of the field names
// 	// sortFields := strings.Split(f.SortBy,", ")
// 	// for i, f := range sortFields {

// 	const DESC = -1
// 	const ASC = 1
// 	var sortFields bson.D
// 	var caDocs []*CADocument
// 	order := ASC // Default Assending

// 	if strings.ToLower(f.Order) == "desc" {
// 		order = DESC
// 	}

// 	sort := bson.E{}
// 	if f.Column == "" {
// 		f.Column = "rept_datetime"
// 	}
// 	sort = bson.E{f.Column, order}
// 	sortFields = append(sortFields, sort)
// 	if len(f.SortBy) > 0 {
// 		for _, s := range f.SortBy {
// 			if s == "visit_num" {
// 				sort = bson.E{"encounter", order}
// 			} else {
// 				sort = bson.E{s, order}
// 			}
// 			sortFields = append(sortFields, sort)
// 		}
// 	}
// 	findOptions.SetSort(sortFields)
// 	log.Debugf("@    1530 -- sort: %v", sortFields)
// 	collection, _ := storage.GetCollection("ca_documents")
// 	ctx := context.Background()
// 	cursor, err := collection.Find(ctx, queryFilter, findOptions)
// 	if err != nil {
// 		log.Debugf("QueryDiagCache:1536 -- for %s returned error: %v", queryFilter, err)
// 		return nil, err
// 	}
// 	defer func() {
// 		if err := cursor.Close(ctx); err != nil {
// 			log.WithError(err).Warn("Got error while closing Diag cursor")
// 		}
// 	}()
// 	for cursor.Next(ctx) {
// 		var caDoc CADocument
// 		err = cursor.Decode(&caDoc)
// 		if err != nil {
// 			log.WithError(err).Warn("Got error while closing cursor")
// 			return nil, err
// 		}
// 		caDocs = append(caDocs, &caDoc)
// 	}
// 	log.Debugf("   1553 -- Finished fetching documents from cursor")
// 	//cursor.Close(context.TODO())
// 	if len(caDocs) == 0 {
// 		err = fmt.Errorf("404|no DiagnosticReports found for %s", queryFilter)
// 		log.Error(err)
// 	} else {
// 		log.Debugf("QueryDiagCache found %d documents ", len(caDocs))
// 	}
// 	return caDocs, err
// }

func (rec *DocumentSummary) UpdateAccess() time.Time {
	filter := bson.M{"_id": rec.CacheID}
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
	encounterID := strings.Trim(f.EncounterID, " ")
	id := strings.Trim(f.ID, " ")
	enterpriseID := strings.Trim(f.EnterpriseID, " ")
	category := strings.Trim(f.Category, " ")
	reptDatetime := strings.Trim(f.ReptDatetime, " ")
	count := strings.Trim(f.Count, " ")
	//UseCache := strings.Trim(f.UseCache, "")

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
		log.Debugf("     @@@ Setting Count from config: %s", count)
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

	// log.Debugf("\n\nNew map: \n")

	f.queryMap = m
	log.Debugf("1652 -- SessionID: %s", f.SessionId)

	return nil
}

// Was named makeQUeryFilter. It is to hold the query of cache and renamed as such
func (df *DocumentFilter) makeCacheFilter() error {
	err := df.makeQueryMap()
	if err != nil {
		fmt.Printf("!@!@!@!    1662 -- makeQueryMap Failed: %s", err.Error())
		return err
	}
	df.MakeFhirQuery()
	//f.QueryFilter, _ = com.FilterFromMap(f.queryMap)
	log.Debugf("@#@#@#   1667 -- makeCacheFilter:1169 working with: %v", df.queryMap)
	layout := "2006-01-02"
	mq := []bson.M{}
	for k := range df.queryMap {
		val := df.queryMap[k]
		//log.Debugf("k: %s,  v: %s\n", k, q[k])
		if k == "count" {
			continue
		}
		// log.Debugf("\n    Document makeCacheFilter  Current Session: \n\n")

		if k == "sessionid" {
			log.Debugf("      1679 -- Adding seleted sessionID: %s", df.SessionId)
			mq = append(mq, bson.M{"sessionid": df.SessionId})
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
				log.Errorf("makeCacheFilter:1694 -- -Effective Date Error: %v", err)
				continue
			}
			q := bson.M{"effectivedate": bson.M{condition: useDate}}
			mq = append(mq, q)
		} else if k == "patient" {
			continue
		} else {
			log.Debugf("    1702 -- Adding %s to QueryFilter", k)
			mq = append(mq, bson.M{k: val})
		}
	}
	//f.QueryFilter = bson.M{}
	if len(mq) > 0 {
		df.CacheFilterBase = mq
		df.CacheFilterBase = append(df.CacheFilterBase, bson.M{"$and": mq})
	}
	log.Debugf("    1771 -- Final Document QueryFilter:1214 %s", spew.Sdump(df.CacheFilter))
	return nil
}

// was named makeQueryString whis the FHIR query string renamed to match

// func (f *DocumentFilter) makeFhirQuery() {
// 	fmt.Printf("\n\n\n       ### queryMap:1219 %v\n", f.queryMap)
// 	f.fhirQuery = fmt.Sprintf("?%s=%s", "patient", f.queryMap["patient"])
// 	if f.Count != "" {
// 		f.fhirQuery = fmt.Sprintf("%s&%s=%s", f.fhirQuery, "_count", f.queryMap["count"])
// 	}
// 	fmt.Printf("\n\n#### MakeFhirQuery:1221 f.FhirQuery: %s\n", f.fhirQuery)
// 	//f.queryString = ""
// 	// for k := range f.queryMap {
// 	// 	f.queryString = fmt.Sprintf("%s=%s", "patient", f.queryMap["patient"])
// 	// if f.queryString == "" {
// 	// 	f.queryString = fmt.Sprintf("%s=%s", k, f.queryMap[k])
// 	// } else {
// 	// 	f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, k, f.queryMap[k])
// 	// }
// 	// }
// }

// func (d *DocumentSummary) ByEnterpriseID() error {
// 	filter := bson.M{"enterprise_id": d.EnterpriseID}
// 	d, err := GetCachedDocumentSummary(filter)
// 	if err != nil { //Cache not found for ths Document
// 		return err
// 	}
// 	return nil
// }

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

// func ConvertDocumentsToVS(docs []*DocumentSummary) {
// 	startTime := time.Now()
// 	for _, d := range docs {
// 		d.ImageURL = fmt.Sprintf("%s%s", config.ImageURL(), d.EnterpriseID)
// 		// doc := d.ToCA()
// 		// caDocuments = append(caDocuments, doc)
// 	}
// 	log.Infof("Convert %d documents to VS took %s\n", len(docs), time.Since(startTime))

// }

// ProcessRemainingDocuments reads the next set of documents and if there is another link spins up
// another gorouting to process it while this one is filling the cache.

*/
// ////////////////////////////////////////////////////////////////////////////////////////////////////////
// //                              Process Next Set of DiagRepts returned from  FHIR                     //
// ////////////////////////////////////////////////////////////////////////////////////////////////////////

// func (df *DocumentFilter) FollowDiagNextLinks(links []fhir.Link) {
// 	// df.Session.UpdateDiagStatus("filling")
// 	// df.Session.Status.Diagnostic = "filling"
// 	fmt.Printf("\n\n\n### FollowNextLinks:2025\n")
// 	url := NextDiagPageLink(links)
// 	i := 1
// 	for {
// 		startTime := time.Now()
// 		if url == "" {
// 			log.Info("FollowDiagNextLinks: 2191 for url is blank, done")
// 			//df.Session.Status.Diagnostic = "done"
// 			//df.Session.UpdateDiagStatus( "done")
// 			break
// 		}
// 		//TODO: Get the next page and start its next page while processing current. Do in paraallel
// 		fmt.Printf("\nFilling Diagnostic Page %d\n", i)
// 		links, _ = df.ProcessDiagPage(url)
// 		fmt.Printf("    Page: %d  added in %f seconds\n\n", i, time.Since(startTime).Seconds())
// 		i = i + 1
// 		url = NextDiagPageLink(links)
// 	}
// 	df.Session.Status.Diagnostic = "done"
// 	df.Session.UpdateDiagStatus( "done")
// }

// // func NextDiagPage(links []fhir.Link) {
// // 	url := NextPageLink(links)
// // 	if url == "" {
// // 		return
// // 	}
// // }

// func NextDiagPageLink(links []fhir.Link) string {
// 	for _, l := range links {
// 		//fmt.Printf("Looking at link: %s\n", spew.Sdump(l))
// 		if l.Relation == "next" {
// 			fmt.Printf("##NextDiagLink:2059 %s \n\n", l.URL)
// 			return l.URL
// 		}
// 	}
// 	fmt.Printf("NextDiagPageLink:2063 - No next link\n")
// 	return ""
// }
// /*
// // func (df *DocumentFilter) ProcessCaDiagPage(url string) ([]fhir.Link, error) {
// // 	//fmt.Printf("\n\n\n###ProcessCaDiagPage:1998\n\n")
// // 	//startTime := time.Now()

// // 	fhirDiagResult, err := fhirC.NextFhirDiagRepts(url)
// // 	if err != nil {
// // 		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
// // 		return nil, err
// // 	}
// // 	//fmt.Printf("\n\nNext Page of results: %s\n\n", spew.Sdump(fhirPatients))
// // 	err = InsertFhirDiagResults(fhirDiagResult, df.Session.DocSessionId)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	diags := []*fhir.DiagnosticReport{}
// // 	entry := fhirDiagResult.Entry

// // 	for _, r := range entry {
// // 		rpt := r.DiagnosticReport

// // 		diags = append(diags, &rpt)
// // 	}

// // 	_, err = df.FhirDiagRptsToCADocuments(diags)
// // 	//_, err = FhirDiagReptsToCADocuments(fhirDiagRepts, df) //Documents are cached
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	//patients := parsePatientResults(fhirPatients, f.SessionId)
// // 	//spew.Dump()

// // 	if diags == nil {
// // 		return nil, fmt.Errorf("404|no more patients found")
// // 	}
// // 	return fhirDiagResult.Link, nil
// // }

// func (df *DocumentFilter) FollowNextDiagLinks(links []fhir.Link) {
// 	df.Session.UpdateDiagStatus("filling")
// 	df.Session.Status.Diagnostic = "filling"
// 	fmt.Printf("\n\n\n### FollowNextDiagLinks:2025\n")
// 	url := NextDiagPageLink(links)
// 	i := 1
// 	for {
// 		startTime := time.Now()
// 		if url == "" {
// 			log.Info("CachePages for url is blank, done")
// 			//df.Session.Status.Diagnostic = "done"
// 			//df.Session.UpdateDiagStatus( "done")
// 			break
// 		}
// 		//TODO: Get the next page and start its next page while processing current. Do in paraallel
// 		fmt.Printf("\nFilling Patient Page %d\n", i)
// 		links, _ = df.ProcessDiagPage(url)
// 		fmt.Printf("    Page: %d  added in %f seconds\n\n", i, time.Since(startTime).Seconds())
// 		i = i + 1
// 		url = NextDiagPageLink(links)
// 	}
// 	df.Session.Status.Diagnostic = "done"
// 	df.Session.UpdateDiagStatus( "done")
// }
// */
// func (df *DocumentFilter) ProcessDiagPage(url string) ([]fhir.Link, error) {
// 	//fmt.Printf("\n\n\n###ProcessCaDiagPage:1998\n\n")
// 	//startTime := time.Now()

// 	fhirDiagResult, err := fhirC.NextFhirDiagRepts(url)
// 	if err != nil {
// 		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
// 		return nil, err
// 	}
// 	//fmt.Printf("\n\nNext Page of results: %s\n\n", spew.Sdump(fhirPatients))
// 	docs, err := InsertFhirDocResults(fhirDiagResult, df.Session.DocSessionId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// docs := []*fhir.Document{}
// 	// entry := fhirDiagResult.Entry

// 	// for _, r := range entry {
// 	// 	rpt := r.Document

// 	// 	docs = append(docs, &rpt)
// 	// }

// 	//_, err = df.FhirDiagRptsToCADocuments(diags)
// 	//_, err = FhirDiagReptsToCADocuments(fhirDiagRepts, df) //Documents are cached
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//patients := parsePatientResults(fhirPatients, f.SessionId)
// 	//spew.Dump()

// 	if docs == nil {
// 		return nil, fmt.Errorf("404|no more patients found")
// 	}
// 	return fhirDiagResult.Link, nil
// }

// func InsertFhirDocResults(results *fhir.DocumentResults, sessionId string)  ([]*fhir.Document, error) {
// 	log.Info("InsertFhirDocResults: 2325")
// 	//entry := results.Entry
// 	docs := []*fhir.Document{}
// 	for _, entry := range results.Entry {

// 		doc := entry.Document
// 		doc.FullURL = entry.FullURL
// 		docs = append(docs, &doc)
// 		err := InsertFhirDoc(&doc, sessionId)
// 		if err != nil {
// 			msg := fmt.Sprintf("InsertFhirDocResults:2333 --  failed: %s", err.Error())
// 			log.Error(msg)
// 			return nil, errors.New(msg)
// 		}
// 	}
// 	return docs, nil
// }

// func InsertFhirDocs(docs []*fhir.Document, sessionId string) error {
// 	for _, doc := range docs {
// 		err := InsertFhirDoc(doc, sessionId)
// 		if err != nil {
// 			if !storage.IsDup(err) {
// 				msg := fmt.Sprintf("InsertFhirDiags:2346 --   failed: %s", err.Error())
// 				log.Error(msg)
// 				return errors.New(msg)
// 			}
// 			err = nil
// 		}
// 	}
// 	return nil
// }

// func InsertFhirDoc(doc *fhir.Document, sessionId string) error {

// 	doc.SessionID = sessionId
// 	// _, err := FindByPhone(c.FaxNumber, c.Facility)
// 	// log.Fatal(err)
// 	collection, _ := storage.GetCollection("documents")
// 	//diag.CacheId = primitive.NewObjectID()
// 	doc.CacheID = primitive.NewObjectID()
// 	_, err := collection.InsertOne(context.TODO(), doc)
// 	if err != nil {
// 		if !storage.IsDup(err) {
// 			msg := fmt.Sprintf("InsertFhirDiag:2373 -- Insert Error: %s", err.Error())
// 			log.Error(msg)
// 			return errors.New(msg)
// 		}
// 		err = nil
// 	}
// 	return err
// }
/*
// func InsertFhirDocument(doc *fhir.Document, sessionId string) error {
// 	log.Info("InsertFhirDocument: 2362")

// 	doc.SessionID = sessionId
// 	// _, err := FindByPhone(c.FaxNumber, c.Facility)
// 	// log.Fatal(err)
// 	collection, _ := storage.GetCollection("documents")
// 	//diag.CacheId = primitive.NewObjectID()
// 	doc.CacheID = primitive.NewObjectID()
// 	_, err := collection.InsertOne(context.TODO(), doc)
// 	if err != nil {
// 		if !storage.IsDup(err) {
// 			msg := fmt.Sprintf("InsertFhirDoc:2394 -- Insert Error: %s", err.Error())
// 			log.Error(msg)
// 			return errors.New(msg)
// 		}
// 		err = nil
// 	}
// 	return err
// }
// func (f *DocumentFilter) ProcessRemainingDiagRepts(link string) {
// 	err := f.ReadNextSet(link)
// 	if err != nil {
// 		log.Errorf("!!!   ReadNextSet returned error: %v", err)
// 	}
// 	//ds := BuildDocumentReferences(ref, &session)
// }

// func (df *DocumentFilter) ReadNextSet(url string) error {
// 	//url = "https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/Binary/"https://fhir-open.sandboxcerner.com/dstu2/0b8a0111-e8e6-4c26-a91c-5069cbc6b1ca/"

// 	//TODO:ReadNextSet should possibly use fhirongo:GetFhir
// 	startTime := time.Now()
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	req.Header.Set("Accept", "application/json+fhir")
// 	resp, err := client.Do(req)
// 	//resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()
// 	log.Debugf("Status Code: %d - %s", resp.StatusCode, resp.Status)
// 	if resp.StatusCode > 299 || resp.StatusCode < 200 {
// 		err = fmt.Errorf("%d|%v", resp.StatusCode, resp.Status)
// 		return err
// 	}
// 	var diagResp fhir.DiagnosticReportResponse

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.Body)
// 	respByte := buf.Bytes()

// 	if err := json.Unmarshal(respByte, &diagResp); err != nil {
// 		return err
// 	}
// 	log.Infof("$  ReadNextSet:1991 -- set of documents took %s", time.Since(startTime))
// 	for _, l := range diagResp.Link {
// 		if l.Relation == "next" {
// 			nextLink := l.URL
// 			//log.Debugf("NextLink: %s\n", nextLink)
// 			go f.ProcessRemainingDiagRepts(nextLink)
// 			break
// 		} // if no next link there are no more DiagRepts
// 	}

// 	entries := diagResp.Entry
// 	diagRepts := []*fhir.DiagnosticReport{}
// 	for _, ent := range entries {
// 		//rpt := fhir.DiagnosticReport{}
// 		rpt := ent.DiagnosticReport
// 		rpt.FullURL = ent.FullURL
// 		log.Debugf("ReadNextSet: - Report.ID: %s", rpt.ID)
// 		diagRepts = append(diagRepts, &rpt)
// 	}

// 	// f.FindFhirDiagnosticReports
// 	_, err = df.FhirDiagRptsToCA(diagRepts)
// 	//err = f.BuildDocumentReferences(&ref, session)
// 	return err
// 	//return &buf, err
// }

/////////////////////////////////////////////////////////////////////////////////////////
//                                Cache Routines                                        /
/////////////////////////////////////////////////////////////////////////////////////////

// func (f *DocumentFilter) QueryCACache() ([]*fhir.DiagnosticReport, error) {
// 	var limit int = 20
// 	var skip int = 0

// 	if f.Limit > 0 {
// 		limit = f.Limit
// 	}
// 	if f.Page > 0 {
// 		skip = (f.Page - 1) * limit
// 	}
// 	if f.Skip > 0 {
// 		skip = f.Skip
// 	}
// 	// q, err := f.QueryFhirCacheByEncounter()

// 	log.Debugf("@@@   1860 -- Document QueryFhirCache Using filter %s", f.QueryFilter)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(int64(limit))
// 	findOptions.SetSkip(int64(skip))

// 	// Multi sort fields are separated by [, ].
// 	// Order is based upon the order of the field names
// 	// sortFields := strings.Split(f.SortBy,", ")
// 	// for i, f := range sortFields {

// 	const DESC = -1
// 	const ASC = 1
// 	var sortFields bson.D
// 	var diagRepts []*fhir.DiagnosticReport
// 	var documents []*DocumentSummary
// 	order := ASC // Default Assending

// 	if strings.ToLower(f.Order) == "desc" {
// 		order = DESC
// 	}

// 	sort := bson.E{}
// 	if f.Column == "" {
// 		f.Column = "rept_datetime"
// 	}
// 	sort = bson.E{f.Column, order}
// 	sortFields = append(sortFields, sort)
// 	if len(f.SortBy) > 0 {
// 		for _, s := range f.SortBy {
// 			if s == "visit_num" {
// 				sort = bson.E{"encounter", order}
// 			} else {
// 				sort = bson.E{s, order}
// 			}
// 			sortFields = append(sortFields, sort)
// 		}
// 	}
// 	findOptions.SetSort(sortFields)
// 	log.Debugf("@    1897 -- sort: %v", sortFields)

// 	// }

// 	collection, _ := storage.GetCollection("diag_repts")
// 	ctx := context.Background()
// 	cursor, err := collection.Find(ctx, f.QueryFilter, findOptions)
// 	if err != nil {
// 		log.Debugf("QueryFhirCache:1906 -- for %s returned error: %v", f.QueryFilter, err)
// 		//cursor.Close(ctx)
// 		return nil, err
// 	}
// 	defer func() {
// 		if err := cursor.Close(ctx); err != nil {
// 			//log.Debugf("document queryCache error while closing cursor: %v\n", err)
// 			log.WithError(err).Warn("Got error while closing cursor")
// 		}
// 	}()
// 	//log.Debugf("\n    No Error on QueryFhirCache\n\n")
// 	for cursor.Next(ctx) {
// 		var diagRept fhir.DiagnosticReport
// 		var document DocumentSummary
// 		err = cursor.Decode(&document)
// 		if err != nil {
// 			//log.Debugf("   Next error: %v\n", err)
// 			//cursor.Close(context.TODO())
// 			log.WithError(err).Warn("Got error while closing cursor")
// 			return nil, err
// 		}
// 		//log.Debugf("  Added one to documents\n")
// 		diagRepts = append(diagRepts, &diagRept)
// 	}
// 	log.Debugf("   Finished fetching documents from cursor")
// 	//cursor.Close(context.TODO())
// 	if len(diagRepts) == 0 {
// 		err = fmt.Errorf("404|no documents found for %s", f.QueryFilter)
// 		log.Error(err)
// 	} else {
// 		log.Debugf("QueryFhirCache:1936 -- found %d documents ", len(documents))
// 	}
// 	return diagRepts, err
// }

func (df *DocumentFilter) QueryFhirCache() ([]*fhir.DiagnosticReport, error) {
	var limit int64 = 20
	var skip int64 = 0

	if df.Limit > 0 {
		limit = df.Limit
	}
	if df.Page > 0 {
		skip = (df.Page - 1) * limit
	}
	if df.Skip > 0 {
		skip = df.Skip
	}
	// q, err := f.QueryFhirCacheByEncounter()

	log.Debugf("@@@   1956 -- Document QueryFhirCache Using filter %s", df.CacheFilter)
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	// Multi sort fields are separated by [, ].
	// Order is based upon the order of the field names
	// sortFields := strings.Split(f.SortBy,", ")
	// for i, f := range sortFields {

	const DESC = -1
	const ASC = 1
	var sortFields bson.D
	var diagRepts []*fhir.DiagnosticReport
	var documents []*DocumentSummary
	order := ASC // Default Assending

	if strings.ToLower(df.Order) == "desc" {
		order = DESC
	}

	sort := bson.E{}
	if df.Column == "" {
		df.Column = "rept_datetime"
	}
	sort = bson.E{df.Column, order}
	sortFields = append(sortFields, sort)
	if len(df.SortBy) > 0 {
		for _, s := range df.SortBy {
			if s == "visit_num" {
				sort = bson.E{"encounter", order}
			} else {
				sort = bson.E{s, order}
			}
			sortFields = append(sortFields, sort)
		}
	}
	findOptions.SetSort(sortFields)
	log.Debugf("@    1993 -- sort: %v", sortFields)

	// }

	collection, _ := storage.GetCollection("diag_repts")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, df.CacheFilter, findOptions)
	if err != nil {
		log.Debugf("QueryFhirCache:2002 -- for %s returned error: %v", df.CacheFilter, err)
		//cursor.Close(ctx)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			//log.Debugf("document queryCache error while closing cursor: %v\n", err)
			log.WithError(err).Warn("Got error while closing cursor")
		}
	}()
	//log.Debugf("\n    No Error on QueryFhirCache\n\n")
	for cursor.Next(ctx) {
		var diagRept fhir.DiagnosticReport
		var document DocumentSummary
		err = cursor.Decode(&document)
		if err != nil {
			//log.Debugf("   Next error: %v\n", err)
			//cursor.Close(context.TODO())
			log.WithError(err).Warn("Got error while closing cursor")
			return nil, err
		}
		//log.Debugf("  Added one to documents\n")
		diagRepts = append(diagRepts, &diagRept)
	}
	log.Debugf("   2036 --Finished fetching documents from cursor")
	//cursor.Close(context.TODO())
	if len(diagRepts) == 0 {
		err = fmt.Errorf("404|no documents found for %s", df.CacheFilter)
		log.Error(err)
	} else {
		log.Debugf("QueryFhirCache:2032 -- found %d documents", len(documents))
	}
	return diagRepts, err
}

// DeleteCache:  Deletes the cahe for the
//func (df *DocumentFilter) DeleteCache()

func DeleteDocuments(docSessionId string) {

	startTime := time.Now()
	log.Debugf("@@@!!!  2040 -- Deleting Documents for session %s", docSessionId)
	collection, _ := storage.GetCollection("documents")
	filter := bson.M{"sessionid": docSessionId}
	log.Debugf("DeleteDocument:2227 -- bson filter delete: %v\n", filter)
	deleteResult, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Errorf("DeleteDocuments:2230 -- for session %s failed: %v", docSessionId, err)
		return
	}
	log.Infof("DeleteDocuments:2233 -- Deleted %v Documents for session: %v in %s",
				deleteResult.DeletedCount, docSessionId, time.Since(startTime))
}
*/
// /////////////////////////////////////////////////////////////////////////////////////////
// //                                 FHIR Getters                                         /
// /////////////////////////////////////////////////////////////////////////////////////////

// func GetImage(fd *fhir.DiagnosticReport, imageType string) string {
// 	for _, form := range fd.PresentedForm {
// 		switch form.ContentType {
// 		case imageType:
// 			return form.URL
// 		}
// 	}
// 	return ""
// }

// func GetFhirPerson(per fhir.Person, item string) string {
// 	if item == "ID" {
// 		id := strings.Split(per.Reference, "/")
// 		if len(id) > 1 {
// 			return id[1]
// 		} else {
// 			return ""
// 		}
// 	} else {
// 		return per.Display
// 	}
// }

// func GetFhirReference(per fhir.Reference) string {

// 	id := strings.Split(per.Reference, "/")
// 	if len(id) > 1 {
// 		return id[1]
// 	} else {
// 		return ""
// 	}
// }

// // SplitReference: Accepts a string of ssss/dddd and returns the second part
// func SplitReference(ref string) string {
// 	id := strings.Split(ref, "/")
// 	if len(id) > 1 {
// 		return id[1]
// 	} else {
// 		return ""
// 	}
// }

/*
//func GetFhirEncounterID(enc fhir.EncounterReference) string {
// 	GetFhirReference(enc.Reference)
// 	id := strings.Split(enc.Reference, "/")
// 	if len(id) > 1 {
// 		return id[1]
// 	} else {
// 		return ""
// 	}
// 	return GetFhirReference(enc.Reference)
// }

// func BuildDocumentReferences(docR *fhir.DocumentReferences, sessionp *AuthSession) error {
// 	// *DocumentSummary {
// 	// fmt.Println("\n\n================docR===")

// 	session := *sessionp
// 	cacheName := session.SessionID
// 	// ///dbSession, err := storage.GetSession()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	//var docs []*DocumentSummary
// 	//var nextLink string
// 	startTime := time.Now()
// 	for _, l := range docR.Link {
// 		if l.Relation == "next" {
// 			//nextLink = l.URL
// 			//log.Debugf("NextLink: %s\n", nextLink)
// 			//go ProcessRemainingDocuments(nextLink, session)
// 			break
// 		}
// 	}
// 	numDocs := 0
// 	for _, entry := range docR.Entry {
// 		var doc DocumentSummary

// 		doc.Source = ActiveConfig().Source()
// 		doc.SourceType = "Reference"
// 		doc.Text = entry.DocumentReference.Text
// 		//doc.Code = entry.Resource.Code.Text
// 		// 	doc.Category = entry.Resource.Category.Text
// 		doc.ReptDatetime = entry.DocumentReference.Created
// 		doc.FullLink = entry.FullURL
// 		doc.EnterpriseID = entry.DocumentReference.ID
// 		doc.Subject = entry.DocumentReference.Subject
// 		//TODO: check both patient and subject for the patient information
// 		doc.PatientID = strings.Split(doc.Subject.Reference, "/")[1]
// 		doc.PatientGPI = doc.PatientID
// 		doc.Performer = entry.DocumentReference.Authenticator
// 		enc := strings.Split(entry.DocumentReference.Context.Encounter.Reference, "/")
// 		if len(enc) > 1 {
// 			doc.Encounter = enc[1]
// 		}
// 		doc.Description = entry.DocumentReference.Description
// 		doc.Category = entry.DocumentReference.Type.Text
// 		for _, item := range entry.DocumentReference.Content {
// 			//log.Debugf("Attachment: %v\n", item.Attachment.URL)
// 			if item.Attachment.ContentType == "application/pdf" {
// 				doc.FhirImageURL = item.Attachment.URL
// 				doc.makeImageURL()
// 			}
// 		}
// 		doc.SessionID = cacheName
// 		// doc.setDates()
// 		// (&doc).insertWithSession(&dbSession)
// 		//(&doc).Insert(cacheName, &dbSession)
// 		numDocs++
// 	}
// 	log.Infof("#    Cached %d documents in %s\n", numDocs, time.Since(startTime))
// 	return nil
// }

func (df *DocumentFilter) FindLatestCachedFhir() (*DocumentSummary, error) {
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

	mq := append(df.QueryFilterBase, bson.M{"source": df.Source})
	filter := bson.M{"$and": mq}

	sort := bson.E{"rept_datetime", order}
	sortFields = append(sortFields, sort)
	findOptions.SetSort(sortFields)
	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Debugf("QueryFhirCache for %s returned error: %v\n", df.CacheFilter, err)
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
		err = fmt.Errorf("404|no documents found for %s", df.CacheFilter)
		log.Error(err)
	} else {
		log.Debugf("QueryLatest found %d documents \n", len(documents))
		document = documents[0]
	}
	return document, err
}

func (df *DocumentFilter) MakeFhirQuery() {
	qry := "patient=" + df.PatientGPI
	if df.BeginDate != "" {
		mdyDate, _ := common.MDYToFhir(df.BeginDate)
		bDate := common.FhirDateToString(mdyDate, "full")
		qry = fmt.Sprintf("%s&date=ge%s", qry, bDate)
	}
	if df.EndDate != "" {
		mdyDate, _ := common.MDYToFhir(df.EndDate)
		eDate := common.FhirDateToString(mdyDate, "full")
		qry = fmt.Sprintf("%s&date=lt%s", qry, eDate)
	}
	if df.Count != "" {
		qry = fmt.Sprintf("%s&_count=%s", qry, df.Count)
	}
	//fmt.Printf("\n\n####  Query: %s\n", qry)
	df.fhirQuery = qry
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
