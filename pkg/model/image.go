package model

// Patient allows filtering by id, mrn, names, birthdate, Encounter, andpossibly ssn
// id, mrn encounter and ssn SHOULD only return one patient and funcs are identified by Getxxx.
// The others can return an array of patients and are the functions are called Findxxx
//

import (
	"fmt"
	//"net/http"
	//"io/ioutil"
	//"encoding/json"
	//"time"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
)

//Patient contains the cached and returned information for one patient
type Image struct {
	//ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ResourceType string        `json:"resourceType"`
	Id           string        `json:"id"`
	Meta         fhir.MetaData `json:"meta"`
	ContentType  string        `json:"contentType"`
	Content      string        `json:"content"`
	URL          string        `json:"-"`
}

// Retrieve takes the URL of the caller retrieves the Image.
func RetrieveImage(docId string) (*Image, error) {
	fhirC = config.Fhir()
	fhirPdfURL := FhirPdfUrl()
	fmt.Printf("fhirPdfURL: %s\n", fhirPdfURL)
	//url := fmt.Sprintf("%s/%s",fhirPdfURL, docId)
	fImage, err := fhirC.GetImage(docId)
	if err != nil {
		return nil, err
	}
	image := Image{}
	image.Content = fImage.Content
	image.ContentType = fImage.ContentType
	image.Id = fImage.Id
	image.Meta = fImage.Meta
	image.ResourceType = fImage.ResourceType
	return &image, nil

	//url := fmt.Sprintf("%s%s", ActiveConfig().ImageURL(), i.URL)

	// fmt.Printf("Looking for Document: %s\n", docId)
	// fmt.Printf("FhirURL: %s\n", url)
	// timeout := time.Duration(60 * time.Second)
	// client := http.Client{
	// 	Timeout: timeout,
	// }
	// doc := new(DocumentSummary)
	// doc.S
	// err := FhirImageUrlForDocument()

	// req, err := http.NewRequest("GET", i.URL, nil)
	// if err != nil {
	// 	return err
	// }
	// req.Header.Add("Accept", "application/json+fhir")
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return err
	// }
	// defer resp.Body.Close()
	// b, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// if resp.StatusCode < 200 || resp.StatusCode > 299 {
	// 	err = fmt.Errorf("%d|%s", resp.StatusCode, string(b))
	// 	return err
	// }

	// if err := json.Unmarshal(b, i); err != nil {
	// 	return err
	// }
	// return nil
}
