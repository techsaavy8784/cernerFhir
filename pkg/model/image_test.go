package model

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRetrieveImage(t *testing.T) {
	setupTest("")
	Convey("Given a valid dopcument id ", t, func() {
		docId := "197198634"
		image, err := RetrieveImage(docId)
		So(err, ShouldBeNil)
		So(image, ShouldNotBeNil)
		pdf, err := base64.StdEncoding.DecodeString(image.Content)	
		fileName := docId+".pdf"
		err = ioutil.WriteFile(fileName, []byte(pdf), 0666)
		if err != nil {
			fmt.Printf("Error writing pdf file: %s\n", err.Error())
		}
		So(len(image.Content), ShouldEqual, 16192)
	})
}