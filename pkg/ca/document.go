package ca

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	//"github.com/davecgh/go-spew/spew"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CADocument struct {
	CacheID      primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	SessionId    string             `json:"-" bson:"session_id"`
	ID           string             `json:"document_id" bson:"document_id"`
	VersionID    uint64             `json:"version_id" bson:"version_id"`
	Encounter    string             `json:"visit_num" bson:"visit_num"`
	Repository   string             `json:"repository" bson:"repository"`
	Category     string             `json:"category" bson:"category"`
	Class        string             `json:"class" bson:"class"`
	Source       string             `json:"source" bson:"source"`
	Description  string             `json:"description" bson:"description"`
	ImageURL     string             `json:"image_url" bson:"image_url"`
	Pages        int                `json:"pages" bson:"pages"`
	ReptDateTime *time.Time         `json:"rept_datetime" bson:"rept_datetime"`
	Subtitle     string             `json:"subtitle" bson:"subtitle"`
	PatientGPI   string             `json:"patient_gpi" bson:"patient_gpi"`
	Text         string             `json:"text" bson:"text"`
	Type         string             `json:"type" bson:"type"`
	DocStatus    string             `json:"doc_status" bson:"doc_status"`
	CreatedAt    *time.Time         `json:"-" bson:"created_at"`
	UpdatedAt    *time.Time         `json:"-" bson:"updated_at"`
	AccessedAt   *time.Time         `json:"-" bson:"accessed_at"`
	//Text         fhir.TextData `json:"text"`
}

type CaDocumentResponse struct {
	StatusCode   int           `json:"status_code"`
	Message      string        `json:"message"`
	CacheStatus  string        `json:"cache_status"`
	TotalInCache int64         `json:"totaldocs"`
	PagesInCache int64         `json:"pages_in_cache"`
	NumberInPage int64         `json:"docs_in_page"`
	Page         int64         `json:"page"`
	SessionId    string        `json:"session_id"`
	Documents    []*CADocument `json:"documents"`
	Document     *CADocument   `json:"document"`
}

// type DocContent       []struct {
// 	Attachment struct {
// 		ContentType string `json:"contentType" bson:"content_type"`
// 		URL         string `json:"url"`
// 	} `json:"attachment"`
// }

func WriteCaDocumentResponse(w http.ResponseWriter, resp *CaDocumentResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}
	return nil
}

func FhirDocumentsToCA(w http.ResponseWriter, total, pages, inPage, page int64, cacheStatus string, fDocs []*fhir.Document) {
	fmt.Printf("\n################################ FhirDocumentsToCa ##############################################\n")
	caDocuments := FhirDocsToCa(fDocs)
	resp := CaDocumentResponse{}
	resp.StatusCode = 200
	resp.Message = "Ok"
	resp.CacheStatus = cacheStatus
	resp.TotalInCache = total
	resp.PagesInCache = pages
	resp.NumberInPage = inPage
	resp.Page = page
	resp.Documents = caDocuments
	WriteCaDocumentResponse(w, &resp)
}

func FhirDocsToCa(fDocs []*fhir.Document) []*CADocument {
	caDocuments := []*CADocument{}
	for _, d := range fDocs {
		doc := FhirDocumentToCA(d)
		caDocuments = append(caDocuments, doc)
	}
	return caDocuments
}

// func FhirDocumentsToCA(fDocs []*fhir.Document) []*CADocument{
// 	caDocuments := []*CADocument{}
// 	for _, d := range fDocs {
// 		doc := FhirDocumentToCA(d)
// 		caDocuments = append(caDocuments, doc)
// 	}
// 	return caDocuments
// }

func FhirDocumentToCA(fd *fhir.Document) *CADocument {
	//fmt.Printf("fhir.Doc: %s\n", spew.Sdump(fd))
	caDoc := CADocument{}
	config := m.ActiveConfig()

	caDoc.PatientGPI = ExtractFhirPatientID(fd.Subject.Reference)
	if fd.Encounter.Reference != "" {
		caDoc.Encounter = ExtractFhirEncounterNum(fd.Encounter)
	} else {
		caDoc.Encounter = ExtractFhirEncounterNum((fd.Context.Encounter))
	}
	dt := fd.Created
	caDoc.ID = fd.ID
	caDoc.ReptDateTime = &dt
	caDoc.Class = fd.Type.Text
	caDoc.Subtitle = fd.Code.Text
	caDoc.Text = fd.Text.Div
	caDoc.Type = fd.Type.Text
	caDoc.Repository = config.Source()
	//fmt.Printf("DocStatus: %s\n", spew.Sdump(fd.DocStatus))
	caDoc.DocStatus = fd.DocStatus.Text
	//caDoc.VersionID = fd.Meta.VersionID
	caDoc.Source = config.Source() // The fhir host.

	switch fd.ResourceType {
	case "DiagnosticReport":
		caDoc.Category = fd.ResourceType
		caDoc.Class = fd.Code.Text
		//fmt.Printf("DiagnosticReport: id: %s  - %s\n", fd.ID, spew.Sdump(fd.PresentedForm))
		caDoc.ImageURL = ExtractDiagnosticImage(fd.PresentedForm, "application/pdf")
		description := ExtractDiagnosticTitle(fd.PresentedForm, "application/pdf")
		if description != "" {
			caDoc.Description = description
		} else {
			caDoc.Description = caDoc.Class
		}
	case "DocumentReference":
		caDoc.Class = fd.Type.Text
		caDoc.Category = fd.ResourceType
		//fmt.Printf("ReportReference: id: %s  -- %s\n", fd.ID, spew.Sdump(fd.Content))
		caDoc.ImageURL = ExtractAttachmentImage(fd, "application/pdf")
		description := ExtractAttachmentText(fd, "application/pdf")
		if description != "" {
			caDoc.Description = description
		} else {
			caDoc.Description = caDoc.Class
		}
	}
	fmt.Printf("\n\n")
	return &caDoc
}

func ExtractDiagnosticImage(forms []fhir.Attachment, imageType string) string {
	for _, att := range forms {
		if att.ContentType == imageType {
			//fmt.Printf("Returning Diagnostic Image: %s\n", att.URL)
			return att.URL
		}
	}
	return ""
}

func ExtractDiagnosticTitle(forms []fhir.Attachment, imageType string) string {
	for _, att := range forms {
		if att.ContentType == imageType {
			return att.Title
		}
	}
	return ""
}

func ExtractAttachmentImage(fd *fhir.Document, imageType string) string {
	for _, ctn := range fd.Content {
		if ctn.Attachment.ContentType == imageType {
			//fmt.Printf("Returning image url: %s\n", ctn.Attachment.URL)
			return ctn.Attachment.URL
		}
	}
	return ""
}

func ExtractAttachmentText(fd *fhir.Document, imageType string) string {
	for _, ctn := range fd.Content {
		if ctn.Attachment.ContentType == imageType {
			return ctn.Attachment.Title
		}
	}
	return ""
}

func ExtractFhirPatientID(pat string) string {
	p := strings.Split(pat, "/")
	if len(p) > 1 {
		return p[1]
	}
	return ""
}

func ExtractFhirEncounterNum(e fhir.EncounterReference) string {
	p := strings.Split(e.Reference, "/")
	if len(p) > 1 {
		return p[1]
	}
	return ""
}

// func FhirDiagDocsToCA(fds []*fhir.DiagnosticReport) []*CADocument{
// 	caDocuments := []*CADocument{}
// 	for _, d := range fds {
// 		doc := FhirDiagDocToCA(d)
// 		caDocuments = append(caDocuments, doc)
// 	}
// 	return caDocuments
// }

// func FhirDiagDocToCA(fd* fhir.DiagnosticReport) *CADocument {
// 	var caDoc CADocument
// 	var err error
// 	caDoc.PatientGPI = fd.ID

// 	caDoc.PatientGPI = fd.ID
// 	caDoc.VersionID, err = strconv.ParseUint(fd.Meta.VersionID, 10, 64)
// 	if err != nil {
// 		log.Errorf("Invalid VersionID: [%s] error:%s\n", fd.Meta.VersionID, err.Error())
// 	}
// 	enc := strings.Split(fd.Encounter.Reference, "/")
// 	if len(enc) > 1 {
// 		caDoc.Encounter =enc[1]
// 	}			// Encounter is not available
// 	caDoc.Repository = "FHIR"
// 	rpdt := fd.EffectiveDateTime
// 	caDoc.ReptDateTime = &rpdt
// 	caDoc.ImageURL = GetImage(fd, "application/pdf")
// 	caDoc.Pages = 0            // Unavailable
// 	caDoc.Subtitle = fd.Code.Text
// 	caDoc.Text = fd.Text
// 	cfg := ActiveConfig()
// 	caDoc.Source = strings.ToLower(cfg.Source())

// 	if caDoc.Description == "" {
// 		caDoc.Description = fd.Code.Text
// 	}
// 	// config := ActiveConfig()

// 	// switch source {
// 	// case "cerner":
// 	// 	caDoc.ImageURL = fmt.Sprintf("%s%s", config.ImageURL(), caDoc.ID)
// 	// case "ca":
// 	// 	caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
// 	// case "QC":
// 	// 	caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
// 	// case "HPF":
// 	// 	caDoc.ImageURL = fmt.Sprintf("%s/%d", config.Env("caImageURL"), caDoc.VersionID)
// 	// }
// 	return &caDoc
// }
