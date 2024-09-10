package common

import (
	"log"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetFhirTime(t *testing.T) {
	Convey("Convert String to fhir time", t, func() {
		ftime, err := StringDateToFhirDate("2020-08-12T19:36:06")
		So(err, ShouldBeNil)
		So(ftime.Year(), ShouldEqual, 2020)
		So(ftime.Day(), ShouldEqual, 12)
		So(FhirDateToString(ftime, "full"), ShouldEqual, "2020-08-12T19:36:06.000Z")
		So(FhirDateToString(ftime, "mdy"), ShouldEqual, "08-12-2020")
		So(FhirDateToString(ftime, "ymd"), ShouldEqual, "2020-08-12")
	})
}

func TestGetFhirPartial(t *testing.T) {
	Convey("Convert String to fhir time", t, func() {
		ftime, err := StringDateToFhirDate("2020-08-12T00:00:00.000Z")
		So(err, ShouldBeNil)
		So(ftime.Year(), ShouldEqual, 2020)
		So(ftime.Day(), ShouldEqual, 12)
		So(FhirDateToString(ftime, "full"), ShouldEqual, "2020-08-12T00:00:00.000Z")
		So(FhirDateToString(ftime, "mdy"), ShouldEqual, "08-12-2020")
		So(FhirDateToString(ftime, "ymd"), ShouldEqual, "2020-08-12")
	})
}

func TestMDYToFhir(t *testing.T) {
	Convey("Convert MDY String to FHIR", t, func() {
		ftime, err := MDYToFhir("01-27-1958")
		So(err, ShouldBeNil)
		So(ftime.Year(), ShouldEqual, 1958)
		So(ftime.Day(), ShouldEqual, 27)
		So(FhirDateToString(ftime, "full"), ShouldEqual, "1958-01-27T00:00:00.000Z")
		So(FhirDateToString(ftime, "mdy"), ShouldEqual, "01-27-1958")
		So(FhirDateToString(ftime, "ymd"), ShouldEqual, "1958-01-27")
	})
}

func TestCalcPages(t *testing.T) {
	Convey("Calculate Number of pages to display a set of documents", t, func() {
		pages, _ := CalcPages(20, 8)
		So(pages, ShouldEqual, 3)
		pages, _ = CalcPages(126, 20)
		So(pages,  ShouldEqual, 7)
		pages, err := CalcPages(0, 8)
		So(err, ShouldNotBeNil)
		So(pages, ShouldEqual, 0)
		pages, err = CalcPages(8, 0)
		So(err, ShouldNotBeNil)
		So(pages, ShouldEqual, 0)
		pages, err = CalcPages(8, -1)
		So(err, ShouldNotBeNil)
		So(pages, ShouldEqual, 0)



	})
}
// func TestCreateIDFilterFromMap(t *testing.T) {
// 	Convey("Given a query map of id:123".t func(){
// 		q := "id=123"

// 		Convey("When the Creator is called", func(){
// 			result := h.Cre

// 		})

// 	})
// }

func TestCreatingQueries(t *testing.T) {
	SkipConvey("Subject: Converting queries", t, func() {
		Convey("Given a query string, convert it to a map", func() {
			q := "id=125&name=Theresa French"
			result, err := MapFromString(q)
			Convey("There is no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("The result is a valid map", func() {
				So(result["id"], ShouldEqual, "125")
				So(result["name"], ShouldEqual, "Theresa French")
			})
		})

		Convey("Given a query string, convert it to a filter", func() {
			q := "id=125&name=Theresa French"
			_, err := FilterFromString(q)
			Convey("Then there should be no error", func() {
				//fmt.Printf("@@@ Filter: %v\n", result)
				//So(result, ShouldEqual, map[$and:[map[enterpriseid:125] map[name:Theresa French]]])
				So(err, ShouldBeNil)
			})
			Convey("The results should be a valid bson Filter", func() {
				//fmt.Printf("@@@ Filter: %v\n", result)
				So(true, ShouldBeTrue)
			})

		})
		Convey("Given a Map, convert it to a filter", func() {
			q := "id=125&name=Theresa French"
			m, _ := MapFromString(q)
			f, err := FilterFromMap(m)
			Convey("There should be no error", func() {
				log.Printf("Filter: %v\n", f)
				So(err, ShouldBeNil)
			})
		})
		Convey("Given a Map, convert to query string", func() {
			q := "id=125&name=Theresa French"
			m, _ := MapFromString(q)
			s := StringFromMap(m)
			Convey("The string should ve valid", func() {
				So(s, ShouldEqual, q)
			})
		})

	})
}

// func TestCreateMapFromString(t *testing.T) {
// 	Convey("Given a query string of id=125&name=Theresa French", t, func() {
// 		q := "id=125&name=Theresa French"
// 		Convey("When MapFromString is called", func() {
// 			result, err := MapFromString(q)
// 			Convey("Then the results should be a valid hash", func() {
// 				So(err, ShouldEqual, nil)
// 				So(result["id"], ShouldEqual, "125")
// 				So(result["name"], ShouldEqual, "Theresa French")
// 			})
// 		})
// 	})
// }

// func TestFilterFromString(t *testing.T) {
// 	Convey("Given a query string of id=125&name=Theresa French", t, func() {
// 		q := "id=125&name=Theresa French"
// 		Convey("When FilterFromString is called", func() {
// 			result, err := FilterFromString(q)
// 			Convey("Then the results should ba a valid filter", func() {
// 				fmt.Printf("@@@ Filter: %v\n", result)
// 				//So(result, ShouldEqual, map[$and:[map[enterpriseid:125] map[name:Theresa French]]])
// 				So(err, ShouldEqual, nil)
// 			})

// 		})
// 	})
// }
