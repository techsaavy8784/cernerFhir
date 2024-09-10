package model

import (
	"fmt"
	//"os"
	//"strings"
	"testing"
	"time"
	//log "github.com/sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPatientExactName(t *testing.T) {
	as :=setupTest("")

	Convey("Subject: patients by Exact family and given name", t, func() {
		token, err := CreateToken("dhf")
		So(err, ShouldBeNil)
		So(token, ShouldNotBeNil)
		
		fmt.Printf("\n\n     #### Setup PatientByExactlName session: [%s]\n\n", as.DocSessionId)

		pf := PatientFilter{Family: "smart",  Mode: "ca", Cache: "reset", SessionId: as.DocSessionId}
		Convey("Given a query string of family_exact=smart&given=baby", func() {
			//pf := PatientFilter{FamilyExact: "smart", Given: "babyb", Session: *session, Mode: "ca"}
			//pf := PatientFilter{Family: "smart", Session: *session, Mode: "ca"}

			Convey("When GetPatient is called", func() {
				fmt.Printf("Performing the search: %s\n", spew.Sdump(pf))
				patients, inPage, pages, totalInCache, err := pf.CaPatientSearch()
				So(err, ShouldBeNil)
				So(inPage, ShouldEqual, 20)
				So(pages, ShouldEqual, 3)
				So(totalInCache, ShouldEqual, 20)
				So(patients, ShouldNotBeNil)
				//fmt.Printf("Patients: %s\n", spew.Sdump(patients))
				// patient := patients[0]
				// Convey("Then the results should be a BabyBoy Smart", func() {
				// 	So(err, ShouldBeNil)
				// 	So(patient.Name, ShouldEqual, "SMART, babygirl")
				// })
			})
		}) 
		time.Sleep(30 * time.Second)
		pf.Page = 4
		caPats, inPage, pages, totalInCache, err := pf.GetPatientPage()
		if err != nil {
			fmt.Printf("err: %s\n", err.Error())
		}
		fmt.Printf("InPage: %d, pages: %d, total: %d\n",inPage, pages, totalInCache)
		fmt.Printf("Patient InPage: %d   patients: %s\n", len(caPats), spew.Sdump(caPats[0]))
	})
}

// func TestPatientName(t *testing.T) {
// 	session := setupTest("")

// 	Convey("Subject: patients by family and given name", t, func() {
	
// 		//session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup PatientByFamilyName session: [%s]\n\n", session.CacheName)
// 		DeletePatients(session.SessionID)
// 		Convey("Given a query string of family=smart&given=b", func() {
// 			pf := PatientFilter{Family: "sm", Given: "", Session: *session}
// 			Convey("When GetPatient is called", func() {
// 				fmt.Printf("Performing the search: %s\n", spew.Sdump(pf))
// 				patients, err := pf.Search()
// 				So(err, ShouldBeNil)
// 				So(patients, ShouldNotBeNil)
// 				patient := patients[0]
// 				time.Sleep(30 * time.Second)
// 				Convey("Should be able to get page 2", func(){
// 					pf.Page = 2
// 					patients, err := pf.Search()
// 					So(err, ShouldBeNil)
// 					So(patients, ShouldNotBeNil)
// 					So(patient.ID, ShouldNotEqual, patients[0].ID)
// 				})
// 				Convey("Then the results should be a BabyBoy Smart", func() {
// 					So(err, ShouldBeNil)
// 					So(len(patients), ShouldEqual, 20)
// 					//So(patient.Name, ShouldEqual, "SMART, BABY BOY")
// 				})
// 			})
// 		})
// 	})
// }

// func TestPatientInvalidName(t *testing.T) {
// 	setupTest("")
// 	Convey("Subject: patients by Invalid Name", t, func() {
// 		session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup PatientInvalidName session: [%s]\n\n", session.CacheName)
// 		Convey("Given a non existing Patient Name", func() {
// 			pf := PatientFilter{Family: "zsma", Given: "Na", Session: *session}
// 			Convey("When GetPatient is called", func() {
// 				patients, err := pf.Search()
// 				Convey("Then the results should be nothing", func() {
// 					s := strings.Split(err.Error(), "|")
// 					So(s[0], ShouldEqual, "404")
// 					So(len(patients), ShouldEqual, 0)
// 				})
// 			})
// 		})
// 	})
// }

// func TestPatientByMRN(t *testing.T) {
// 	setupTest("")
// 	Convey("Subject: patients by MRN", t, func() {
// 		session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup PatientByMRN session: [%s]\n\n", session.CacheName)
// 		Convey("Given a query by mrn", func() {
// 			pf := PatientFilter{MRN: "6930", Session: *session}

// 			caPats, inPage, pages, totalInCache, err := pf.CaSearch()
// 			So(err, ShouldEqual, nil)
// 			So(len(caPats), ShouldBeGreaterThan, 0)
// 			So(len(caPats), ShouldEqual, 1)
// 			//ids := extractIDs(patients[0].Identifier)
// 			//So(ids["MRN"], ShouldEqual, "6930")	
// 				//So(patients[0].Name, ShouldEqual, "Creevey, Colin Carl")
// 			//caPat :=FhirPatientToCA(patients[0])
// 			//So(caPat.MRN, ShouldEqual, "6930")
// 			//pat := patients[0]
// 			// := FhirPatientToCA(pat)
// 			So(caPats, ShouldNotBeNil)
// 			So(caPats[0].MRN, ShouldEqual, "6930")		
// 		})
// 	})
// }

// func TestPatientByQueryID(t *testing.T) {
// 	setupTest("")
// 	Convey("Subject: patients by ID Query", t, func() {
// 		session, _ := ValidateAuth("test")
// 		fmt.Printf("\n\n     #### Setup PatientByID session: [%s]\n\n", session.CacheName)
// 		Convey("given a query by id 12724066", func() {
// 			pf := PatientFilter{PatientGPI: "12724066", UseCache: "false", Session: *session}
// 			Convey("The request is processed", func() {
// 				fmt.Printf("Request is processed\n")
// 				results, err := pf.Search()
// 				patients := results
// 					So(err, ShouldBeNil)
// 					So(patients[0].Name, ShouldEqual, "SMART, NANCYU")
// 			})
// 		})
// 	})
// }
func TestPatientByGPI(t *testing.T) {
	as := setupTest("")
	Convey("Subject: patients by ID Query", t, func() {
		//session, _ := ValidateAuth("test")
		fmt.Printf("\n\n     #### Setup PatientByID as: [%s]\n\n", as.DocSessionId)
		Convey("given a query by id 12724066", func() {
			pf := PatientFilter{PatientGPI: "12724066"}
			fmt.Printf("Request is processed\n")
			results, err := pf.Search()
			patients := results
			So(err, ShouldBeNil)
			So(patients[0].Name, ShouldEqual, "SMART, NANCYU")
		})
	})
}

func TestPatientByPatientGPI(t *testing.T) {
	//godotenv.Load("env_test")
	as := setupTest("")
	Convey("Subject: patient by ID Query", t, func() {
		So(as, ShouldNotBeNil)
		patient, err := ForPatientGPI("12742397")
		So(err, ShouldBeNil)
		So(patient, ShouldNotBeNil)
		So(patient.ID, ShouldEqual,"12742397" )
		fmt.Printf("Patient: %s\n", spew.Sdump(patient.Name)) 
		firstName := extractName(patient,"official", "given")
		So(firstName, ShouldEqual, "BABY BOY")
	})
}

// func TestMakeCacheFilter(t *testing.T) {
// 	session := setupTest("")
// 	Convey("Subject: Proper CacheFilter is created", t, func() {
// 		//session, _ := ValidateAuth("test")
// 		//fmt.Printf("\n\n     #### Setup PatientByEncounter session: [%s]\n\n", session.CacheName)
// 		Convey("Given: a sessionId", func() {
// 			pf := PatientFilter{Session: *session, Family:"smart"}
// 			//err := pf.MakeCacheFilter()
// 			err := pf.MakeCacheFilter()
// 			So(err, ShouldBeNil)
// 			fmt.Printf("cacheFilter: %s\n",spew.Sdump(pf))
// 		})
// 	})
// }

// func TestQueryCashe(t *testing.T) {
// 	session := setupTest("")
// 	Convey("Subject: PatientFilter returns selected from Cache", t, func() {
// 		//session, _ := ValidateAuth("test")
// 		//fmt.Printf("\n\n     #### Setup PatientByEncounter session: [%s]\n\n", session.CacheName)
// 		pf := PatientFilter{Session: *session, Family:"smart", Limit: 2, Count: "5"}
// 		Convey("Given: a valid filter", func() {
// 			caPats, err := pf.CaSearch()
// 			So(err, ShouldBeNil)
// 			So(caPats, ShouldNotBeNil)
// 			//caPats := FhirPatientsToCA(pats)
// 			//So(caPats, ShouldNotBeNil)
// 			fmt.Printf("caPats: %s\n", spew.Sdump(caPats))
// 			fmt.Printf("Number of caPats: %d\n", len(caPats))
// 		})
// 	})
// }


// func TestCaSearch(t *testing.T) {
// 	session := setupTest("")
// 	Convey("Subject: PatientFilter returns selected from Cache", t, func() {
// 		//session, _ := ValidateAuth("test")
// 		//fmt.Printf("\n\n     #### Setup PatientByEncounter session: [%s]\n\n", session.CacheName)
// 		pf := PatientFilter{Session: *session, Family:"smart", Limit: 2, Count: "5"}
// 		Convey("Given: a valid filter", func() {
// 			caPats, err := pf.CaSearch()
// 			So(err, ShouldBeNil)
// 			So(caPats, ShouldNotBeNil)
// 			//caPats := FhirPatientsToCA(pats)
// 			//So(caPats, ShouldNotBeNil)
// 			//fmt.Printf("caPats: %s\n", spew.Sdump(caPats))
// 			fmt.Printf("Number of caPats: %d\n", len(caPats))
// 		})
// 	})
// }
