package model

import (
	"fmt"
	"testing"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

//func TestGetEncounters(t *testing.T) {
// 	session := setupTest("")
// 	Convey("Given there is an Authorized Session ", t, func() {
// 	//session, _ := ValidateAuth("test")
// 		//		DeleteDocuments(session.CacheName)
// 		fmt.Printf("\n\n     #### TestGetEncounters session: [%s]\n\n", session.EncSessionId)
// 		Convey("Find encounters for a patient", func() {
// 			Convey("Find encounters for 4342009", func() {
// 				ef := EncounterFilter{PatientID: "4342009"}
// 				encounters, err := ef.SearchEncounters()
// 				//fmt.Printf("Number received: %d\n", len(encounters))
// 				Convey("Then the results should be a Nancy Smart", func() {
// 					So(err, ShouldBeNil)
// 					//spew.Dump(encounters)
// 					So(len(encounters), ShouldEqual, 7)
// 					So(encounters[0].PatientName, ShouldEqual, "SMART, NANCY")
// 				})
// 			})
// 		})
// 		Convey("Find encounters for a patient before 2018-06-30", func() {
// 			Convey("Setup the search", func() {
// 				ef := EncounterFilter{PatientID: "4342009", AdmitDate: "$lt|2018-03-01"}
// 				encounters, err := ef.SearchEncounters()
// 				//fmt.Printf("Number received: %d\n", len(encounters))
// 				Convey("Then there should be 1 of the 7 Encounters", func() {
// 					So(err, ShouldBeNil)
// 					So(len(encounters), ShouldEqual, 1)
// 				})
// 			})
// 		})
// 	})
// }

func TestSearchEncounters(t *testing.T) {
	session := setupTest("")
	Convey("Given there is an Authorized Session ", t, func() {
	//session, _ := ValidateAuth("test")
		//		DeleteDocuments(session.CacheName)
		fmt.Printf("\n\n     #### TestGetEncounters session: [%s]\n\n", session.EncSessionId)
		Convey("Find encounters for a patient", func() {
			Convey("Find encounters for 4342009", func() {
				ef := EncounterFilter{PatientID: "12765407", Session: session}
				encounters, err := ef.SearchEncounters()
				//fmt.Printf("Number received: %d\n", len(encounters))
				Convey("Then the results should be a Nancy Smart", func() {
					So(err, ShouldBeNil)
					//spew.Dump(encounters)
					So(len(encounters), ShouldBeGreaterThan, 0)
					So(encounters[0].Patient.Display, ShouldEqual, "SMART, NANCY")
					//PatientName, ShouldEqual, "SMART, NANCY")
				})
			})
		})
	})
}

func TestGetFhirEncounterByID(t *testing.T) {
	session := setupTest("")
	Convey("Given there is an Authorized Session ", t, func() {
		fmt.Printf("\n\n     #### TestGetEncounter session: [%s]\n", session.EncSessionId)
		Convey("Find encounters for a patient", func() {
			Convey("Find encounters for 4342009", func() {
				encounter, err := GetFhirEncounterByID("97964891")
				//fmt.Printf("Number received: %d\n", len(encounters)){
				So(err, ShouldBeNil)
				So(encounter, ShouldNotBeNil)
				So(encounter.Patient.Display, ShouldEqual, "SMART, NANCY")
				//fmt.Printf("Enconter: %s\n", spew.Sdump(encounter))
				
			})
		})
	})
}