package fhirongo

import (
	//log "github.com/sirupsen/logrus"
	//. "github.com/smartystreets/goconvey/convey"

	"fmt"
	//"os"
	"testing"
	//"time"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
	//"github.com/davecgh/go-spew/spew"
	//log "github.com/sirupsen/logrus"
)

//const pid = "4342009"
const baseurla = "https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/"

func TestGetDiagnostic(t *testing.T) {
	fmt.Printf("Test run a FHIR query\n")
	c := New(baseurl)
	Convey("Run a query", t, func(){

		fmt.Printf("GetDiagReport\n")
		//url := fmt.Sprintf("%sDiagnosticReport?patient=12724066",baseurla)
		// data, err := c.GetDiagnosticReports("?patient=12724066")
		// So(err, ShouldBeNil)
		// So(data, ShouldNotBeNil)
		data, err := c.GetPatientDiagnosticReports("12724066")
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	fmt.Printf("Data: %s\n", spew.Sdump(data))

	})
}
//https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/DiagnosticReport?patient=12714066
//https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/DiagnosticReport?patient=12724066