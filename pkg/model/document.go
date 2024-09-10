package model

import (
	"context"
	"errors"
	"fmt"
	"strings"

	//"time"

	storage "github.com/dhf0820/cernerFhir/pkg/storage"
	log "github.com/sirupsen/logrus"

	//"time"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentFilter struct {
	Skip          int64    `schema:"skip"`
	Page          int64    `schema:"page"`
	Limit         int64    `schema:"limit"`
	SortBy        []string `schema:"sort"`
	Column        string   `schema:"column"`
	Order         string   `schema:"order"`
	ResultFormat  string   `schema:"format"` //Header
	Count         string   `schema:"count"`  //Header
	UserId        string
	Cache         string `schema:"cache"` //reset, stats
	Session       *AuthSession
	SessionId     string `Schema:"session_id"`
	AccessDetails *AccessDetails

	Source string
	// Specified filters.
	// PatientID/ID is the only one to query FHIR with
	// Once all documents are in the mongo cache we do the requested query using mongo query

	//VisitID     string `schema:"visit_num"`
	EncounterID string `schema:"visit_num"`
	//Visit      string `schema:"visit"`
	PatientID  string `schema:"patient_id"`
	PatientGPI string `schema:"patient_gpi"`

	//Encounter    string `schema:"encounter"`
	MRN          string `schema:"mrn"`
	ID           string `schema:"id"`
	DocID        string `schema:"doc_id"`
	EnterpriseID string `schema:"enterpriseid"`
	ReptDatetime string `schema:"rept_datetime"`
	Category     string `schema:"category"`

	SourceValues []string `schema:"source_values"`
	BeginDate    string   `schema:"begin_date"`
	EndDate      string   `schema:"end_date"`
	TabID        string   `schema:"tab_id"`

	//queryString     string
	fhirQuery string
	//queryMap        map[string]string
	CacheFilterD    bson.D
	CacheFilter     bson.M
	QueryFilterBase []bson.M
	CacheFilterBase []bson.M
}

// type Document struct {
// 	CacheID 		primitive.ObjectID			`bson:"cache_id" json:"cacheId"`
// 	SessionID		string						`bson:"session_id" json:"sessionId"`
// 	ResourceType 	string						`bson:"resource_id" json:"resourceType"`
// 	ID 				string						`bson:"id" json:"id"`
// 	FullURL			string						`bson:"full_url" json:"fullURL"`
// 	EffectiveDateTime time.Time 				`bson:"effective_date_time", json:"effectiveDateTime"`
// 	Meta 			fhir.MetaData 				`bson:"meta" json:"meta"`
// 	Text			fhir.TextData 				`bson:"text" json:"text"`
// 	Status			string						`bson:"status" json:"status"`
// 	Category 		fhir.CodeableConcept 		`bson:"category" json:"category"`
// 	Code 			fhir.CodeableConcept		`bson:"code" json:"code"`
// 	Subject			fhir.Person					`bson:"subject" json:"subject"`
// 	Type			fhir.Concept				`bson:"type" json:"type"`
// 	Encounter 		fhir.EncounterReference		`bson:"encounter" json:"encounter"`
// 	Issued 			time.Time 					`bson:"issued" json:"issued"`
// 	Performer		Person						`bson:"performer" json:"performer`
// 	PresentedForm	[]fhir.Attachment			`bson:"presented_form" json:"presentedForm"`
// 	Request			fhir.Thing					`bson:"request" json:"request"`
// 	Result 			fhir.Thing 					`bson:"result" json:"result"`
// 	Authenticator	Person						`bson:"authenticator" json:"authenticator"`
// 	Created			time.Time					`bson:"created" json:"created"`
// 	Indexed 		time.Time					`bson:"indexed" json:"indexed"`
// 	DocStatus		fhir.Concept				`bson:"doc_status" json:"docSatus"`
// 	Description 	string						`bson:"description" json:"description"`
// 	Context 		fhir.EncounterContext 		`bson:"context" json:"context"`
// 	Content 		[]fhir.Attachment			`bson:"content" json:"content"`
// }

// SearchReports returns  []*fhir.Document, cacheStatus(string), numInPage(int64), pagesInCache(int64), totalInCache(int64), error
func (df *DocumentFilter) SearchReports() ([]*fhir.Document, string, int64, int64, int64, error) {
	fhirC = config.Fhir()
	//activeDocumentFilter = df

	// df.Session.Status.Reference = "filling"
	// df.Session.UpdateRefStatus( "filling")
	// df.Session.UpdateRefStatus( "filling")
	// Start GetFhirDiags

	status, pagesInCache, totalInCache, err := df.DocumentCacheStats()
	if err != nil {
		log.Errorf("SearchReports:107 -  error: %s", err.Error())
		return nil, "", 0, 0, 0, err
	}
	log.Debugf("SearchReports:110 -- Status: %s, PagesInCache: %d, TotalInCache: %d\n", status, pagesInCache, totalInCache)
	if df.Page <= pagesInCache {
		log.Debugf("SearchReports:112 -- Enough pages  to fill page: %d request", df.Page)
		return df.GetFhirDocumentPage()
	}

	// if status == "done" {

	// 	return nil, status, 0, pagesInCache, totalInCache, errors.New("page not available")
	// }

	// Fill the Cache
	//fmt.Printf("\n\n#### Calling GetDiagnosticRept\ns")
	go df.FindFhirDiagRepts()

	// Start GetDocRef
	//fmt.Printf("\n\n### Calling DocumentReferences\n\n")
	df.GetFhirDocRefs()

	//Wait for either both done status or 20 in the cache

	// for {
	// 	// is status Done for both?
	// 	status, _, inCache, err := df.DocumentCacheStats()
	// 	if err != nil {
	// 		msg := fmt.Sprintf("SearchReports:110 - DocumentCacheStats %s", err.Error())
	// 		log.Errorf("%s", msg)
	// 		return nil, "", 0, 0, 0, errors.New(msg)
	// 	}

	// 	if status == "done" {
	// 		fmt.Printf("\n\n### Both are done\n\n")
	// 		//if it is get and return the page requested
	// 		//fDocs, cacheStatus, inPage, pagesInCache, totalInCache, err :=  df.GetFhirDocumentPage()
	// 		return df.GetFhirDocumentPage()
	// 	}
	// 	if inCache >= 20  {
	// 		fmt.Printf("\n\n\n### Enough documents: %d in cache\n\n", inCache)
	// 		return df.GetFhirDocumentPage()
	// 	}
	// 	time.Sleep(400* time.Millisecond)
	// 	fmt.Printf("\n\n### SearchReports:129 -- Currently %d in cache\n", inCache)
	// }
	//fmt.Printf("######SearchReports: 153 -- Returning requested Page\n")
	return df.GetFhirDocumentPage()
}

// _, err := df.FindFhirDiagRepts()
// if err != nil {
// 	return nil, "", 0, 0, 0, err
// }

// fmt.Printf("\n\n### Calling DocumentReferences\n\n")
// _, err = df.GetFhirDocRefs()
// //_, err = df.GetCaDocumentReferences()
// if err != nil {
// 	log.Warnf("SearchReports:190 GetDocRefs returned error: %s", err.Error())
// }
// // fhirReptRef, err := df.FindFhirDo()
// // if err != nil {
// // 	return nil, err
// // }
// // caRepts := FhirDiagRptsToCA(fhirDiag)

// //TODO: Cache the CA document version

// //return fhirDiag, nil
// //log.Debugf("CacheQuery:212 -- %v", f.CacheFilterBase)
// totalInCache, err := df.CountCachedFhirDocuments()
// if err != nil {
// 	return nil, "", 0, 0, 0, err
// }
// log.Infof("Total Cached Documents: %d", totalInCache)
// //log.Debugf("Counted:217 -- %d documents", totalDocs)
// //docs, cacheStatus, numInPage, pagesInCache, totalInCache, err := df.GetDocumentPage()
// // if err != nil {
// // 	return nil, "", 0, 0, 0, err
// // }
// //log.Debugf("SearchCAReports:206 -- returning %d Documents", len(caRepts))

// return nil, "", 0, 0, 0, nil
// //return nil, cacheStatus, numInPage, pagesInCache, totalInCache, nil

func (df *DocumentFilter) MakeFhirQuery() string {
	qry := "patient=" + df.PatientGPI

	if df.BeginDate != "" {
		mdyDate, _ := MDYToFhir(df.BeginDate)
		bDate := FhirDateToString(mdyDate, "full")
		qry = fmt.Sprintf("%s&date=ge%s", qry, bDate)
	}
	if df.EndDate != "" {
		mdyDate, _ := MDYToFhir(df.EndDate)
		eDate := FhirDateToString(mdyDate, "full")
		qry = fmt.Sprintf("%s&date=lt%s", qry, eDate)
	}
	if df.Count != "" {
		qry = fmt.Sprintf("%s&_count=%s", qry, df.Count)
	}
	df.fhirQuery = qry
	return qry
}

func InsertFhirDocResults(results *fhir.DocumentResults, sessionId string) ([]*fhir.Document, error) {
	//entry := results.Entry
	docs := []*fhir.Document{}
	for _, entry := range results.Entry {

		doc := entry.Document
		doc.FullURL = entry.FullURL
		docs = append(docs, &doc)
		err := InsertFhirDoc(&doc, sessionId)
		if err != nil {
			msg := fmt.Sprintf("InsertFhirDocResults:224 --  failed: %s", err.Error())
			log.Error(msg)
			return nil, errors.New(msg)
		}
	}
	return docs, nil
}

func InsertFhirDocs(docs []*fhir.Document, sessionId string) error {
	for _, doc := range docs {
		err := InsertFhirDoc(doc, sessionId)
		if err != nil {
			if !storage.IsDup(err) {
				msg := fmt.Sprintf("InsertFhirDiags:237 --   failed: %s", err.Error())
				log.Error(msg)
				return errors.New(msg)
			}
			err = nil
		}
	}
	return nil
}

func InsertFhirDoc(doc *fhir.Document, sessionId string) error {

	doc.SessionID = sessionId
	// _, err := FindByPhone(c.FaxNumber, c.Facility)
	// log.Fatal(err)
	collection, _ := storage.GetCollection("documents")
	//diag.CacheId = primitive.NewObjectID()
	doc.CacheID = primitive.NewObjectID()

	enc := doc.Context.Encounter.Reference
	doc.Encounter.Reference = enc
	_, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		if !storage.IsDup(err) {
			msg := fmt.Sprintf("InsertFhirDiag:259 -- Insert Error: %s", err.Error())
			log.Error(msg)
			return errors.New(msg)
		}
		err = nil
	}
	return err
}

func (df *DocumentFilter) CountCachedFhirDocuments() (int64, error) {

	filter := bson.M{"subject.reference": "Patient/" + df.PatientGPI} //f.CacheFilterBase //append(f.QueryFilterBase, bson.M{"session_id": f.SessionId})
	log.Debugf("CountCachedFhirdDocuments:271 --  matching: %v\n", filter)
	collection, err := storage.GetCollection("documents")
	if err != nil {
		return -1, err
	}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("CountchedFhirDocuments:278 -- returned error: %v\n", err)
		return 0, err
	}
	return count, nil
}

// GetFhirDocumentPage: returns fDocs, cacheStatus, int64(len(fDocs)), pagesInCache, totalInCache, error

func (df *DocumentFilter) GetFhirDocumentPage() ([]*fhir.Document, string, int64, int64, int64, error) {
	fmt.Printf("GetFhirDocumentPage:291")
	var linesPerPage int64 = LinesPerPage()
	var skip int64 = 0
	//var caDocs []*CADocument
	var fDocs []*fhir.Document
	const DESC = -1
	const ASC = 1
	cacheFilter := bson.M{}

	if df.Limit > 0 {
		//fmt.Printf("GetFhirDocumentPage:301 - setting linesPerPage: %d\n", df.Limit)
		linesPerPage = df.Limit
	}
	if df.Page > 0 {
		skip = (df.Page - 1) * linesPerPage
		//fmt.Printf("GetFhirDocumentPage:306 -- setting skip: %d\n", skip)
	}
	if df.Skip > 0 {
		skip = df.Skip
	}
	mq := []bson.M{}
	fmt.Println()
	fmt.Println("###")
	fmt.Println()
	mq = append(mq, bson.M{"subject.reference": "Patient/" + df.PatientGPI})
	log.Debugf("GetFhirDocumentPage:316 -- mq: %v", mq)
	if df.EncounterID != "" {
		mq = append(mq, bson.M{"context.encounter.reference": "Encounter/" + df.EncounterID})
		log.Debugf("GetFhirDocumentPage:319 -- mq: %v", mq)
	}
	if len(mq) > 1 {

		cacheFilter = bson.M{"$and": mq}
	}
	log.Debugf("GetFhirDocumentPage:325 -- filter: %v", cacheFilter)
	findOptions := options.Find()
	findOptions.SetLimit(linesPerPage)
	findOptions.SetSkip(skip)
	sortOrder := ASC // Default Assending
	var sortFields bson.D

	if strings.ToLower(df.Order) == "desc" {
		sortOrder = DESC
	}

	//sortFields = append(sortFields, bson.E{"rept_datetime", sortOrder})

	sort := bson.E{}
	if df.Column == "" {
		df.Column = "rept_datetime"
	}
	sort = bson.E{df.Column, sortOrder}
	sortFields = append(sortFields, sort)
	if len(df.SortBy) > 0 {
		for _, s := range df.SortBy {
			if s == "visit_num" {
				sort = bson.E{"encounter", sortOrder}
			} else {
				sort = bson.E{s, sortOrder}
			}
			sortFields = append(sortFields, sort)
		}
	}
	findOptions.SetSort(sortFields)

	log.Debugf("GeFhirDocumentPage:356 - cacheFilter: %v", cacheFilter)

	findOptions.SetSort(sortFields)
	log.Debugf("GetFhirDocumentPage:359-- sort: %v", sortFields)
	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, cacheFilter, findOptions)
	if err != nil {
		log.Debugf("GetFhirDocumentPage364 for filter: %s returned error: %s\n", cacheFilter, err.Error())
		return nil, "", 0, 0, 0, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.WithError(err).Warnf("GetFhirDocumentPage:369 -- Got error while closing Document cursor: %s", err.Error())
		}
	}()
	for cursor.Next(ctx) {
		var fDoc fhir.Document
		err = cursor.Decode(&fDoc)
		if err != nil {
			log.WithError(err).Warnf("GetFhirDocumentPage:376 -- Got error while closing cursor: %s", err.Error())
			return nil, "", 0, 0, 0, err
		}
		fDocs = append(fDocs, &fDoc)
	}
	//log.Debugf("GetFhirDocumentPage:381 -- Finished fetching %d documents from cursor", len(fDocs))
	//cursor.Close(context.TODO())
	if len(fDocs) == 0 {
		log.Infof("GetFhirDocumentPage384 -- no Documents found for %s", cacheFilter)

	}
	//else {
	// 	log.Debugf("GetFhirDocumentPage:388 -- GetQueryDiagCache found %d documents \n", len(fDocs))
	// }
	cacheStatus, pagesInCache, totalInCache, err := df.DocumentCacheStats()
	return fDocs, cacheStatus, int64(len(fDocs)), pagesInCache, totalInCache, err

}

// Get document by fhirId
func FhirDocumentById(id string) (*fhir.Document, error) {
	fDoc := &fhir.Document{}
	filter := bson.M{"id": id}
	findOneOptions := options.FindOneOptions{}

	collection, _ := storage.GetCollection("documents")
	ctx := context.Background()
	results := collection.FindOne(ctx, filter, &findOneOptions)
	if results.Err() != nil {
		return nil, results.Err()
	}
	err := results.Decode(fDoc)
	return fDoc, err
}

//DocumentCacheStats returns cacheStatus, pagesInCache, totalInCache, error
func (df *DocumentFilter) DocumentCacheStats() (string, int64, int64, error) {
	totalInCache, err := df.DocumentsInCache()
	if err != nil {
		msg := fmt.Sprintf("DocumentCacheStats:399 -- err: %s", err.Error())
		return "", 0, 0, errors.New(msg)
	}
	pageSize := LinesPerPage()
	pagesInCache, _ := CalcPages(totalInCache, pageSize)
	//pages := inCache/pageSize
	log.Debugf("DocumentCacheStats:405 -- pageSize: %d  InCache: %d, pagesInCaches: %d", pageSize, totalInCache, pagesInCache)
	cacheStatus := df.Session.GetDocumentStatus()
	return cacheStatus, pagesInCache, totalInCache, nil
}

func (df *DocumentFilter) DocumentPagesInCache() (int64, error) {
	numInCache, err := df.DocumentsInCache()
	if err != nil {
		return 0, err
	}
	pageSize := LinesPerPage()
	pagesInCache, _ := CalcPages(numInCache, pageSize)
	//pages := inCache/pageSize
	log.Debugf("DocumentPagesInCache:418 -- pageSize: %d  InCache: %d, pagesInCaches: %d", pageSize, numInCache, pagesInCache)
	return int64(pagesInCache), nil
}

func (df *DocumentFilter) DocumentsInCache() (int64, error) {
	//mq := append(f.CacheFilterBase, bson.M{"source": f.Source})
	//filter := bson.M{"$and": mq}
	c, err := storage.GetCollection("documents")
	if err != nil {
		log.Errorf("DocumentsInCache:4427 -- Settng Collection(%s) failed: %s", "documents", err.Error())
	}
	filter := bson.M{"subject.reference": "Patient/" + df.PatientGPI}
	log.Infof("DocumentsInCache:430 For Patient matching: [%v]\n", filter)
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("DocumentInCache:433  Count returned error: %s\n", err.Error())
		return 0, err
	}
	//log.Debugf("DocumentsInCache:436 -- Counted %d documents", count)
	return count, nil
}

func DeleteDocuments(patientID string) error {
	//startTime := time.Now()
	//log.Infof("Deleting Documents:423 -- for Patient: %s", patientID)
	col, _ := storage.GetCollection("documents")
	filter := bson.M{"subject.reference": "Patient/" + patientID}
	log.Debugf("DeleteDocument445 -- bson filter delete: %v\n", filter)
	res, err := col.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Errorf("DeleteDocuments:448 -- for patientID %s failed: %v", filter, err)
		return err
	}
	fmt.Printf("\nDeleted %d documents\n", res.DeletedCount)
	return nil
}
