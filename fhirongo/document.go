package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	//"github.com/tidwall/pretty"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetDocumentReference will return a result whichh has an array of document references
func (c *Connection) FindDocumentReferences(qry string) (*DocumentResults, error) {
	fmt.Printf("\n%sDocumentReference%s\n", c.BaseURL, qry)
	body, err := c.Query(fmt.Sprintf("DocumentReference/%s", qry))
	if err != nil {
		return nil, err
	}
	data := DocumentResults{}

	b := body
	//bodyStr := pretty.Pretty(b[:])
	//fmt.Printf("\n\n\n@@@ RAW DocumentReference: %s\n\n\n", bodyStr)

	//json.NewDecoder(body).Decode(&data)
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Printf("GetDocumentReference err: %s\n", err.Error())
		return nil, err
	}
	return &data, nil
}

//FindDocumentReference will return one document
func (c *Connection) GetDocumentReference(qry string) (*DocumentResults, error) {
	fmt.Printf("%sDocumentReference%s\n", c.BaseURL, qry)
	body, err := c.Query(fmt.Sprintf("DocumentReference%s", qry))
	if err != nil {
		log.Errorf("FhirDocumentReference cerner returned error: %s", err.Error())
		return nil, err
	}
	//fmt.Printf("\n\n\n@@@ RAW DocumentReference: %s\n\n\n", pretty.Pretty(b))
	data := &DocumentResults{}
	if err := json.Unmarshal(body, data); err != nil {
		log.Errorf("FhirDocumentReference Unmarshal error: %s", err.Error())
		return nil, err
	}
	fmt.Printf("FindDocumentReference:50 returning all %s\n", spew.Sdump(data))
	return data, nil
}

// Process the next page of DocRefs
func (c *Connection) NextFhirDocRefs(url string) (*DocumentResults, error) {
	//fmt.Printf("Next retrieving : %s\n", url)
	bytes, err := c.GetFhir(url)
	if err != nil {
		msg := fmt.Sprintf("NextPatient returned error: %s", err.Error())
		log.Errorf("%s", msg)
		return nil, err
	}

	data := DocumentResults{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

type DocumentResponse struct {
	//Bundle Bundle
	SearchResult
	Entry []struct {
		FullURL  string   `json:"fullUrl"`
		Document Document `json:"resource"`
	} `json:"entry"`
}
type DocumentResults struct {
	//Bundle Bundle
	SearchResult
	Entry []struct {
		FullURL  string   `json:"fullUrl"`
		Document Document `json:"resource"`
	}
}

// DocumentReference is a FHIR document
type ReferenceResults struct {
	Bundle
	Entry []struct {
		FullURL string            `json:"fullUrl"`
		DocRef  DocumentReference `json:"resource"`
	} `json:"entry"`
}

// DocumentReference is a single FHIR DocumentReference.
// Use DocumentReferences for a bundle.
type DocumentReference struct {
	CacheID           primitive.ObjectID `json:"cache_id" bson:"_id,omitempty"`
	SessionId         string             `json:"-" bson:"sessionid"`
	ResourceType      string             `json:"resourceType" bson:"resource_type"`
	ID                string             `json:"id" bson:"id"`
	FullURL           string             `json:"fullUrl" bson:"full_url"`
	EffectiveDateTime time.Time          `json:"effectiveDateTime" bson:"effective_date_time"`
	Meta              MetaData           `json:"meta" bson:"meta"`
	Text              TextData           `json:"text" bson:"text"`
	Status            Code               `json:"status" bson:"status"`
	Subject           Person             `json:"subject" bson:"subject"`
	Type              Concept            `json:"type" bson:"type"`
	Authenticator     Person             `json:"authenticator" bson:"authenticator"`
	Created           time.Time          `json:"created" bson:"created"`
	Indexed           time.Time          `json:"indexed" bson:"indexed"`
	DocStatus         Concept            `json:"docStatus" bson:"doc_status"`
	Description       string             `json:"description" bson:"description"`
	PresentedForm     []Attachment       `bson:"content" json:"presentedForm"`
	Context           EncounterContext   `bson:"context" json:"context"`
	Content           []struct {
		Attachment struct {
			ContentType string 			 `json:"contentType" bson:"content_type"`
			URL         string 			 `json:"url" bson:"url"`
			Title       string 			 `json:"title" bson:"title"`
		} 								 `json:"attachment" bson:"attachment"`
	} 									 `json:"content" bson:"content"`
	//} `json:"content"`
	//Content       []Attachment `json:"content"`
	// Context struct {
	// 	Encounter struct {
	// 		Reference string `json:"reference"`
	// 	} `json:"encounter"`
	// 	Period Period 		`json:"period"`
	// } `json:"context"`
}

type Document struct {
	CacheID           primitive.ObjectID `bson:"cache_id" json:"cacheId"`
	SessionID         string             `bson:"session_id" json:"sessionId"`
	ResourceType      string             `bson:"resource_type" json:"resourceType"`
	ID                string             `bson:"id"jbson:"id"`
	FullURL           string             `bson:"full_url" json:"fullURL"`
	EffectiveDateTime time.Time          `bson:"effective_date_time", json:"effectiveDateTime"`
	Meta              MetaData           `bson:"meta" json:"meta"`
	Text              TextData           `bson:"text" json:"text"`
	Status            string             `bson:"status" json:"status"`
	Category          CodeableConcept    `bson:"category" json:"category"`
	Code              CodeableConcept    `bson:"code" json:"code"`
	Subject           Person             `bson:"subject" json:"subject"`
	Type              Concept            `bson:"type" json:"type"`
	Encounter         EncounterReference `bson:"encounter" json:"encounter"`
	Issued            time.Time          `bson:"issued" json:"issued"`
	Performer         Person             `bson:"performer" json:"performer`
	PresentedForm     []Attachment       `bson:"presented_form" json:"presentedForm"`
	Request           Thing              `bson:"request" json:"request"`
	Result            Thing              `bson:"result" json:"result"`
	Authenticator     Person             `bson:"authenticator" json:"authenticator"`
	Created           time.Time          `bson:"created" json:"created"`
	Indexed           time.Time          `bson:"indexed" json:"indexed"`
	DocStatus         CodeableConcept    `bson:"doc_status" json:"docSatus"`
	Description       string             `bson:"description" json:"description"`
	Context           EncounterContext   `bson:"context" json:"context"`
	//Content           []Attachment       `bson:"content" json:"content"`
	Content []struct {
		Attachment struct {
			ContentType string `json:"contentType" bson:"content_type"`
			URL         string `json:"url" bson:"url"`
			Title       string `json:"title" bson:title"`
		} `json:"attachment" bson:"attachment"`
	} `json:"content" bson:"content"`
}
