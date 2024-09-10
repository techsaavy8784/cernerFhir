package model

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	"github.com/dhf0820/cernerFhir/pkg/storage"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Encounter struct {
	_ID             primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	SessionID       string             `json:"-"`
	VisitNum        string             `json:"visit_num" bson:"visit_num"`
	AdmitDate       time.Time          `json:"admit_date" bson:"admit_date"`
	DischargeDate   time.Time          `json:"discharge_date" bson:"discharge_date"`
	Facility        string             `json:"facility"`
	PatientID       string             `json:"patient_id" bson:"patient_id"`
	PatientName     string             `json:"patient_name" bson:"patient_name"`
	MRN             string             `json:"mrn"`
	EncounterID     string             `json:"encounter_id" bson:"encounter_id"`
	AccountNumber   string             `json:"account_number" bson:"account_number"`
	Description     string             `json:"description"`
	Type            string             `json:"type"`
	PatientGPI      string             `json:"patient_gpi" bson:"patient_gpi"`
	VersionID       string             `json:"version_id" bson:"version_id"`
	LastUpdated     time.Time          `json:"last_updated" bson:"last_updated_at"`
	Status          string             `json:"status"`
	Source          string             `json:"source"`
	Text            string             `json:"text"`
	TextStatus      string             `json:"text_status" bson:"text_status"`
	Reason          string             `json:"reason"`
	Class           string             `json:"Class"`
	LocationID      string             `json:"location_id" bson:"location_id"`
	LocationDisplay string             `json:"location_display" bson:"location_display"`
	//encounter         encounter            `bson:"encounter,omitempty`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	AccessedAt time.Time `json:"accessed_at"`
}

type EncounterSummary struct {
	ID            string    `json:"id"`
	VisitNum      string    `json:"visit_num"`
	AdmitDate     time.Time `json:"admit_date"`
	DischargeDate time.Time `json:"discharge_date"`
	PatientID     string    `json:"patientid"`
	PatientName   string    `json:"patient_name"`
	Reason        string    `json:"reason"`
	AccountNumber string    `json:"account_number"`
	Status        string    `json:"status"`
	Class         string    `json:"Class"`
}

type EncounterFilter struct {
	Skip         int64    `schema:"skip"`
	Page         int64    `schema:"page"`
	PageStr      string   `schema:"page_str"`
	Limit        int64    `schema:"limit"`
	SortBy       []string `schema:"column"`
	Order        string   `schema:"order"`
	Mode         string   `schema:"mode"`
	ResultFormat string
	Count        string
	UseCache     string `schema:"useCashe"`
	FillCache    string `schema:"fillCache"`
	Cache        string `schema:"cache"` // reset, stats
	Session      *AuthSession
	SessionId    string `schema:"session_id"`
	UserId       string
	//AccessDetails *AccessDetails

	EncounterID string `schema:"id"`
	PatientID   string `schema:"patientGPI"`
	//MRN           string `schema:"mrn"`
	//EnterpriseID  string `schema:"enterpise_id"`
	// Encounter     string `schema:"encounter"`
	AccountNumber    string `schema:"accountnumber"`
	Class            string `schema:"class"`
	FhirQueryString  string
	CacheQueryFilter bson.M
}

var fhirC *fhir.Connection

func (ef *EncounterFilter) SearchEncounters() ([]*fhir.Encounter, error) {
	//activeEncounterFilter = ef
	println("In Encounter Search\n")
	// spew.Dump(f)
	// println()
	fhirC = config.Fhir()
	// f.fhirC = fhirC
	//activeEncounterFilter = ef
	//fmt.Printf("SearchEncounters:103  ef: %s\n", spew.Sdump(ef))
	ef.makeFilters() // sets up QueryString, queryMap, and queryFilter

	var encounters []*fhir.Encounter
	startTime := time.Now()
	fmt.Printf("Query that should be used to retrieve the requested: %s\n", ef.FhirQueryString)
	encounters, err := ef.FindFhirEncounters()
	fmt.Printf("%d encounters returned in %s\n", len(encounters), time.Since(startTime))
	if err != nil {
		fmt.Printf("f.FindPatient Err: %v\n encounters: %v\n", err, encounters)
		return nil, err

	}
	if len(encounters) == 0 {
		err = fmt.Errorf("no encounters found for: %s", ef.FhirQueryString)
	}
	return encounters, err
}

func (ef *EncounterFilter) FindFhirEncounters() ([]*fhir.Encounter, error) {
	//fmt.Printf("FindFhirEncounter query: %s\n", query)
	config := ActiveConfig()

	c := config.Fhir()

	start := time.Now()
	fmt.Printf("FindFhirEncounters: %s\n", ef.FhirQueryString)
	encResults, err := c.FindFhirEncounters(ef.FhirQueryString)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n\n\n Initial EncounterQuery length : %d\n\n", len(encResults.Entry))
	log.Debug("Start following next links")
	go ef.FollowNextFhirEncLinks(encResults.SearchResult.Link)
	elapsed := time.Since(start)
	fmt.Printf("GetEncounters took %s\n", elapsed)
	if len(encResults.Entry) == 0 {
		return nil, fmt.Errorf("404, err: %v", err)
	}
	encounters := []*fhir.Encounter{}

	for _, e := range encResults.Entry {
		err = InsertFhirEnc(&e.Encounter, ef.Session.EncSessionId)
		if err != nil {
			if storage.IsDup(err) {
				err = nil
			} else {
				return nil, fmt.Errorf("InsertFhirEnc error: %s", err.Error())
			}
		}
		encounters = append(encounters, &e.Encounter)
	}
	return encounters, err
}

func (ef *EncounterFilter) FollowNextFhirEncLinks(links []fhir.Link) {
	fmt.Printf("\n\n\n\n\n####  FollowNextFhirEncLinks:124\n\n\n")
	url := NextPageLink(links)
	i := 2
	ef.Session.UpdateEncStatus("filling")

	//time.Sleep(10 * time.Second)
	for {
		startTime := time.Now()
		if url == "" {

			//ef.Session.UpdateEncStatus("done")

			log.Info("CachePages for url is blank, done")
			break
		}
		links, _ = ef.ProcessFhirEncPage(url)
		fmt.Printf("\n--------Page: %d  added in %f seconds\n\n\n", i, time.Since(startTime).Seconds())
		i = i + 1
		url = NextPageLink(links)
	}

	ef.Session.UpdateEncStatus("done")

}

func (ef *EncounterFilter) ProcessFhirEncPage(url string) ([]fhir.Link, error) {
	//startTime := time.Now()
	fhirEncounterResults, err := fhirC.NextFhirEncounters(url)
	if err != nil {
		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
		return nil, err
	}

	var fhirEncounters []*fhir.Encounter
	for _, entry := range fhirEncounterResults.Entry {
		enc := entry.Encounter
		err = InsertFhirEnc(&enc, ef.Session.EncSessionId)
		if err != nil {
			if storage.IsDup(err) {
				err = nil
			} else {
				return nil, err
			}
		}

		fhirEncounters = append(fhirEncounters, &enc)
	}
	if fhirEncounters == nil {
		return nil, fmt.Errorf("404|no fhirEncounters found for %s", ef.FhirQueryString)
	}
	return fhirEncounterResults.Link, nil
}

func InsertFhirEnc(enc *fhir.Encounter, sessionId string) error {
	//fmt.Printf("add Fhir Encounter: %s\n", spew.Sdump(enc))
	enc.CacheID = primitive.NewObjectID()
	tn := time.Now()
	enc.CreatedAt = &tn
	collection, _ := storage.GetCollection("encounters")
	enc.SessionId = sessionId
	_, err := collection.InsertOne(context.TODO(), enc)
	if err != nil {
		if storage.IsDup(err) {
			log.Warnf("InsertFhirEnc:277 - Encounter-%s exists", *enc.Id)
			err = nil // Not realy a nil, just not inserted
		} else {
			err = fmt.Errorf("InsertFhirEnc:281 - Error: %s", err.Error())
		}
	}
	// } else {
	// 	enc_ID := insertResult.InsertedID.(primitive.ObjectID)
	// 	enc.CacheID = enc_ID
	// 	fmt.Printf("\n\nCacheID: %s for ID: %s\n\n", enc_ID.String(), enc.ID)
	// }
	return err
}

func GetFhirEncounter(id string) (*fhir.Encounter, error) {
	fmt.Printf("GetFhirEncounter:239 -- ID: %s\n", id)
	enc, err := fhirC.GetEncounter(id)
	if err != nil {
		log.Errorf("FhirGetEncounter:242 -- err: %s", err.Error())
		return nil, err
	}
	return enc, nil
}

func GetFhirEncounterByID(id string) (*fhir.Encounter, error) {
	c := config.Fhir()
	enc, err := c.GetEncounter(id)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// func parseEncounterResults(result *fhir.EncounterResults) []*fhir.Encounter {
// 	var encounters []*fhir.Encounter
// 	for _, entry := range result.Entry {
// 		enc := entry.Encounter

// 		encounters = append(encounters, &enc)
// 	}
// 	return encounters
// }

// func parseEncounterResultsSummary(result *fhir.EncounterResults) []*EncounterSummary {
// 	var encounters []*EncounterSummary
// 	for _, entry := range result.Entry {
// 		enc, _ := parseEncounterSummary(entry.Encounter)
// 		encounters = append(encounters, &enc)
// 	}
// 	return encounters
// }

// func parseEncounter(e fhir.Encounter) (*Encounter, error) {
// 	//ids := extractIDs(p.Identifier)
// 	var enc Encounter
// 	// fmt.Printf("\n\n\nParse Encounters: \n")
// 	// spew.Dump()
// 	// fmt.Printf("End of Parse Spew\n")
// 	enc.EncounterID = e.ID
// 	enc.EnterpriseID = e.ID
// 	enc.VisitNum = e.ID

// 	enc.AdmitDate = e.Period.Start
// 	enc.DischargeDate = e.Period.End
// 	if e.Subject.Reference != "" {
// 		enc.PatientID = strings.Split(e.Subject.Reference, "/")[1]
// 		enc.PatientName = e.Subject.Display
// 	} else {
// 		enc.PatientID = strings.Split(e.Patient.Reference, "/")[1]
// 		enc.PatientName = e.Patient.Display
// 	}

// 	enc.AccountNumber = parseAccount(e.Identifiers)
// 	if e.Reasons != nil {
// 		enc.Reason = e.Reasons[0].Text
// 	}
// 	enc.Status = e.Status
// 	enc.Text = e.Text.Div
// 	enc.TextStatus = e.Text.Status
// 	enc.Class = e.Class

// 	// Get the encounter information

// 	// fmt.Printf("Dump Encounters\n")
// 	// spew.Dump(e)
// 	// fmt.Printf("End of spew\n")
// 	err := enc.Insert()
// 	return &enc, err
// }

// func parseEncounterSummary(e fhir.Encounter) (EncounterSummary, error) {
// 	//ids := extractIDs(p.Identifier)
// 	var enc EncounterSummary
// 	enc.VisitNum = e.ID
// 	enc.ID = e.ID
// 	enc.AdmitDate = e.Period.Start
// 	enc.DischargeDate = e.Period.End
// 	if e.Subject.Reference != "" {
// 		enc.PatientID = strings.Split(e.Subject.Reference, "/")[1]
// 		enc.PatientName = e.Subject.Display
// 	} else {
// 		enc.PatientID = strings.Split(e.Patient.Reference, "/")[1]
// 		enc.PatientName = e.Patient.Display
// 	}
// 	enc.AccountNumber = parseAccount(e.Identifiers)
// 	if e.Reasons != nil {
// 		enc.Reason = e.Reasons[0].Text
// 	}
// 	enc.Status = e.Status
// 	enc.Class = e.Class
// 	//enc.Text = e.Text.Div
// 	//enc.TextStatus = e.Text.Status

// 	// fmt.Printf("Dump Encounters\n")
// 	// spew.Dump(e)
// 	// fmt.Printf("End of spew\n")
// 	return enc, nil
// }

func ExtractAccountNum(idents []fhir.Identifier) string {
	var finNbr string
	for _, id := range idents {
		if id.Type.Text == "FIN NBR" {
			finNbr = id.Value
		}
	}
	return finNbr
}

// func GetCachedEncounter(id string) (*fhir.Encounter, error) {

// 	filter := bson.M{"enterpriseid": id}
// 	//filter := bson.D{{"encounterid", id}}
// 	collection, err := storage.GetCollection("encounters")
// 	if err != nil {
// 		fmt.Printf(" Error getting fhir_encounter Collection: %s\n", err)
// 		return nil, err
// 	}
// 	fenc := fhir.Encounter{}
// 	fmt.Printf("\nGetCachedEncounter Filter: %v\n", filter)
// 	err = collection.FindOne(context.TODO(), filter).Decode(&fenc)
// 	if err != nil {
// 		fmt.Printf("error in find encounter: %v\n", err)
// 		return nil, err
// 	}

// 	// fmt.Printf("\n\n============ Found Cached Encounter\n")
// 	// spew.Dump(enc)
// 	// fmt.Printf("================== id: %v\n\n", enc.ID)
// 	//fenc.AccessedAt = fenc.UpdateAccess()

// 	return &fenc, err
// }

// GetFhirEncounterPage: Queries the EncounterCache based upon session and PatientID(patientGPI) and any filters provided
// returns the slice of the max of pageSize matching Encounters, numberInPage, numPages, totalInCache
//Was QueryCache
func (ef *EncounterFilter) GetFhirEncounterPage() ([]*fhir.Encounter, int64, int64, int64, string, error) {
	var linesPerPage int64 = LinesPerPage()
	var skip int64 = 0
	const DESC = -1
	const ASC = 1
	mq := []bson.M{}
	//mq = append(mq, bson.M{"session_id": ef.Session.EncSessionId})
	mq = append(mq, bson.M{"patient.reference": "Patient/" + ef.PatientID})
	filter := bson.M{"$and": mq}
	//fmt.Printf("GetFhirEncounterPage:391 -- EncounterFilter: %s\n", spew.Sdump(ef))

	if ef.Limit > 0 {
		fmt.Printf("GetFhirEncounterPage:394 -- setting linesPerPage: %d\n", ef.Limit)
		linesPerPage = ef.Limit
	}
	if ef.Page > 0 {
		skip = (ef.Page - 1) * linesPerPage
		//fmt.Printf("GetEncounterPage:399 -- setting skip: %d\n", skip)
	}
	if ef.Skip > 0 {
		skip = ef.Skip
	}

	startTime := time.Now()
	findOptions := options.Find()
	findOptions.SetLimit(linesPerPage)
	findOptions.SetSkip(skip)
	//findOptions.SetSort(bson.D{bson.E{"family", 1}, bson.E{"given", 1}})
	//if pf.SortBy == "" || pv.SortBy {}
	sortOrder := ASC // Default Assending
	var sortFields bson.D

	if strings.ToLower(ef.Order) == "desc" {
		sortOrder = DESC
	}
	sortFields = append(sortFields, bson.E{Key: "id", Value: sortOrder})
	findOptions.SetSort(sortFields)
	log.Debugf("GetFhirEncounterPage:439 - sortFields: %s", sortFields)
	collection, _ := storage.GetCollection("encounters")
	cursor, err := collection.Find(context.TODO(), filter, findOptions)

	var encounters []*fhir.Encounter

	if err != nil {
		msg := fmt.Sprintf("GetFhirEncounterPage:447 -  Query for %s returned error: %s\n", filter, err.Error())
		log.Error(msg)
		//cursor is not open
		return nil, 0, 0, 0, "", errors.New(msg)
	}
	log.Printf("GetFhirEncounterPage:450 - took %f seconds\n", time.Since(startTime).Seconds())

	//log.Debugf("GetFhirEncounterPage:456 - remaining %d",cursor.RemainingBatchLength())
	for cursor.Next(context.TODO()) {
		enc := fhir.Encounter{}
		err = cursor.Decode(&enc)
		if err != nil {
			msg := fmt.Sprintf("GetFhirEncounter -- Decode failed: %s", err.Error())
			//fmt.Printf("   Next error: %v\n", err)
			cursor.Close(context.TODO())
			return nil, 0, 0, 0, "", errors.New(msg)
		}

		//fmt.Printf("  Added one to Encounters\n")
		encounters = append(encounters, &enc)
	}
	//fmt.Printf("   Finished fetching documents: error: %v\n\n", cursor.Err())
	cursor.Close(context.TODO())
	numberInPage := int64(len(encounters))
	if numberInPage == 0 {
		msg := fmt.Sprintf("GetFhirEncounterPage:474 - No Encounters found for %s", filter)
		log.Error(msg)
		return nil, 0, 0, 0, "", fmt.Errorf("notFound")
	} else {
		log.Printf("GetFhirEncounterPage:484 - found %d encounters \n", numberInPage)
	}
	// totalInCache, err := ef.FhirEncountersInCache()
	// if err != nil {
	// 	log.Errorf("fhirEncountersInCache error: %s", err.Error())
	// }
	numPages, totalInCache, cacheStatus, err := ef.FhirEncounterPagesInCache()
	return encounters, numberInPage, numPages, totalInCache, cacheStatus, err
}

//EncounterCacheStats:  Queries the cache for the number of Encounters for the session, calculates how many pages that is
//based upon thepage size specified in the environment, returns numPages, totalInCache, error
func (ef *EncounterFilter) FhirEncounterCacheStats() (int64, int64, string, error) {
	// totalInCache, err := ef.FhirEncountersInCache()
	// if err != nil {
	// 	msg := fmt.Sprintf("EncounterCacheStats:496 -- err: %s", err.Error())
	// 	return 0, 0, errors.New(msg)
	// }
	numPages, totalInCache, cacheStatus, err := ef.FhirEncounterPagesInCache()
	if err != nil {
		msg := fmt.Sprintf("EncounterCacheStats:501 -- err: %s", err.Error())
		return 0, 0, "", errors.New(msg)
	}
	return numPages, totalInCache, cacheStatus, nil

}

func (ef *EncounterFilter) FhirEncounterPagesInCache() (int64, int64, string, error) {
	totalInCache, cacheStatus, err := ef.FhirEncountersInCache()
	if err != nil {
		return 0, 0, "", err
	}
	pageSize := LinesPerPage()
	pagesInCache, _ := CalcPages(totalInCache, pageSize)
	//pages := inCache/pageSize
	log.Debugf("EncounterPagesInCache:519 -- pageSize: %d  InCache: %d, pagesInCaches: %d", pageSize, totalInCache, pagesInCache)
	return pagesInCache, totalInCache, cacheStatus, nil
}

//FhirEncountersInCache: returns the number of documents in cache for the session and patient.
func (ef *EncounterFilter) FhirEncountersInCache() (int64, string, error) {
	//mq := append(f.CacheFilterBase, bson.M{"source": f.Source})
	//filter := bson.M{"$and": mq}
	mq := []bson.M{}
	//mq = append(mq, bson.M{"session_id": ef.Session.EncSessionId})
	mq = append(mq, bson.M{"patient.reference": "Patient/" + ef.PatientID})
	filter := bson.M{"$and": mq}

	c, err := storage.GetCollection("encounters")
	if err != nil {
		log.Errorf("Settng Collection(%s) failed: %s", "encounters", err.Error())
	}
	//filter := bson.M{"session_id": ef.Session.EncSessionId}
	//log.Infof("EncountersInCache:530 For session matching: [%v]\n", filter)
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("encountersInCache:535  Count returned error: %v\n", err.Error())
		return 0, "", err
	}
	//log.Debugf("EncountersInCache:538 -   Counted %d Encounters", count)
	cacheStatus := ef.Session.GetEncounterStatus()
	return count, cacheStatus, nil
}

//Get Chached entries first. If none, get them from FHIR
// func getCachedEncounters(filter bson.M) ([]*Encounter, error) {
// 	var encounters []*Encounter
// 	collection, err := storage.GetCollection("encounters")
// 	if err != nil {
// 		fmt.Printf(" Error getting Collection: %s\n", err)
// 		return nil, err
// 	}
// 	//var encounter = new(Encounter)
// 	fmt.Printf("\nFilter: %v\n", filter)
// 	cursor, err := collection.Find(context.TODO(), filter)
// 	if err != nil {
// 		cursor.Close(context.TODO())
// 		return nil, err
// 	}
// 	for cursor.Next(context.TODO()) {
// 		var encounter Encounter
// 		err = cursor.Decode(&encounter)
// 		if err != nil {
// 			cursor.Close(context.TODO())
// 			return nil, err
// 		}

// 		encounter.AccessedAt = encounter.UpdateAccess()
// 		encounters = append(encounters, &encounter)
// 	}
// 	return encounters, err
// }

func (enc *Encounter) UpdateAccess() time.Time {
	filter := bson.M{"_id": enc._ID}
	fmt.Printf("@@ Update finder: %v\n", filter)
	loc, _ := time.LoadLocation("UTC")
	accessed := time.Now().In(loc)
	update := bson.M{"$set": bson.M{"accessedat": accessed}}
	collection, _ := storage.GetCollection("encounters")
	res, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Printf(" Update error ignored: %s\n", err)
	}
	fmt.Printf("Matched: %d  -- modified: %d\n", res.MatchedCount, res.ModifiedCount)
	(*enc).AccessedAt = accessed
	return accessed
}

func (ef *EncounterFilter) makeFilters() {
	ef.makeFhirQueryString()
	ef.makeCacheQueryFilter()
}

func (ef *EncounterFilter) makeCacheQueryFilter() {
	fmt.Printf("makeEncounter:574 - session: %s\n", spew.Sdump(ef.Session))
	mq := []bson.M{}
	sess := bson.M{"id": ef.Session.EncSessionId}
	fmt.Printf("sess : %v\n", sess)
	mq = append(mq, sess)
	if ef.EncounterID != "" {
		mq = append(mq, bson.M{"id": ef.EncounterID})
	} //DTSU-2 does not support others. We may add
	ef.CacheQueryFilter = bson.M{"$and": mq}
}

func (ef *EncounterFilter) makeFhirQueryString() {

	ef.FhirQueryString = ""
	if ef.EncounterID != "" {
		ef.FhirQueryString = fmt.Sprintf("_id=%s", ef.EncounterID)
	} else {
		ef.FhirQueryString = fmt.Sprintf("patient=%s", ef.PatientID)
	}
	if ef.Count == "" {
		ef.FhirQueryString = fmt.Sprintf("%s&_count=%d", ef.FhirQueryString, 20)
	} else {
		ef.FhirQueryString = fmt.Sprintf("%s&_count=%s", ef.FhirQueryString, ef.Count)
	}
}
