package handler

import (
	"net/http"

	m "github.com/dhf0820/cernerFhir/pkg/model"
	log "github.com/sirupsen/logrus"
)

// type ImageResponse []byte

// func WriteImageResponse(w http.ResponseWriter, status_code int, data *[]byte) error {
// 	w.Header().Set("Content-Type", "application/pdf")

// 	resp := *data

// 	w.WriteHeader(status_code)
// 	err := json.NewEncoder(w).Encode(resp)
// 	if err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return err
// 	}
// 	return nil
// }

//GetDocumentImage returns the image from the url passed in url=
func UpdateEnv(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["loglevel"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'loglevel' is missing")
		return
	}

	level := keys[0]
	log.Debugf("Level: %s\n", level)
	m.ActiveConfig().SetLogLevel(level)

	// params := mux.Vars(r)
	// fmt.Printf("Update Params: %v\n", params)
	// fmt.Printf("id: %v\n", params["id"])
	// id := params["id"]

	// fmt.Printf("Handler id: %s\n", id)

	// doc := m.DocumentSummary{EnterpriseID: id}
	// err := doc.GetDocumentImage()
	// if err != nil {
	// 	fmt.Printf("GetDocumentImage handler-433: Err: %v\n", err)
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
