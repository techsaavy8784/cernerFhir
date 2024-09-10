package handler

import (
	"encoding/base64"
	//"encoding/json"
	//http "net/http"
	//"net/http"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	//ca "github.com/dhf0820/cernerFhir/pkg/ca"
	//m "github.com/dhf0820/cernerFhir/pkg/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetImage(t *testing.T) {
	as := setupTest("")

	Convey("Subject: GetImage responds properly", t, func() {
		Convey("Given a request for an image for document 197466431", func() {
			req := httptest.NewRequest("GET", "/api/rest/v1/image/197466431", nil)

			//w := httptest.NewRecorder()
			Convey("When the request is handled by the router", func() {
				//req := httptest.NewRequest("GET", "/api/v1/documents?mrn=10002701&mode=ca&limit=20&page=1", nil)
				req.Header.Set("SESSION", as.DocSessionId)
				resp := httptest.NewRecorder()
				NewRouter().ServeHTTP(resp, req)
				So(resp.Result().StatusCode, ShouldEqual, 200)
				fmt.Printf("Length of body: %d\n", resp.Body.Len())
				content, err := ioutil.ReadAll(resp.Body)
				So(err, ShouldBeNil)
				Decode("197466431.pdf", string(content))
				//json.NewDecoder(resp.Body)

				NewRouter().ServeHTTP(resp, req)
				//content, err := ioutil.ReadFile("../../sample.pdf")     // the file is inside the local directory
				//content, err := ioutil.ReadFile("./test1.txt")     // the file is inside the local directory

				fmt.Printf("Length of data: %d\n", len(content))
				// if err != nil {
				// 	fmt.Println("Err")
				// 	WriteImageResponse(w, 400, nil)
				// 	return
				// }

				//Decode()
				// GetImage(w, req)
				// res := w.Result()
				// defer res.Body.Close()
				// data, err := ioutil.ReadAll(res.Body)
				// So(err, ShouldBeNil)
				//pdf:= base64.StdEncoding.EncodeToString(content)
				//pdf, err := base64.StdEncoding.DecodeString(string(content))
				// So(err, ShouldBeNil)
				//fmt.Printf("saving %d bytes\n", len(pdf))
				//WriteBytesToFile(pdf)
				//WriteStringToFile("./test3.pdf", string(pdf))
				// if err := os.WriteFile("file.pdf", string(data), 0666); err != nil {
				// 	log.Fatal(err)
				// }
				// f, err := os.Create("./file.pdf")
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// // remember to close the file
				// defer f.Close()

				// // write bytes to the file
				// _, err = f.Write(data)
				// if err != nil {
				// 	log.Fatal(err)
				// }

				// NewRouter().ServeHTTP(resp, req)
				// defer resp.Body.Close()
				// b, _ := ioutil.ReadAll(resp.Body)
				// fmt.Printf("Err: %v\n")
				// fmt.Printf("b: %v\n", string(b))
				// docs, _ := DocumentResults(resp)

				// Convey("Then the response should be a 200", func() {
				// 	//So(true, ShouldBeTrue)
				// 	So(resp.Code, ShouldEqual, 200)
				// 	//So(encs[0].Name, ShouldEqual, "Creevey, Colin Carl")
				// 	//spew.Dump(b)
				// })
			})
		})
	})
}

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
}

func WriteBytesToFile(byteSlice []byte) {
	file, err := os.OpenFile(
		"test.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Write bytes to file
	// byteSlice := []byte("Bytes!\n")
	bytesWritten, err := file.Write(byteSlice)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d bytes.\n", bytesWritten)
}

func Decode(fileName string, b64 string) {
	//b64, _ := ioutil.ReadFile("../../sample.pdf")
	//b64, _ := ioutil.ReadFile("../model/sample.bas64")
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

// func ImageResults(resp *httptest.ResponseRecorder) ([]*m.CADocument, error) {
// 	b, _ := ioutil.ReadAll(resp.Body)
// 	respData := ca.CaDocumentResponse{}

// 	if err := json.Unmarshal(b, &respData); err != nil {
// 		fmt.Printf("@      CADocumentResults Error: %v\n", err)
// 		return nil, err
// 	}
// 	documents := respData.Documents
// 	fmt.Printf("Test returns %d documents\n", len(documents))
// 	return documents, nil
// }
