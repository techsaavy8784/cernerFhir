package model

// Patient allows filtering by id, mrn, names, birthdate, Encounter, andpossibly ssn
// id, mrn encounter and ssn SHOULD only return one patient and funcs are identified by Getxxx.
// The others can return an array of patients and are the functions are called Findxxx
//

import (
	"context"
	"net/http"

	//"errors"
	"errors"
	"fmt"

	//"net/http"
	//"strconv"
	"strings"

	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/dhf0820/cernerFhir/fhirongo"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	"github.com/dhf0820/cernerFhir/pkg/common"
	com "github.com/dhf0820/cernerFhir/pkg/common"
	"github.com/dhf0820/cernerFhir/pkg/storage"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Patient contains the cached and returned information for one patient
type Patient struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	SessionID     string             `json:"-"`
	MPI           string             `json:"mpi"`
	Facility      string             `json:"facility"`
	MRN           string             `json:"mrn"`
	SSN           string             `json:"ssn"`
	Name          string             `json:"name"`
	Given         string             `json:"given"`
	MiddleName    string             `json:"middle_name"`
	Family        string             `json:"family"`
	Sex           string             `json:"sex"`
	MaritalStatus string             `json:"marital_status"`
	BirthDate     time.Time          `json:"birth_date"`
	DeathDate     time.Time          `json:"death_date"`
	Address1      string             `json:"address_1"`
	Address2      string             `json:"address_2"`
	City          string             `json:"city"`
	State         string             `json:"state"`
	PostalCode    string             `json:"postal_code"`
	Country       string             `json:"country"`
	Email         string             `json:"email"`
	HomePhone     string             `json:"home_phone"`
	CellPhone     string             `json:"cell_phone"`
	WorkPhone     string             `json:"work_phone"`
	EnterpriseID  string             `json:"enterprise_id"`
	Source        string             `json:"source"`
	Text          fhir.TextData      `json:"text" bson:"text"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	AccessedAt    time.Time          `json:"accessed_at"`
	DeletedAt     time.Time          `json:"deleted_at"`
}

type CAPatient struct {
	ID            primitive.ObjectID `json:"-" bson:"_id"`
	SessionID     string             `json:"-" bson:"session_id"`
	MRN           string             `json:"mrn"`
	SSN           string             `json:"ssn"`
	PatientGPI    string             `json:"patient_gpi"`
	Name          string             `json:"name"`
	FirstName     string             `json:"given" bson:"firstname"`
	MiddleName    string             `json:"middle_name" bson:"middlename"`
	LastName      string             `json:"family" bson:"lastname"`
	Text          string             `json:"text"`
	BirthDate     *time.Time         `json:"birth_date"`
	DeathDate     string             `json:"death_date"`
	Sex           string             `json:"sex"`
	MaritalStatus string             `json:"marital_status"`
	Address1      string             `json:"address_1"`
	Address2      string             `json:"address_2"`
	City          string             `json:"city"`
	State         string             `json:"state"`
	PostalCode    string             `json:"postal_code"`
	Country       string             `json:"country"`
	Email         string             `json:"email"`
	HomePhone     string             `json:"home_phone"`
	CellPhone     string             `json:"cell_phone"`
	WorkPhone     string             `json:"work_phone"`
	Identifiers   []*Identifier      `json:"identifiers"`
	Source        string             `json:"source"`
	CreatedAt     *time.Time         `json:"-" bson:"created_at"`
	AccessedAt    *time.Time         `json:"-" bson:"accessed_at"`
}

type Text struct {
	Status string `json:"status" bson:"status"`
	Div    string `json:"div" bson:"div"`
}

type Identifier struct {
	System   string `json:"system"`
	Facility string `json:"facility"`
	Value    string `json:"value"`
}

type PatientFilter struct {
	Skip          int64    `schema:"skip"`
	Page          int64    `schema:"page"`
	PageStr       string   `schema:"page_str"`
	Limit         int64    `schema:"limit"`
	SortBy        []string `schema:"column"`
	Order         string   `schema:"order"`
	Mode          string   `schema:"mode"`
	ResultFormat  string
	Count         string
	UseCache      string `schema:"useCashe"`
	FillCache     string `schema:"fillCache"`
	Cache         string `schema:"cache"` // reset, stats
	Session       *AuthSession
	SessionId     string `schema:"session_id"`
	UserId        string
	JWTokenStr    string
	JWToken       *jwt.Token
	TokenCookie   *http.Cookie
	AccessDetails *AccessDetails

	//ID           string `schema:"id"`  //EnterpriseID shoud be used only
	MRN string `schema:"mrn"`
	SSN string `schema:"ssn"`
	//GivenExact   string `schema:"given_exact"`
	FamilyExact string `schema:"family_exact"`
	Given       string `schema:"given"`
	FirstName   string `schema:"first_name"`
	Family      string `schema:"family"`
	LastName    string `schema:"last_name"`
	BirthDate   string `schema:"birth_date"`
	Email       string `schema:"email"`
	PatientGPI  string `schema:"patient_gpi"`
	EncounterID string `schema:"encounter"`
	Source      string `schema:"source"`

	queryMap    map[string]string
	queryString string
	queryFilter bson.M
	cacheFilter bson.M
}

//fhirC *fhir.Connection

var activePatientFilter *PatientFilter

// func GetFhirPatients(query string) ([]*Patient, error) {
// 	return nil, fmt.Errorf("Not implemented")
// }

/////////////////////////////////////////////////////////////////////////////////
//                           Session Filtering                                 //
/////////////////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////////////////////////////
// SessionQueryCache: Uses the PatientFilter , Creates a new CcheFilter based
// upon the new query values and returns the found in the sort order requested
// and selected fields matched.
// Additionally the total patients in the resulting from the filter, number in
// a page, the calculated number of pages, the current page returned, and error.
/////////////////////////////////////////////////////////////////////////////////

// func (pf *PatientFilter) SessionQueryCache() ([]*CAPatient, int, int, int, int, error) {
// 	err := pf.MakeCacheFilter()
// 	if err != nil {
// 		return nil, 0, 0, 0, 0, err
// 	}

// 	patients, total, inPage, pages, page, err := pf.GetPatientPage()
// 	return patients, total, inPage, pages, page, err

// }

// func (pf *PatientFilter) GetPatientPage(sessionId string) ([]*CAPatient, int, int, int, int, error) {
// 	err := pf.MakeCacheFilter()
// 	if err != nil {
// 		return nil, 0, 0, 0, 0, err
// 	}

// 	return nil, 0, 0, 0, 0, nil
// }

// GetPatientPage: Queries the PatientCache based upon session and any filters provided
// returns the slice of the max of pageSize matching patients, numberInPage, numPages, totalInCache
//Was QueryCache
func (pf *PatientFilter) GetPatientPage() ([]*CAPatient, int, int64, int64, error) {
	var linesPerPage int64 = LinesPerPage()
	var skip int64 = 0
	const DESC = -1
	const ASC = 1

	fmt.Printf("GetPatientPage:203 -- PatientFilter: %s\n", spew.Sdump(pf))

	if pf.Limit > 0 {
		fmt.Printf("GetPatientPage:206 -- setting linesPerPage: %d\n", pf.Limit)
		linesPerPage = pf.Limit
	}
	if pf.Page > 0 {
		skip = (pf.Page - 1) * linesPerPage
		fmt.Printf("GetPatientPage:211 -- setting skip: %d\n", skip)
	}
	if pf.Skip > 0 {
		skip = pf.Skip
	}

	startTime := time.Now()
	//pf.cacheFilter["session_id"] = pf.Session.SessionID
	log.Debugf("GetPatientPage:209 - Paging limit: %d   skip:%d", linesPerPage, skip)
	log.Debugf("GetPatientPage:210 - Using cacheFilter %s ", pf.cacheFilter)
	log.Debugf("GetPatientPage:210 - Using queryFilter %s ", pf.queryFilter)
	filter := bson.M{"session_id": pf.Session.PatSessionId}
	//log.Debugf("GetPatientPage:211 - cacheFilter: %v", pf.cacheFilter)
	findOptions := options.Find()
	findOptions.SetLimit(linesPerPage)
	findOptions.SetSkip(skip)
	//findOptions.SetSort(bson.D{bson.E{"family", 1}, bson.E{"given", 1}})
	//if pf.SortBy == "" || pv.SortBy {}
	sortOrder := ASC // Default Assending
	var sortFields bson.D

	if strings.ToLower(pf.Order) == "desc" {
		sortOrder = DESC
	}

	// if len(pf.SortBy) > 0 {
	// 	if len(pf.SortBy) > 0 {
	// 		fmt.Printf("GetPatientPage:226 -- SortBy: %v\n", pf.SortBy)
	// 		for _, s := range pf.SortBy {
	// 			fmt.Printf("   Adding [%s]\n", s)

	// 			sortFields = append(sortFields, bson.E{s, sortOrder})
	// 		}
	// 	} else {

	log.Infof("GetPatientPage:241 -- Setting sort fields to Last_name, first_name")
	sortFields = append(sortFields, bson.E{"lastname", sortOrder})
	sortFields = append(sortFields, bson.E{"firstname", sortOrder})
	//TODO: Add sort by birthdate
	// 	}
	// }
	findOptions.SetSort(sortFields)
	log.Debugf("GetPatientPage:251 - sortFields: %s", sortFields)
	collection, _ := storage.GetCollection("capatients")
	log.Debugf("GetPatientPage:253 - queryFilter: %v", pf.queryFilter)
	cursor, err := collection.Find(context.TODO(), filter, findOptions)

	var patients []*CAPatient

	if err != nil {
		msg := fmt.Sprintf("GetPatientPage:259 -  Query for %s returned error: %v\n", filter, err)
		log.Error(msg)
		//cursor is not open
		return nil, 0, 0, 0, errors.New(msg)
	}
	log.Printf("GetPatientPage:264 - took %f seconds\n", time.Since(startTime).Seconds())
	//defer cursor.Close(context.TODO)  // Need to get real context
	//fmt.Printf("\n    No Error on QueryCache\n\n")

	//log.Debugf("GetPatientPage:244 - remaining %d",cursor.RemainingBatchLength())
	for cursor.Next(context.TODO()) {
		var patient CAPatient
		err = cursor.Decode(&patient)
		if err != nil {
			msg := fmt.Sprintf("GetPatientPate:273 -- Decode failed: %s", err.Error())
			//fmt.Printf("   Next error: %v\n", err)
			cursor.Close(context.TODO())
			return nil, 0, 0, 0, errors.New(msg)
		}

		fmt.Printf("  Added one to documents\n")
		patients = append(patients, &patient)
	}
	//fmt.Printf("   Finished fetching documents: error: %v\n\n", cursor.Err())
	cursor.Close(context.TODO())
	numberInPage := len(patients)
	if numberInPage == 0 {
		msg := fmt.Sprintf("GetPatientPage:288 - No Documents found for %s", filter)
		log.Error(msg)
		return nil, 0, 0, 0, fmt.Errorf(msg)
	} else {
		log.Printf("GetPatientPage:291 - found %d documents \n", numberInPage)
	}
	totalInCache, err := pf.PatientsInCache()
	numPages, err := pf.PatientPagesInCache()
	return patients, numberInPage, numPages, totalInCache, err
}

//PatientCacheStats:  Queries the cache for the number ther for the session, calculates how many pages that is
//based upon thepage size specified in the environment, returns numPages, totalInCache, error
func (pf *PatientFilter) PatientCacheStats() (int64, int64, error) {
	totalInCache, err := pf.PatientsInCache()
	if err != nil {
		msg := fmt.Sprintf("PatientCacheStats:301 -- err: %s", err.Error())
		return 0, 0, errors.New(msg)
	}
	numPages, err := pf.PatientPagesInCache()
	if err != nil {
		msg := fmt.Sprintf("PatientCacheStats:306 -- err: %s", err.Error())
		return 0, 0, errors.New(msg)
	}
	return numPages, totalInCache, nil

}

func (pf *PatientFilter) QueryCache() ([]*fhir.Patient, error) {
	var linesPerPage int64 = LinesPerPage()
	var skip int64 = 0
	const DESC = -1
	const ASC = 1

	fmt.Printf("QueryCache319 -- QueryCache filter: %s\n", spew.Sdump(pf))

	if pf.Limit > 0 {
		fmt.Printf("QueryCache:322 -- setting linesPerPage: %d\n", pf.Limit)
		linesPerPage = pf.Limit
	}
	if pf.Page > 0 {
		skip = (pf.Page - 1) * linesPerPage
		fmt.Printf("QueryCache:327 -- setting skip: %d\n", skip)
	}
	if pf.Skip > 0 {
		skip = pf.Skip
	}
	// q, err := pf.QueryCacheByEncounter()

	fmt.Printf("QueryCache:334 -- Paging perPage: %d   skip:%d\n", linesPerPage, skip)
	fmt.Printf("QueryCache:335 -- Using filter %s \n", pf.queryFilter)
	startTime := time.Now()
	// spew.Dump(q)
	// spew.Dump(f.QueryFilter)
	findOptions := options.Find()
	findOptions.SetLimit(linesPerPage)
	findOptions.SetSkip(skip)
	//findOptions.SetSort(bson.D{bson.E{"family", 1}, bson.E{"given", 1}})
	//if pf.SortBy == "" || pv.SortBy {}
	sortOrder := ASC // Default Assending
	var sortFields bson.D

	if strings.ToLower(pf.Order) == "desc" {
		sortOrder = DESC
	}

	if len(pf.SortBy) > 0 {
		fmt.Printf("QueryCache:352 -- SortBy: %v\n", pf.SortBy)
		for _, s := range pf.SortBy {
			fmt.Printf("   Adding [%s]\n", s)

			sortFields = append(sortFields, bson.E{s, sortOrder})
		}
	} else {
		log.Infof("QueryCache:359 -- Setting sort fields to Last_name, first_name")
		sortFields = append(sortFields, bson.E{"last_name", sortOrder})
		sortFields = append(sortFields, bson.E{"first_name", sortOrder})
	}
	findOptions.SetSort(sortFields)
	fmt.Printf("QueryCache:364 -- sortFields- %s\n", sortFields)
	collection, _ := storage.GetCollection("capatients")

	cursor, err := collection.Find(context.TODO(), pf.queryFilter, findOptions)

	var patients []*fhir.Patient

	if err != nil {
		log.Printf("QueryCache:372 -- QueryCache for %s returned error: %v\n", pf.queryFilter, err)
		//cursor is not open
		return nil, err
	}
	fmt.Printf("QueryCache:376 -- QueryCache took %f seconds\n", time.Since(startTime).Seconds())
	//defer cursor.Close(context.TODO)  // Need to get real context
	//fmt.Printf("\n    No Error on QueryCache\n\n")
	for cursor.Next(context.TODO()) {
		var patient fhir.Patient
		err = cursor.Decode(&patient)
		if err != nil {
			//fmt.Printf("   Next error: %v\n", err)
			cursor.Close(context.TODO())
			return nil, err
		}

		//fmt.Printf("  Added one to documents\n")
		patients = append(patients, &patient)
	}
	//fmt.Printf("   Finished fetching documents: error: %v\n\n", cursor.Err())
	cursor.Close(context.TODO())
	if len(patients) == 0 {
		err = fmt.Errorf("404|QueryCache:394 -- no documents found for %s", pf.queryFilter)
	} else {
		fmt.Printf("QueryCache:396 -- QueryCache found %d documents \n", len(patients))
	}
	return patients, err
}

func (pf *PatientFilter) MakeCacheFilter() error {
	err := pf.MakeQueryMap()
	if err != nil {
		return err
	}
	layout := "2006-01-02"
	mq := []bson.M{}
	mq = append(mq, bson.M{"session_id": pf.Session.PatSessionId})
	for k := range pf.queryMap {
		val := pf.queryMap[k]
		//fmt.Printf("k: %s,  v: %s\n", k, q[k])
		if k == "patient_gpi" {
			//fmt.Printf("Converting search for id %s to search for enterpriseid\n", val)
			mq = append(mq, bson.M{"id": val})
		} else if k == "given" {
			mq = append(mq, bson.M{"given": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "family" {
			mq = append(mq, bson.M{"family": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "given:exact" {
			mq = append(mq, bson.M{"given": primitive.Regex{strings.Replace("given", "\"", "", -1), "i"}})
		} else if k == "family:given" {
			mq = append(mq, bson.M{"family": primitive.Regex{strings.Replace("family", "\"", "", -1), "i"}})
		} else if k == "email" {
			mq = append(mq, bson.M{"email": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "birthdate" {
			condition := "eq"
			input := ""
			s := strings.Split(val, "|")
			if len(s) > 1 {
				condition = s[0]
				input = s[1]
			} else {
				condition = "$eq"
				input = s[0]
			}
			useDate, err := time.Parse(layout, input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			q := bson.M{"birthdate": bson.M{condition: useDate}}
			mq = append(mq, q)

		} else {
			mq = append(mq, bson.M{k: val})
		}
	}

	if len(mq) > 0 {
		pf.cacheFilter = bson.M{"$and": mq}
	}
	fmt.Printf("\n\nMakeCacheFilter:453  cacheFilter: %v\n\n", pf.cacheFilter)
	fmt.Printf("\n\nMakeCacheFilter:453  queryFilter: %v\n\n", pf.queryFilter)

	return nil
}

func (f *PatientFilter) MakeQueryMap() error {
	fmt.Printf("MakeQueryMap:434 -- Starting MakeQueryMap: %s\n", f.queryMap)
	m := make(map[string]string)
	session := f.Session.PatSessionId
	mrn := strings.Trim(f.MRN, " ")
	given := strings.Trim(f.Given, " ")
	family := strings.Trim(f.Family, " ")
	family_exact := strings.Trim(f.FamilyExact, " ")
	//given_exact := strings.Trim(f.GivenExact, " ")
	encounter := strings.Trim(f.EncounterID, " ")
	patientGPI := strings.Trim(f.PatientGPI, " ")

	email := strings.Trim(f.Email, " ")
	//id := strings.Trim(f.ID, " ")
	birthdate := strings.Trim(f.BirthDate, " ")

	fmt.Printf("MakeQueryMap:474 -- SessionId: %s\n", f.Session.SessionID)
	m["session"] = session
	if family != "" {
		fmt.Printf("MakeQueryMap:477 - Family is set\n")
		m["family"] = family
		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
		f.SortBy = append(f.SortBy, "given")
	}
	if family_exact != "" {

		fam := strings.Split(family_exact, ":")
		fmt.Printf("MakeQueryMap:85 - Family_exact \n")
		m["family"] = fam[0]
		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
		f.SortBy = append(f.SortBy, "given")
	}
	if given != "" {
		if family != "" || family_exact != "" {
			m["given"] = given
		} else {
			log.Warn("MakeQueryMap:494 - Invalid search: given alone is invalid")
			//return fmt.Errorf("400|MakeQueryMap:442 - Invalid search: given alone is invalid")
		}
	}

	if mrn != "" {
		m["mrn"] = mrn
	}
	if encounter != "" {
		m["encounter"] = encounter
	}
	if patientGPI != "" {
		m["patient_gpi"] = patientGPI
		//m["id"] = enterpriseID
	}

	if email != "" {
		m["email"] = email
	}
	if birthdate != "" {
		m["birthdate"] = birthdate
	}
	f.queryMap = m
	fmt.Printf("MakeQueryMap:517 -- created map: %v\n", m)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//                           Convert to CA
////////////////////////////////////////////////////////////////////////////////
//

func (p *CAPatient) InsertCaCache(pf *PatientFilter) error {
	//fmt.Printf("adding: %T: %v\n\n", p, p)
	//fmt.Printf("\n\n\n\n#######  add patient SessionID: %s  Name: %s to cache \n", p.SessionID, p.Name)
	//
	//p.SessionID = p.SessionID

	// _, err := FindByPhone(c.FaxNumber, c.Facility
	// log.Fatal(err)
	//startTime := time.Now()
	collection, err := storage.GetCollection("capatients")
	if err != nil {
		log.Errorf("InsertCaCache Get capatients failed: %s", err.Error())
	}
	id := primitive.NewObjectID()
	p.ID = id
	tn := time.Now().UTC()
	p.AccessedAt = &tn
	p.SessionID = pf.Session.PatSessionId
	insertResult, err := collection.InsertOne(context.TODO(), p)
	if err != nil {
		log.Debug("InsertCaCache:542 err")
		if !storage.IsDup(err) {
			fmt.Printf("InsertCaCache:544 - Patient %s already in Cache\n", p.PatientGPI)
			err = nil
		}
		//fmt.Printf(" Err: %s\n", err.Error())
		//cnt, err1 := pf.CountCachedCaPatients()
		// if err1 != nil {
		// 	log.Errorf("CountCaPatients err: %s", err1.Error())
		// }
		//log.Debugf("Found %d caPatients", cnt)
		//err = nil
	} else {
		//id := insertResult.InsertedID.(primitive.ObjectID)
		//log.Debugf("   Setting Patient ID : %s", id.Hex())
		p.ID = insertResult.InsertedID.(primitive.ObjectID)
		//fmt.Printf("   Set IDid: %s\n", p.ID.String())
		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
	}
	//log.Infof("Elapsed time to add to cache: %f seconds\\nn", time.Since(startTime).Seconds())
	return err
}

// PatientCacheStats(): Calculates the number of numInCache, and pagesInCache
// func (pf *PatientFilter) PatientCacheStats() (int64, int64) {
// 	numInCache, err := pf.PatientsInCache()
// 	if err != nil {
// 		return 0, 0
// 	}
// 	pageSize := LinesPerPage()
// 	pagesInCache, _ := CalcPages(numInCache, pageSize)
// 	return pagesInCache, numInCache
// }

func (pf *PatientFilter) PatientPagesInCache() (int64, error) {
	numInCache, err := pf.PatientsInCache()
	if err != nil {
		return 0, err
	}
	pageSize := LinesPerPage()
	pagesInCache, _ := common.CalcPages(numInCache, pageSize)
	//pages := inCache/pageSize
	log.Debugf("PatientPagesInCache:584 -- pageSize: %d  InCache: %d, pagesInCaches: %d", pageSize, numInCache, pagesInCache)
	return int64(pagesInCache), nil
}

func (f *PatientFilter) CountCachedCaPatients() (int64, error) {
	return f.PatientsInCache()
}

func (f *PatientFilter) PatientsInCache() (int64, error) {
	//mq := append(f.CacheFilterBase, bson.M{"source": f.Source})
	//filter := bson.M{"$and": mq}

	c, err := storage.GetCollection("capatients")
	if err != nil {
		log.Errorf("Settng Collection(%s) failed: %s", "capatient", err.Error())
	}
	filter := bson.M{"session_id": f.Session.PatSessionId}
	log.Infof("PatientsInCache:606 For session matching: [%v]\n", filter)
	count, err := c.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Errorf("patientsInCache:609  Count returned error: %v\n", err.Error())
		return 0, err
	}
	log.Debugf("PatientsInCache:612 -   Counted %d ca_patients", count)
	return count, nil
}

func FhirResultsToCAPatients(results *fhir.PatientResult, pf *PatientFilter) ([]*CAPatient, error) {
	caPatients := []*CAPatient{}
	entries := results.Entry
	for _, entry := range entries {
		pat, err := FhirPatientToCA(&entry.Patient, pf)
		if err != nil {
			return nil, err
		}
		caPatients = append(caPatients, pat)
	}
	return caPatients, nil
}

func FhirPatientsToCA(fps []*fhir.Patient, pf *PatientFilter) ([]*CAPatient, error) {
	caPatients := []*CAPatient{}
	for _, p := range fps {
		pat, err := FhirPatientToCA(p, pf)
		if err != nil {
			return nil, err
		}
		caPatients = append(caPatients, pat)
	}
	return caPatients, nil
}

func FhirPatientToCA(fp *fhir.Patient, pf *PatientFilter) (*CAPatient, error) {
	fmt.Println("FhirPatientToCa:638 -  entered")
	ids := extractIDs(fp.Identifier)
	var caPat CAPatient
	caPat.SessionID = pf.Session.SessionID
	caPat.LastName = extractName(fp, "official", "family")
	caPat.FirstName = extractName(fp, "official", "given")
	//caPat.Name = extractName(*fp, "official", "text")
	caPat.Name = fmt.Sprintf("%s, %s", caPat.LastName, caPat.FirstName)
	caPat.Sex = fp.Gender
	ident := Identifier{}
	ident.System = "cerner" // TODO: Get ID.System from config
	ident.Value = fp.ID
	ident.Facility = ""
	caPat.Identifiers = append(caPat.Identifiers, &ident)

	caPat.PatientGPI = fp.ID
	caPat.MRN = ids["MRN"]
	caPat.SSN = ids["SSN"]
	caPat.MaritalStatus = fp.MaritalStatus.Text
	caPat.Source = strings.ToLower(config.Source())
	caPat.Text = fp.Text.Div

	caPat.BirthDate = extractCaBirthDate(fp)

	address := extractAddress(fp, "home")
	if address != nil {
		switch len(address.Line) {
		case 1:
			caPat.Address1 = address.Line[0]
		case 2:
			caPat.Address1 = address.Line[0]
			caPat.Address2 = address.Line[1]
		}
		caPat.City = address.City
		caPat.State = address.State
		caPat.PostalCode = address.PostalCode
		caPat.Country = address.Country
	}

	caPat.Email = extractTelecom(fp, "email", "home")
	caPat.HomePhone = extractTelecom(fp, "phone", "home")
	caPat.CellPhone = extractTelecom(fp, "phone", "cell")
	caPat.WorkPhone = extractTelecom(fp, "phone", "work")

	//fmt.Printf("\n-------FhirPatientToCA:594 - calling caPat.InsertCaCache\n")
	err := caPat.InsertCaCache(pf)

	if err != nil {
		if !storage.IsDup(err) {
			//fmt.Printf("FhirPatientToCA:605 - Patient: %s - Already Exisits\n", caPat.PatientGPI)
			err = nil
		} else {
			return nil, err
		}

		return nil, err

	}
	fmt.Printf("------FhirPatientToCa:601 -  InsertCaCache Returning  ID: %s  SessionId: %s\n", caPat.ID.Hex(), caPat.SessionID)
	return &caPat, nil
}

func (f *PatientFilter) FollowNextLinks(links []fhirongo.Link) {
	url := NextPageLink(links)
	i := 1
	f.Session.Status.Patient = "filling"
	f.Session.UpdatePatStatus("filling")
	//time.Sleep(10 * time.Second)
	for {
		startTime := time.Now()
		if url == "" {
			f.Session.UpdatePatStatus("done")
			f.Session.Status.Patient = "done"
			log.Info("CachePages for url is blank, done")
			break
		}
		links, _ = f.ProcessCaPatPage(url)
		fmt.Printf("\n--------Page: %d  added in %f seconds\n\n\n", i, time.Since(startTime).Seconds())
		i = i + 1
		url = NextPageLink(links)
	}

}

func (f *PatientFilter) ProcessCaPatPage(url string) ([]fhirongo.Link, error) {
	//startTime := time.Now()
	fhirPatients, err := fhirC.NextFhirPatients(url)
	if err != nil {
		log.Errorf("NextFHIRPatients returned err: %s\n", err.Error())
		return nil, err
	}
	//fmt.Printf("\n\nNext Page of results: %s\n\n", spew.Sdump(fhirPatients))
	err = InsertFhirPatResults(fhirPatients, f.Session.PatSessionId)
	if err != nil {
		return nil, err
	}
	//TODO: ARE CA DOCUMENTS CACHED
	_, err = FhirResultsToCAPatients(fhirPatients, f) //Documents are cached
	if err != nil {
		return nil, err
	}
	//patients := parsePatientResults(fhirPatients, f.Session.PatSessionID)
	//spew.Dump()

	if fhirPatients == nil {
		return nil, fmt.Errorf("404|no more patients found for %s", f.queryString)
	}
	return fhirPatients.Link, nil
}

// func (p *CAPatient) Insert(sessionId string) error {
// 	//fmt.Printf("adding: %T: %v\n\n", p, p)
// 	fmt.Printf("\n\n\n\n#######  add patient SessionID: %s  Name: %s to cache \n", p.SessionID, p.Name)
// 	//
// 	//p.SessionID = p.SessionID

// 	// _, err := FindByPhone(c.FaxNumber, c.Facility)
// 	// log.Fatal(err)
// 	collection, _ := storage.GetCollection("capatients")

// 	insertResult, err := collection.InsertOne(context.TODO(), p)
// 	if err != nil {
// 		fmt.Printf(" id: %s already exists\n", p.PatientGPI)
// 		err = nil
// 	} else {
// 		p.ID = insertResult.InsertedID.(primitive.ObjectID)
// 		// fmt.Printf(" id: %s\n", p.ID.String())
// 		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
// 	}

// 	return err
// }

//ForMRN returns the patient for the MRN. This is not cached since only one
func (pf *PatientFilter) ForMRN() ([]*fhir.Patient, error) {
	//var patient = new(Patient)
	startTime := time.Now()
	// if pf.UseCache == "true" {
	// //cacheName := pf.SessionId
	// 	filter := bson.M{"mrn": pf.MRN}

	// 	fmt.Printf("ForMRN filter: %v\n", filter)
	// 	patients, err := GetCachedPatients(filter)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Printf("ForMRN Cache FOUND took %s\n", time.Since(startTime))
	// 	return patients, err

	// 	fmt.Printf("\n!!!ForMRN Cache was not found use fhir: %s\n\n", err)
	// }

	mrnQuery := pf.createMRNQuery()

	//config := ActiveConfig()
	c := config.Fhir()
	fmt.Printf("ForMRN calling FHIR using query: %s\n", mrnQuery)
	fhirPatient, err := c.FindFhirPatient(mrnQuery)
	if err != nil {
		fmt.Printf("   FHIR FindFhirPatient returned err: %v\n", err)
		return nil, err
	}
	// fmt.Printf("\nFhir Patient\n")
	// spew.Dump(fhirPatient)
	//patients := parsePatientResults(fhirPatient, pf.Session.PatSessionID)
	// fmt.Printf("\nParsed Patient\n")
	// spew.Dump(patients)

	// (*p) = (*patients[0])
	pats := []*fhir.Patient{}
	for _, entry := range fhirPatient.Entry {
		pat := entry.Patient
		InsertFhirPat(&pat, pf.Session.PatSessionId)
		pats = append(pats, &pat)
	}
	fmt.Printf("ForMRN Fhir FOUND took %s\n", time.Since(startTime))
	return pats, nil
}

func ForPatientGPI(id string) (*fhir.Patient, error) {
	//startTime := time.Now()

	c := config.Fhir()
	fhirPatient, err := c.GetPatient(id)
	if err != nil {
		msg := fmt.Sprintf("GetPatient returned: %s", err.Error())
		return nil, errors.New(msg)
	}
	//fmt.Printf("\nFhir Patient\n")
	//spew.Dump(fhirPatient)
	return fhirPatient, nil
}

// func (pf *PatientFilter) ForID() (*Patient, error) {
// 	//var patient = new(Patient)
// 	fmt.Printf("ForID pf: %v\n", pf)
// 	//fmt.Printf("In ForID using Cache: %s\n", pf.UseCache)
// 	startTime := time.Now()
// 	if pf.UseCache == "true" {
// 		filter := bson.M{"id": pf.EnterpriseID}
// 		fmt.Printf("GetPatientByID filter: %v\n", filter)
// 		patient, err := GetCachedPatient(filter)
// 		if err == nil {
// 			fmt.Printf("GetPatientByID Cache FOUND took %s\n", time.Since(startTime))
// 			return patient, err
// 		}
// 		fmt.Printf("\n!!!GetPatientByID Cache was not found %s\n\n", err)
// 	}

// 	c := config.Fhir()
// 	fhirPatient, err := c.GetPatient(EnterpriseID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	//fmt.Printf("\nFhir Patient\n")
// 	//spew.Dump(fhirPatient)
// 	patient := parsePatient(*fhirPatient, pf.Session.PatSessionID)
// 	fmt.Printf("\nParsed Patient\n")
// 	spew.Dump(patient)
// 	fmt.Printf("\nByID Fhir FOUND took %s\n", time.Since(startTime))
// 	return patient, nil
// }

func (p *Patient) ByEncounter(encID string) error {
	//fmt.Printf("FindFhirPatient query: %s\n", query)
	//config := ActiveConfig()

	//c := config.Fhir()

	// startTime := time.Now()
	// fmt.Printf("GetpatientByEncounter: %s \n", encID)
	// enc := &Encounter{EncounterID: encID}

	// err := enc.ForEncounterID()

	// if err != nil {
	// 	return err
	// }
	// fmt.Printf("Get encounter took %s ms\n", time.Since(startTime))
	// patient := new(Patient)
	// patient.EnterpriseID = enc.PatientID
	// err = patient.ForID(true)
	// fmt.Printf("Get Combined took %s ms\n", time.Since(startTime))
	// if err != nil {
	// 	return err
	// }
	// p = patient
	return nil
}

func extractIDs(ids []fhir.Identifier) map[string]string {
	config := ActiveConfig()
	mrnId := config.MrnID()
	var idents = make(map[string]string)
	//fmt.Printf("\n\nIdents:\n")
	for _, id := range ids {
		//spew.Dump(id)
		if id.System == mrnId {
			id.Type.Text = "MRN"
			idents[id.Type.Text] = id.Value
		}
		if id.System == "http://hl7.org/fhir/sid/us-ssn" {
			//fmt.Printf("!!! SocialSecurity Number: %s:  %s\n", id.Type.Text, id.Value)
			id.Type.Text = "SSN"
			idents[id.Type.Text] = id.Value
		}
		//fmt.Printf("   System: %v  Type: %s - %s\n\n", id.System, id.Type.Text, id.Value)
	}
	return idents
}

// Process each patient in the result bundle. Cache the patient and add to the slice
// of patients to be returned
func PatientsFromResults(result *fhir.PatientResult, cacheName string) []*fhir.Patient {
	pats := []*fhir.Patient{}
	for _, entry := range result.Entry {
		pat := entry.Patient
		pats = append(pats, &pat)
		InsertFhirPat(&pat, cacheName)
	}
	return pats
}

// func parsePatient(p fhir.Patient, cacheName string) *Patient {
// 	//fmt.Prinln("   parsePatient entered")
// 	ids := extractIDs(p.Identifier)
// 	var pat Patient
// 	//fmt.Printf("IDS: %v\n", ids)
// 	// fmt.Printf("\n\nPatient\n")
// 	// spew.Dump(p)

// 	pat.Family = strings.Join(p.Name[0].Family, " ")
// 	pat.Given = strings.Join(p.Name[0].Given, " ")
// 	pat.Name = fmt.Sprintf("%s, %s", pat.Family, pat.Given)
// 	pat.Sex = p.Gender
// 	pat.EnterpriseID = p.ID
// 	pat.MRN = ids["MRN"]
// 	pat.SSN = ids["SSN"]
// 	pat.MaritalStatus = p.MaritalStatus.Text
// 	pat.Source = config.Source()

// 	pat.Text.Div = p.Text.Div
// 	pat.Text.Status = p.Text.Status
// 	extractDates(p, &pat)
// 	extractAddress(p, &pat)
// 	extractTelecom(p, &pat)
// 	//fmt.Printf("BirthDate: %s\n", p.BirthDate)
// 	err := pat.Insert(cacheName)
// 	if err != nil {
// 		fmt.Printf("Patient Insert failed: %s\n", err.Error())
// 	}
// 	return &pat
// }

// func parsePatientResults(result *fhir.PatientResult, cacheName string) []*Patient {
// 	//fmt.Println("parsePatientResults entered")
// 	var patients []*Patient
// 	//fmt.Printf("ParsePatient count: %v\n", len(result.Entry))
// 	for _, entry := range result.Entry {
// 		//fmt.Printf("Working on entry %d\n", i)
// 		//spew.Dump(entry.Patient)

// 		patients = append(patients, parsePatient(entry.Patient, cacheName))
// 	}
// 	return patients
// }

// func parseEncounterResults(result *fhir.PatientEncounterResult) []Encounter {
// 	var encounters []Encounter
// 	//fmt.Printf("ParsePatient count: %v\n", len(result.Entry))
// 	for _, entry := range result.Entry {
// 		//fmt.Printf("Working on entry %d\n", i)
// 		//spew.Dump(entry.Patient)

// 		encounters = append(encounters, parsePatientEncounters(entry.Encounter))
// 	}

// 	return encounters
// }

// ////////////////////////////////////////////////////////////////////////////////
// //                           Convert to CA
// ////////////////////////////////////////////////////////////////////////////////
// func (fps []*fhir.Patient) ConvertPatientsToCA() []*CAPatient {
// 	CApats := []*CAPatient{}
// 	for _, fp := range fps {
// 		p := fp.ConvertToCA()
// 		CApats = append(CApats, p)
// 	}
// 	return CApats
// }

// func (fp *fhir.Patient) ConvertToCA() *CAPatient {
// 	//fmt.Prinln("   parsePatient entered")
// 	ids := extractIDs(p.Identifier)
// 	var pat CAPatient
// 	//fmt.Printf("IDS: %v\n", ids)
// 	// fmt.Printf("\n\nPatient\n")
// 	// spew.Dump(p)

// 	pat.LastName = strings.Join(p.Name[0].Family, " ")
// 	names:= strings.Split(p.Name[0].Given, " ")
// 	pat.FirstName = names[0]
// 	if len(names > 1){
// 		pat.MiddleName = names[1]
// 	}
// 	//pat.FirstName= strings.Join(p.Name[0].Given, " ")
// 	if pat.MiddleName == "" {
// 		pat.Name = fmt.Sprintf("%s, %s", pat.LastName, pat.FirstName)
// 	} else {
// 		pat.Name = fmt.Sprintf("%s, %s %s", pat.LastName, pat.FirstName, pat.MiddleName)
// 	}
// 	pat.Sex = p.Gender
// 	pat.EnterpriseID = p.ID
// 	pat.MRN = ids["MRN"]
// 	pat.SSN = ids["SSN"]
// 	pat.MaritalStatus = p.MaritalStatus.Text
// 	pat.Source = config.Source()
// 	pat.Text.Div = p.Text.Div
// 	pat.Text.Status = p.Text.Status
// 	extractDates(p, &pat)
// 	extractAddress(p, &pat)
// 	extractTelecom(p, &pat)
// 	//fmt.Printf("BirthDate: %s\n", p.BirthDate)
// 	err := pat.Insert(cacheName)
// 	if err != nil {
// 		fmt.Printf("Patient Insert failed: %s\n", err.Error())
// 	}
// 	return &pat
// }

//

func parsePatientEncounter(e fhir.Encounter) Encounter {
	//ids := extractIDs(p.Identifier)
	var enc Encounter
	//fmt.Printf("Encounters: \n")
	//spew.Dump()

	// pat.LastName = strings.Join(p.Name[0].Family, " ")
	// pat.FirstName = strings.Join(p.Name[0].Given, " ")
	// pat.Name = fmt.Sprintf("%s, %s", pat.LastName, pat.FirstName)
	// pat.Sex = p.Gender
	// pat.EnterpriseID = p.ID
	// pat.MRN = ids["MRN"]
	// pat.MaritalStatus = p.MaritalStatus.Text
	// pat.Source = config.Source()
	// extractDates(p, &pat)
	// extractAddress(p, &pat)
	// extractTelecom(p, &pat)
	//fmt.Printf("BirthDate: %s\n", p.BirthDate)
	return enc
}

func (pf *PatientFilter) createMRNQuery() string {
	mrnId := config.MrnID()
	mrnQuery := fmt.Sprintf("identifier=%s|%s", mrnId, pf.MRN)
	//fmt.Printf("Querying for mrn: %s\n", mrnQuery)
	return mrnQuery
}

// //Get Chached entries first. If none, get them from FHIR
// func GetCachedPatients(filter bson.M) ([]*Patient, error) {
// 	var patients []*Patient
// 	collection, err := storage.GetCollection("patients")
// 	if err != nil {
// 		fmt.Printf(" Error getting Collection: %s\n", err)
// 		return nil, err
// 	}
// 	//var encounter = new(Encounter)
// 	fmt.Printf("\nGetCachedPatients Filter: %v\n", filter)
// 	cursor, err := collection.Find(context.TODO(), filter)
// 	if err != nil {
// 		log.Printf("GetCachedPatients search for %s returned error: %v\n", filter, err)
// 		cursor.Close(context.TODO())
// 		return nil, err
// 	}
// 	for cursor.Next(context.TODO()) {
// 		var patient Patient
// 		err = cursor.Decode(&patient)
// 		if err != nil {
// 			cursor.Close(context.TODO())
// 			return nil, err
// 		}
// 		//spew.Dump(patients)
// 		patients = append(patients, &patient)
// 	}
// 	if patients == nil {
// 		err = fmt.Errorf("404|no patients found for %s\n", filter)
// 	}
// 	return patients, err
// }

// // //Get Chached entries first. If none, get them from FHIR
// func GetCachedPatient(filter bson.M) (*Patient, error) {
// 	var patient Patient
// 	collection, err := storage.GetCollection("patients")
// 	if err != nil {

// 		return nil, err
// 	}
// 	fmt.Printf("\nGetCachedPatient Checking for Patient with query: %v\n", filter)
// 	err = collection.FindOne(context.TODO(), filter).Decode(&patient)
// 	//fmt.Printf("FindOne returned err: %v\n", err)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &patient, nil
// }

// func (pf *PatientFilter) QueryCache() ([]*fhir.Patient, error) {
// 	var limit int64 = 20
// 	var skip int64 = 0
// 	const DESC = -1
// 	const ASC = 1

// 	fmt.Printf("QueryCache filter: %s\n", spew.Sdump(pf))
// 	if pf.Limit > 0 {
// 		fmt.Printf("setting limit: %d\n",pf.Limit)
// 		limit = pf.Limit
// 	}
// 	if pf.Page > 0 {
// 		skip = (pf.Page - 1) * limit
// 		fmt.Printf("setting skip: %d\n",skip)
// 	}
// 	if pf.Skip > 0 {
// 		skip = pf.Skip
// 	}
// 	// q, err := pf.QueryCacheByEncounter()

// 	fmt.Printf("\nQueryCache: Paging limit: %d   skip:%d\n", limit, skip)
// 	fmt.Printf("\n\n@@@   QueryCache Using filter %s \n", pf.queryFilter)
// 	startTime := time.Now()
// 	// spew.Dump(q)
// 	// spew.Dump(f.QueryFilter)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(limit)
// 	findOptions.SetSkip(skip)
// 	//findOptions.SetSort(bson.D{bson.E{"family", 1}, bson.E{"given", 1}})
// 	//if pf.SortBy == "" || pv.SortBy {}
// 	sortOrder := ASC // Default Assending
// 	var sortFields bson.D

// 	if strings.ToLower(pf.Order) == "desc" {
// 		sortOrder = DESC
// 	}

// 	if len(pf.SortBy) > 0 {
// 		fmt.Printf("SortBy: %v\n", pf.SortBy)
// 		for _, s := range pf.SortBy {
// 			fmt.Printf("   Adding [%s]\n", s)

// 			sortFields = append(sortFields, bson.E{s, sortOrder})
// 		}
// 	}
// 	findOptions.SetSort(sortFields)
// 	fmt.Printf("sortFields: %s\n", sortFields)
// 	collection, _ := storage.GetCollection("fhir_patients")

// 	cursor, err := collection.Find(context.TODO(), pf.queryFilter, findOptions)

// 	var patients []*fhir.Patient

// 	if err != nil {
// 		log.Printf("QueryCache for %s returned error: %v\n", pf.queryFilter, err)
// 		//cursor is not open
// 		return nil, err
// 	}
// 	fmt.Printf("QueryCache took %f seconds\n", time.Since(startTime).Seconds())
// 	//defer cursor.Close(context.TODO)  // Need to get real context
// 	//fmt.Printf("\n    No Error on QueryCache\n\n")
// 	for cursor.Next(context.TODO()) {
// 		var patient fhir.Patient
// 		err = cursor.Decode(&patient)
// 		if err != nil {
// 			//fmt.Printf("   Next error: %v\n", err)
// 			cursor.Close(context.TODO())
// 			return nil, err
// 		}

// 		//fmt.Printf("  Added one to documents\n")
// 		patients = append(patients, &patient)
// 	}
// 	//fmt.Printf("   Finished fetching documents: error: %v\n\n", cursor.Err())
// 	cursor.Close(context.TODO())
// 	if len(patients) == 0 {
// 		err = fmt.Errorf("404|no documents found for %s", pf.queryFilter)
// 	} else {
// 		fmt.Printf("QueryCache found %d documents \n", len(patients))
// 	}
// 	return patients, err
// }

func (f *PatientFilter) makeCacheQueryFilter() {
	f.makeQueryMap()
	f.MakePatFHIRQueryString()
	f.queryFilter, _ = com.FilterFromMap(f.queryMap)
	//fmt.Printf("\n\ntoFilter returning: %v\n\n", filter)
	layout := "2006-01-02"
	mq := []bson.M{}
	for k := range f.queryMap {
		val := f.queryMap[k]
		//fmt.Printf("k: %s,  v: %s\n", k, q[k])
		if k == "enterprise_id" {
			//fmt.Printf("Converting search for id %s to search for enterpriseid\n", val)
			mq = append(mq, bson.M{"id": val})
		} else if k == "given" {
			mq = append(mq, bson.M{"firstname": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "family" {
			mq = append(mq, bson.M{"lastname": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
			// } else if k == "given:exact" {
			// 	mq = append(mq, bson.M{"firstname": primitive.Regex{strings.Replace("firstname", "\"", "", -1), "i"}})
		} else if k == "family:given" {
			mq = append(mq, bson.M{"lastname": primitive.Regex{strings.Replace("lastname", "\"", "", -1), "i"}})
		} else if k == "email" {
			mq = append(mq, bson.M{"email": primitive.Regex{strings.Replace(val, "\"", "", -1), "i"}})
		} else if k == "birthdate" {
			condition := "eq"
			input := ""
			s := strings.Split(val, "|")
			if len(s) > 1 {
				condition = s[0]
				input = s[1]
			} else {
				condition = "$eq"
				input = s[0]
			}
			useDate, err := time.Parse(layout, input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			q := bson.M{"birthdate": bson.M{condition: useDate}}
			mq = append(mq, q)

		} else {
			mq = append(mq, bson.M{k: val})
		}
	}

	if len(mq) > 0 {
		f.queryFilter = bson.M{"$and": mq}
	}
	fmt.Printf("\n\n@    queryFilter: %v\n", f.queryFilter)
	return
}

func (f *PatientFilter) makeQueryMap() error {
	fmt.Printf("Starting makeQueryMap: %s\n", f.queryMap)
	m := make(map[string]string)
	mrn := strings.Trim(f.MRN, " ")
	given := strings.Trim(f.Given, " ")
	family := strings.Trim(f.Family, " ")
	//given_exact := strings.Trim(f.GivenExact, " ")
	family_exact := strings.Trim(f.FamilyExact, " ")
	encounter := strings.Trim(f.EncounterID, " ")
	enterpriseID := strings.Trim(f.PatientGPI, " ")

	email := strings.Trim(f.Email, " ")
	//id := strings.Trim(f.ID, " ")
	birthdate := strings.Trim(f.BirthDate, " ")

	if family != "" {
		fmt.Printf("Family is set\n")
		m["family"] = family
		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
		f.SortBy = append(f.SortBy, "given")
	}
	if family_exact != "" {

		fam := strings.Split(family_exact, ":")
		fmt.Printf("Family_exact ")
		m["family"] = fam[0]
		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
		f.SortBy = append(f.SortBy, "given")
	}
	if given != "" {
		if family != "" || family_exact != "" {
			m["given"] = given

		} else {
			log.Warn("makeQueryMap:1230 - Invalid search: given alone is invalid")
			//return fmt.Errorf("400|makeQueryMap:1231 - Invalid search: given alone is invalid")
		}
	}

	if mrn != "" {
		m["mrn"] = mrn
	}
	if encounter != "" {
		m["encounter"] = encounter
	}
	if enterpriseID != "" {
		m["enterpriseid"] = enterpriseID
		m["id"] = enterpriseID
	}

	if email != "" {
		m["email"] = email
	}
	if birthdate != "" {
		m["birthdate"] = birthdate
	}
	f.queryMap = m
	fmt.Printf("created map: %v\n", m)
	return nil
}

func (f *PatientFilter) addQueryItem(name, value string) {
	field := ""
	if len(f.queryString) == 0 {
		field = fmt.Sprintf("?%s=%s", name, value)
	} else {
		field = fmt.Sprintf("&%s=%s", name, value)
	}
	f.queryString = f.queryString + field
}

func (f *PatientFilter) MakePatFHIRQueryString() {
	fmt.Printf("\n\nMakePatFHIRQueryString\n")
	if strings.Trim(f.Family, " ") != "" {
		log.Infof("Add Family: [%s]", f.Family)
		f.addQueryItem("family", strings.Trim(f.Family, " "))
		log.Infof("queryString: %s\n", f.queryString)
		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
	}
	if strings.Trim(f.FamilyExact, " ") != "" {
		f.addQueryItem("family:exact", strings.Trim(f.FamilyExact, " "))
		log.Infof("queryString: %s\n", f.queryString)
		f.SortBy = append(f.SortBy, "family") // if querying by names force a sort
	}
	if strings.Trim(f.Given, " ") != "" {
		f.addQueryItem("given", strings.Trim(f.Given, " "))
		log.Infof("queryString: %s\n", f.queryString)
		f.SortBy = append(f.SortBy, "given") // if querying by names force a sort
	}
	// if strings.Trim(f.GivenExact, " ") != "" {
	// 	f.addQueryItem("given:exact", strings.Trim(f.GivenExact, " "))
	// 	log.Infof("queryString: %s\n", f.queryString)
	// 	f.SortBy = append(f.SortBy, "given") // if querying by names force a sort
	// }
	// } else {
	// 	return fmt.Errorf("400|Invalid search: given alone is invalid")
	// }
	if strings.Trim(f.MRN, " ") != "" {
		f.addQueryItem("identifier", fmt.Sprintf("%s|%s", config.MrnID(), strings.Trim(f.MRN, " ")))
		log.Infof("queryString: %s\n", f.queryString)
	}
	if strings.Trim(f.PatientGPI, " ") != "" {
		f.addQueryItem("_id", strings.Trim(f.PatientGPI, " "))
	}
	log.Infof("FinalPatFHIRQueryString: %s\n", f.queryString)
}
func (f *PatientFilter) makeQueryString() {

	f.queryString = ""
	for k := range f.queryMap {

		if f.queryString == "" {
			f.queryString = fmt.Sprintf("%s=%s", k, f.queryMap[k])
		} else {
			f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, k, f.queryMap[k])
		}
	}
	log.Infof("queryString: %s\n", spew.Sdump(f))
	if f.Count != "" {
		f.queryString = fmt.Sprintf("%s&%s=%s", f.queryString, "_count", f.queryMap["count"])
	}
	fmt.Printf("makeQueryString returning:836 %s\n", f.queryString)
}

// Insert adds one emvironment to the pending. Checks if already exists and if there returns existing.
func (p *Patient) Insert(cacheName string) error {
	//fmt.Printf("adding: %T: %v\n\n", p, p)
	//fmt.Printf("add patient ID: %s  Name: %s to cache \n", p.ID, p.Name)
	p.setDates()
	p.SessionID = cacheName
	// _, err := FindByPhone(c.FaxNumber, c.Facility)
	// log.Fatal(err)
	collection, _ := storage.GetCollection("patients")

	insertResult, err := collection.InsertOne(context.TODO(), p)
	if err == nil {
		p.ID = insertResult.InsertedID.(primitive.ObjectID)
		// fmt.Printf(" id: %s\n", p.ID.String())
		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
	} else {
		//fmt.Printf(" id: %s already exists\n", p.EnterpriseID)
		err = nil
	}

	return err

}

// func InsertFhirPatResults(results *fhir.PatientResult, cacheName string) error {
// 	//entry := results.Entry

// 	for _, entry := range results.Entry{
// 		pat := entry.Patient

// 		err := InsertFhirPat( &pat, cacheName)
// 		if err != nil {
// 			msg := fmt.Sprintf("InsertFhirPat using cacheName: %s failed: %s",cacheName, err.Error())
// 			log.Error(msg)
// 			return  errors.New(msg)
// 		}
// 	}
// 	return nil
// }

// func InsertFhirPats(pats []*fhir.Patient, cacheName string) error {
// 	for _, pat := range pats {
// 		err := InsertFhirPat( pat, cacheName)
// 		if err != nil {
// 			msg := fmt.Sprintf("InsertFhirPat using cacheName: %s failed: %s",cacheName, err.Error())
// 			log.Error(msg)
// 			return  errors.New(msg)
// 		}
// 	}
// 	return nil
// }

// func InsertFhirPat(pat *fhir.Patient, cacheName string) error {
// 	//fmt.Printf("adding: %T: %v\n\n", p, p)
// 	//fmt.Printf("add patient ID: %s  Name: %s to cache \n", p.ID, p.Name)

// 	pat.SessionId = cacheName
// 	// _, err := FindByPhone(c.FaxNumber, c.Facility)
// 	// log.Fatal(err)
// 	collection, _ := storage.GetCollection("fhir_patients")

// 	insertResult, err := collection.InsertOne(context.TODO(), pat)
// 	if err == nil {
// 		pat.CacheId = insertResult.InsertedID.(primitive.ObjectID)
// 		// fmt.Printf(" id: %s\n", p.ID.String())
// 		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
// 	} else {
// 		//fmt.Printf(" id: %s already exists\n", p.EnterpriseID)
// 		err = nil
// 	}

// 	return err

// }

func (p *Patient) setDates() {
	t := time.Now()
	p.CreatedAt = t
	p.UpdatedAt = t
	p.AccessedAt = t
}

///////////////////////////////////////////////////////////////////////////////
//                          Cache managers
///////////////////////////////////////////////////////////////////////////////

func Insert(fp *fhir.Patient, sessionId string) error {
	//fmt.Printf("adding: %T: %v\n\n", p, p)
	//fmt.Printf("add patient ID: %s  Name: %s to cache \n", p.ID, p.Name)
	//p.setDates()

	fp.SessionId = sessionId
	fp.CacheID = primitive.NewObjectID()
	fp.LastAccess = time.Now().UTC()
	// _, err := FindByPhone(c.FaxNumber, c.Facility)
	// log.Fatal(err)
	collection, _ := storage.GetCollection("fhir_patients")
	//id := primitive.NewObjectID()
	_, err := collection.InsertOne(context.TODO(), fp)
	if err == nil {
		//p.ID = insertResult.InsertedID.(primitive.ObjectID)
		// fmt.Printf(" id: %s\n", p.ID.String())
		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
	} else {
		//fmt.Printf(" id: %s already exists\n", p.EnterpriseID)
		err = nil
	}

	return err

}

// Insert adds one emvironment to the pending. Checks if already exists and if there returns existing.
//func (p *Patient) Insert(cacheName string) error {
// 	//fmt.Printf("adding: %T: %v\n\n", p, p)
// 	//fmt.Printf("add patient ID: %s  Name: %s to cache \n", p.ID, p.Name)
// 	p.setDates()
// 	p.SessionID = p.SessionID
// 	// _, err := FindByPhone(c.FaxNumber, c.Facility)
// 	// log.Fatal(err)
// 	collection, _ := storage.GetCollection("capatients")

// 	insertResult, err := collection.InsertOne(context.TODO(), p)
// 	if err == nil {
// 		p.ID = insertResult.InsertedID.(primitive.ObjectID)
// 		// fmt.Printf(" id: %s\n", p.ID.String())
// 		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
// 	} else {
// 		//fmt.Printf(" id: %s already exists\n", p.EnterpriseID)
// 		err = nil
// 	}

// 	return err

// }

// func (p *CAPatient) Insert(sessionId string) error {
// 	//fmt.Printf("adding: %T: %v\n\n", p, p)
// 	fmt.Printf("\n\n\n\n#######  add patient SessionID: %s  Name: %s to cache \n", p.SessionID, p.Name)
// 	//
// 	//p.SessionID = p.SessionID

// 	// _, err := FindByPhone(c.FaxNumber, c.Facility)
// 	// log.Fatal(err)
// 	collection, _ := storage.GetCollection("capatients")

// 	insertResult, err := collection.InsertOne(context.TODO(), p)
// 	if err != nil {
// 		fmt.Printf(" id: %s already exists\n", p.PatientGPI)
// 		err = nil
// 	} else {
// 		p.ID = insertResult.InsertedID.(primitive.ObjectID)
// 		// fmt.Printf(" id: %s\n", p.ID.String())
// 		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
// 	}

// 	return err

// }

func InsertFhirPatResults(results *fhir.PatientResult, sessionId string) error {
	//entry := results.Entry

	for _, entry := range results.Entry {
		pat := entry.Patient

		err := InsertFhirPat(&pat, sessionId)
		if err != nil {
			msg := fmt.Sprintf("InsertFhirPat using sessionId: %s failed: %s", sessionId, err.Error())
			log.Error(msg)
			return errors.New(msg)
		}
	}
	return nil
}

func InsertFhirPats(pats []*fhir.Patient, sessionId string) error {
	for _, pat := range pats {
		err := InsertFhirPat(pat, sessionId)
		if err != nil {
			msg := fmt.Sprintf("InsertFhirPat using sessionId: %s failed: %s", sessionId, err.Error())
			log.Error(msg)
			return errors.New(msg)
		}
	}
	return nil
}

func InsertFhirPat(pat *fhir.Patient, sessionId string) error {
	//fmt.Printf("adding: %T: %v\n\n", p, p)
	//fmt.Printf("add patient ID: %s  Name: %s to cache \n", p.ID, p.Name)

	pat.SessionId = sessionId
	// _, err := FindByPhone(c.FaxNumber, c.Facility)
	// log.Fatal(err)
	collection, _ := storage.GetCollection("fhir_patients")
	pat.CacheID = primitive.NewObjectID()
	insertResult, err := collection.InsertOne(context.TODO(), pat)
	if err == nil {
		pat.CacheID = insertResult.InsertedID.(primitive.ObjectID)
		// fmt.Printf(" id: %s\n", p.ID.String())
		// fmt.Printf("\n\nAdded %s to cache\n", spew.Sdump(p))
	} else {
		//fmt.Printf(" id: %s already exists\n", p.EnterpriseID)
		err = nil
	}

	return err

}

//Get Chached entries first. If none, get them from FHIR
func GetCachedPatients(filter bson.M) ([]*Patient, error) {
	var patients []*Patient
	collection, err := storage.GetCollection("patients")
	if err != nil {
		fmt.Printf(" Error getting Collection: %s\n", err)
		return nil, err
	}
	//var encounter = new(Encounter)
	fmt.Printf("\nGetCachedPatients Filter: %v\n", filter)
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("GetCachedPatients search for %s returned error: %v\n", filter, err)
		cursor.Close(context.TODO())
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var patient Patient
		err = cursor.Decode(&patient)
		if err != nil {
			cursor.Close(context.TODO())
			return nil, err
		}
		//spew.Dump(patients)
		patients = append(patients, &patient)
	}
	if patients == nil {
		err = fmt.Errorf("404|no patients found for %s\n", filter)
	}
	return patients, err
}

// //Get Chached entries first. If none, get them from FHIR
func GetCachedPatient(filter bson.M) (*Patient, error) {
	var patient Patient
	collection, err := storage.GetCollection("patients")
	if err != nil {

		return nil, err
	}
	fmt.Printf("\nGetCachedPatient Checking for Patient with query: %v\n", filter)
	err = collection.FindOne(context.TODO(), filter).Decode(&patient)
	//fmt.Printf("FindOne returned err: %v\n", err)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

// func (pf *PatientFilter) QueryCache() ([]*fhir.Patient, error) {
// 	var limit int64 = 20
// 	var skip int64 = 0
// 	const DESC = -1
// 	const ASC = 1

// 	fmt.Printf("QueryCache filter: %s\n", spew.Sdump(pf))
// 	if pf.Limit > 0 {
// 		fmt.Printf("setting limit: %d\n", pf.Limit)
// 		limit = pf.Limit
// 	}
// 	if pf.Page > 0 {
// 		skip = (pf.Page - 1) * limit
// 		fmt.Printf("setting skip: %d\n", skip)
// 	}
// 	if pf.Skip > 0 {
// 		skip = pf.Skip
// 	}
// 	// q, err := pf.QueryCacheByEncounter()

// 	fmt.Printf("\nQueryCache: Paging limit: %d   skip:%d\n", limit, skip)
// 	fmt.Printf("\n\n@@@   QueryCache Using filter %s \n", pf.queryFilter)
// 	startTime := time.Now()
// 	// spew.Dump(q)
// 	// spew.Dump(f.QueryFilter)
// 	findOptions := options.Find()
// 	findOptions.SetLimit(limit)
// 	findOptions.SetSkip(skip)
// 	//findOptions.SetSort(bson.D{bson.E{"family", 1}, bson.E{"given", 1}})
// 	//if pf.SortBy == "" || pv.SortBy {}
// 	sortOrder := ASC // Default Assending
// 	var sortFields bson.D

// 	if strings.ToLower(pf.Order) == "desc" {
// 		sortOrder = DESC
// 	}

// 	if len(pf.SortBy) > 0 {
// 		fmt.Printf("SortBy: %v\n", pf.SortBy)
// 		for _, s := range pf.SortBy {
// 			fmt.Printf("   Adding [%s]\n", s)

// 			sortFields = append(sortFields, bson.E{s, sortOrder})
// 		}
// 	}
// 	findOptions.SetSort(sortFields)
// 	fmt.Printf("sortFields: %s\n", sortFields)
// 	collection, _ := storage.GetCollection("fhir_patients")

// 	cursor, err := collection.Find(context.TODO(), pf.queryFilter, findOptions)

// 	var patients []*fhir.Patient

// 	if err != nil {
// 		log.Printf("QueryCache for %s returned error: %v\n", pf.queryFilter, err)
// 		//cursor is not open
// 		return nil, err
// 	}
// 	fmt.Printf("QueryCache took %f seconds\n", time.Since(startTime).Seconds())
// 	//defer cursor.Close(context.TODO)  // Need to get real context
// 	//fmt.Printf("\n    No Error on QueryCache\n\n")
// 	for cursor.Next(context.TODO()) {
// 		var patient fhir.Patient
// 		err = cursor.Decode(&patient)
// 		if err != nil {
// 			//fmt.Printf("   Next error: %v\n", err)
// 			cursor.Close(context.TODO())
// 			return nil, err
// 		}

// 		//fmt.Printf("  Added one to documents\n")
// 		patients = append(patients, &patient)
// 	}
// 	//fmt.Printf("   Finished fetching documents: error: %v\n\n", cursor.Err())
// 	cursor.Close(context.TODO())
// 	if len(patients) == 0 {
// 		err = fmt.Errorf("404|no documents found for %s", pf.queryFilter)
// 	} else {
// 		fmt.Printf("QueryCache found %d documents \n", len(patients))
// 	}
// 	return patients, err
// }

////////////////////////////////////////////////////////////////////////////////
//                           Extract from FHIR
///////////////////////////////////////////////////////////////////////////////

func extractStringBirthDateDates(fp *fhir.Patient) string {
	layout := "2006-01-02"

	bdate, err := time.Parse(layout, fp.BirthDate)
	if err != nil {
		return ""
	} else {
		return bdate.Format(layout)
	}

	//ddate, err  := time.Parse(layout, p.DeathDate)
}

func extractCaBirthDate(fp *fhir.Patient) *time.Time {
	layout := "2006-01-02"

	bdate, err := time.Parse(layout, fp.BirthDate)
	if err != nil {
		return nil
	} else {
		return &bdate
	}

	//ddate, err  := time.Parse(layout, p.DeathDate)
}

func extractAddress(p *fhir.Patient, use string) *fhir.Address {
	addresses := p.Address

	for _, address := range addresses {
		if address.Use == use {
			return &address
		}
	}
	return nil
}

func extractTelecom(p *fhir.Patient, system, use string) string {

	for _, t := range p.Telecom {

		if t.System == system && t.Use == use {
			return t.Value
		}
	}
	return ""
}
func extractName(p *fhir.Patient, use string, element string) string {
	for _, name := range p.Name {
		if name.Use == use {
			switch element {
			case "family":
				return name.Family[0]
			case "given":
				return name.Given[0]
				// case "text":   // Full name
				// 	return name.
			}
		}
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////
//                           Convert to Summary
////////////////////////////////////////////////////////////////////////////////

func FhirPatientToSum(p *fhir.Patient, cacheName string) *Patient {
	//fmt.Prinln("   parsePatient entered")
	ids := extractIDs(p.Identifier)
	var pat Patient
	pat.Family = strings.Join(p.Name[0].Family, " ")
	pat.Given = strings.Join(p.Name[0].Given, " ")
	pat.Name = fmt.Sprintf("%s, %s", pat.Family, pat.Given)
	pat.Sex = p.Gender
	pat.EnterpriseID = p.ID
	pat.MRN = ids["MRN"]
	pat.SSN = ids["SSN"]
	pat.MaritalStatus = p.MaritalStatus.Text
	pat.Source = config.Source()

	pat.Text.Div = p.Text.Div
	pat.Text.Status = p.Text.Status
	dob := extractCaBirthDate(p)
	pat.BirthDate = *dob
	address := extractAddress(p, "home")
	if address != nil {
		switch len(address.Line) {
		case 1:
			pat.Address1 = address.Line[0]
		case 2:
			pat.Address1 = address.Line[0]
			pat.Address2 = address.Line[1]
		}
		pat.City = address.City
		pat.State = address.State
		pat.PostalCode = address.PostalCode
		pat.Country = address.Country
	}
	pat.HomePhone = extractTelecom(p, "phone", "home")
	pat.CellPhone = extractTelecom(p, "phone", "cell")
	pat.WorkPhone = extractTelecom(p, "phone", "work")
	pat.Email = extractTelecom(p, "email", "home")
	//fmt.Printf("BirthDate: %s\n", p.BirthDate)
	err := pat.Insert(cacheName)
	if err != nil {
		fmt.Printf("Patient Insert failed: %s\n", err.Error())
	}
	return &pat
}

func FhirPatientsToSum(fps []*fhir.Patient, cacheName string) []*Patient {
	patients := []*Patient{}
	for _, ps := range fps {
		pat := FhirPatientToSum(ps, cacheName)
		patients = append(patients, pat)
	}
	return patients
}

// ////////////////////////////////////////////////////////////////////////////////
// //                           Convert to CA
// ////////////////////////////////////////////////////////////////////////////////
// //

// func FhirPatientsToCA(fps []*fhir.Patient, sessionId string) ([]*CAPatient, error) {
// 	caPatients := []*CAPatient{}
// 	for _, p := range fps {
// 		pat, err := FhirPatientToCA(p, sessionId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		caPatients = append(caPatients, pat)
// 	}
// 	return caPatients, nil
// }

// func FhirPatientToCA(fp *fhir.Patient, sessionId string) (*CAPatient, error) {
// 	//fmt.Prinln("   parsePatient entered")
// 	ids := extractIDs(fp.Identifier)
// 	var caPat CAPatient
// 	caPat.SessionID = sessionId
// 	caPat.LastName = extractName(fp, "official", "family")
// 	caPat.FirstName = extractName(fp, "official", "given")
// 	//caPat.Name = extractName(*fp, "official", "text")
// 	caPat.Name = fmt.Sprintf("%s, %s", caPat.LastName, caPat.FirstName)
// 	caPat.Sex = fp.Gender
// 	ident := Identifier{}
// 	ident.System = "cerner" // TODO: Get ID.System from config
// 	ident.Value = fp.ID
// 	ident.Facility = ""
// 	caPat.Identifiers = append(caPat.Identifiers, &ident)

// 	caPat.PatientGPI = fp.ID
// 	caPat.MRN = ids["MRN"]
// 	caPat.SSN = ids["SSN"]
// 	caPat.MaritalStatus = fp.MaritalStatus.Text
// 	caPat.Source = strings.ToLower(config.Source())
// 	caPat.Text.Div = fp.Text.Div
// 	caPat.Text.Status = fp.Text.Status
// 	caPat.BirthDate = extractCaBirthDate(fp)

// 	address := extractAddress(fp, "home")
// 	if address != nil {
// 		switch len(address.Line) {
// 		case 1:
// 			caPat.Address1 = address.Line[0]
// 		case 2:
// 			caPat.Address1 = address.Line[0]
// 			caPat.Address2 = address.Line[1]
// 		}
// 		caPat.City = address.City
// 		caPat.State = address.State
// 		caPat.PostalCode = address.PostalCode
// 		caPat.Country = address.Country
// 	}

// 	caPat.Email = extractTelecom(fp, "email", "home")
// 	caPat.HomePhone = extractTelecom(fp, "phone", "home")
// 	caPat.CellPhone = extractTelecom(fp, "phone", "cell")
// 	caPat.WorkPhone = extractTelecom(fp, "phone", "work")

// 	fmt.Printf("-------FhirPatientToCA:1553 calling caPat.Insert\n")
// 	err := caPat.Insert(sessionId)
// 	//fmt.Printf("BirthDate: %s\n", p.BirthDate)
// 	//err := pat.Insert(cacheName)
// 	if err != nil {
// 		fmt.Printf("caPatient Insert failed: %s\n", err.Error())
// 		return nil, err

// 	}
// 	fmt.Printf("--------FhirPatientToCA:1562 Insert Cache Returning  ID: %s  SessionId: %s\n", caPat.ID.Hex(),  caPat.SessionID)
// 	return &caPat, nil
// }

func (pf *PatientFilter) DeleteCachedPatients() {
	sessionid := pf.Session.PatSessionId
	startTime := time.Now()
	fmt.Printf("Deleting Patients for session %s\n", sessionid)
	collection, _ := storage.GetCollection("patients")
	filter := bson.D{{"sessionid", sessionid}}
	fmt.Printf("    bson filter delete: %v\n", filter)
	deleteResult, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		fmt.Printf("DeletePatients for session %s failed: %vn", sessionid, err)
		return
	}
	fmt.Printf("Deleted %v Patients for session: %v in %s\n", deleteResult.DeletedCount, sessionid, time.Since(startTime))
}

func NextPage(links []fhirongo.Link) {
	url := NextPageLink(links)
	if url == "" {
		return
	}
}

func NextPageLink(links []fhirongo.Link) string {
	for _, l := range links {
		//fmt.Printf("Looking at link: %s\n", spew.Sdump(l))
		if l.Relation == "next" {
			//fmt.Printf("FOUND IT: \n")
			return l.URL
		}
	}
	return ""
}
