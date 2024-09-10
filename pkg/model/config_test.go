package model

import (
	//"fmt"
	"os"
	"testing"

	//mongo "github.com/dhf0820/cernerFhir/pkg/storage"
	//"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInitConfig(t *testing.T) {
	Convey("Given an uninitialized app", t, func() {
		os.Setenv("DATABASE", "fhir")
		Convey("When InitConfig is called", func() {
			name := "epic"
			mode := "dev"
			//mongoURL := "mongodb://dhf:Sacj0nhat1@linode.ihids.com:27017/?serverSelectionTimeoutMS=5000&connectTimeoutMS=10000&authSource=admin&authMechanism=SCRAM-SHA-256"
			mongoURL := "mongodb+srv://dhfadmin:Sacj0nhati@cluster1.24b12.mongodb.net/fhir?retryWrites=true&w=majority"

			config := initConfig(name, mode, mongoURL)

			//Convey("There should be an active Configuration", func() {
			So(config, ShouldNotEqual, nil)
			//fmt.Printf("Config: %s\n", spew.Sdump(config))
			//So(ActiveConfig(), ShouldEqual, config)
			//})
		})
	})
}

// func TestSimple(t *testing.T) {
// 	Convey("Given Simple Test", t, func() {
// 		i := 1
// 		Convey("When Tested", func() {
// 			//	Convey("There should be a result", func() {
// 			i = 2
// 			So(i, ShouldEqual, i)
// 			//	})
// 		})
// 	})
// }

// func TestGetDocument(t *testing.T) {
// 	Convey("Given there is a Database with a document", t, func() {
// 		var url = "postgres://chartarchive:Sacj0nhat!@vertisoft.com:5432/chartarchive_dev"
// 		os.Setenv("POSTGRES_URL", url)
// 		//err := fmt.Errorf("Test %s", " some value")
// 		//driver.Open("")
// 		//driver.SetSchema("demo")
// 		Convey("When a valid document is requested", func() {
// 			// doc := Document{}
// 			// doc.DocID = 72
// 			// d, err := doc.Get()
// 			//So(err, ShouldEqual, nil)
// 			// So(d.MedRecNum.String, ShouldEqual, "ll0819")
// 			// So(d.Description.String, ShouldEqual, "Test OP Lab Document")
// 			i := 1
// 			So(i, ShouldEqual, 1)
// 		})
// 	})
// }
