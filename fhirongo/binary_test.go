package fhirongo

import (
	//log "github.com/sirupsen/logrus"
	//. "github.com/smartystreets/goconvey/convey"

	"fmt"
	//"os"
	"testing"
	//"time"
	//"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
	//"github.com/davecgh/go-spew/spew"
	//log "github.com/sirupsen/logrus"
)

//const pid = "4342009"
//const baseurl = "https://fhir-open.cerner.com/dstu2/ec2458f2-1e24-41c8-b71b-0e701af7583d/"

func TestGetImage(t *testing.T) {
	fmt.Printf("Get a PDF image\n")
	c := New(baseurl)
	Convey("Request the PDF image by documentId", t, func(){
		image, err := c.GetImage("197198634")
		So(err, ShouldBeNil)
		So(image, ShouldNotBeNil)
		So(image.ContentType, ShouldEqual, "application/pdf")
	})
}

func TestGetPDF(t *testing.T) {
	fmt.Printf("Get a PDF image\n")
	c := New(baseurl)
	Convey("Request the PDF image by documentId", t, func(){
		image, err := c.GetPDF("197466431")
		So(err, ShouldBeNil)
		So(image, ShouldNotBeNil)

	})
}
func TestGetPDFb64(t *testing.T) {
	fmt.Printf("Get a PDF image\n")
	c := New(baseurl)
	Convey("Request the PDF image by documentId", t, func(){
		docId := "197369077"
		image, err := c.GetPDFb64(docId)
		So(err, ShouldBeNil)
		So(image, ShouldNotBeNil)
		Decode(docId+".pdf", image )
	})
}
