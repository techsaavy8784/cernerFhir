package handler

import (
	//"encoding/json"
	//"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	//"io/ioutil"

	m "github.com/dhf0820/cernerFhir/pkg/model"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type ImageResponse []byte

func WriteImageResponse(w http.ResponseWriter, status_code int, data []byte) error {
	w.Header().Set("Content-Type", "application/pdf")

	resp := data

	w.WriteHeader(status_code)
	w.Write(resp)
	// err := json.NewEncoder(w).Encode(resp)
	// if err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// 	return err
	// }
	return nil
}

//GetDocumentImage returns the image from the url passed in url=
func GetImage(w http.ResponseWriter, r *http.Request) {
	// content, err := ioutil.ReadFile("./sample.pdf")     // the file is inside the local directory
	// fmt.Printf("Length of data: %d\n", len(content))
	// if err != nil {
	//     fmt.Println("Err")
	// 	WriteImageResponse(w, 400, nil)
	// 	return

	// }
	//WriteImageResponse(w, 200, content)
	fmt.Printf("\nAccept Header: %s\n\n", r.Header.Get("Accept"))
	params := mux.Vars(r)
	log.Debugf("\n\n\n\nGetImage Params: %v", params)
	log.Debugf("id: %s", params["id"])
	id := params["id"]
	StatusCode := 200
	ImageResponse, err := m.RetrieveImage(id)
	if err != nil {
		StatusCode = 400
		log.Errorf("GetImage:53 -- RetrieveImage error: %s", err.Error())
		WriteImageResponse(w, StatusCode, nil)
		return
	}
	dec, err := base64.StdEncoding.DecodeString(ImageResponse.Content)
	if err != nil {
		StatusCode = 400
		log.Errorf("GetImage:60 -- Base64 decoding image err: %s", err.Error())
		WriteImageResponse(w, StatusCode, nil)
		return
	}
	fmt.Printf("\n\n###ContentType: %s\n", ImageResponse.ContentType)
	fmt.Printf("\n\n###Length of Content: %d\n", len(ImageResponse.Content))
	WriteImageResponse(w, StatusCode, dec)

	// doc := m.DocumentSummary{EnterpriseID: id}
	// err := doc.GetDocumentImage()
	// if err != nil {
	// 	log.Errorf("GetImage handler-41: Err: %v\n", err)
	// 	HandleFhirError("GetDocumentImage-Handler", w, err)
	// 	return
	// }

	// pdfBytes, err := b64.StdEncoding.DecodeString(doc.Image)

	// b := bytes.NewBuffer(pdfBytes)
	// if _, err := b.WriteTo(w); err != nil {
	// 	fmt.Fprintf(w, "%s", err)
	// 	HandleFhirError("GetImage", w, err)
	// }

}
