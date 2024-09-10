package handler

import (
	//"bytes"
	"encoding/json"
	"net/http"

	//"time"
	l "log"
	//"os"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	//"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	//"github.com/dgrijalva/jwt-go"
	m "github.com/dhf0820/cernerFhir/pkg/model"
)

type LoginFilter struct {
	UserName string `schema:"user_name"`
	Password string `schema:"password"`
}

type LoginResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	SessionId string `json:"session_id"`
}

func WriteLoginResponse(w http.ResponseWriter, resp *LoginResponse) error {
	log.Debug("WriteLoginResponse:26")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(resp.Status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorf("Error marshaling JSON: %s", err.Error())
		return err
	}
	return nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	//var lbuf bytes.Buffer
	//var logIt = l.New(os.Stderr, "app: ", l.Ldate | l.Ltime| l.Lshortfile)
	l.SetFlags(l.Ldate | l.Ltime | l.Lshortfile)

	//l.SetFlags(l.Ldate | l.Ltime| l.Lshortfile )
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	var params LoginFilter
	err := decoder.Decode(&params, r.URL.Query())
	if err != nil {
		log.Println("Error in GET parameters : ", err)
	} else {
		//log.Println("GET parameters : ", spew.Sprint(patientFilter))
	}
	//logger := l.Ldate | l.Ltime| l.Lshortfile | l.Llongfile
	l.Println("something to log", l.LstdFlags)
	l.Println("Using standard Printf without setting flags")
	l.Println("Using standard Printf with setting flags", l.LstdFlags)

	sessionId, err := m.Login(params.UserName, params.Password)
	resp := LoginResponse{}
	if err != nil {
		resp.Status = 400
		resp.Message = err.Error()
	} else {
		resp.Status = 200
		resp.Message = "Ok"
		resp.SessionId = sessionId
		log.Debugf("AccessHandler.Login:73 -- %s", sessionId)

	}
	WriteLoginResponse(w, &resp)

	// keys, ok := r.URL.Query()["loglevel"]
	// if !ok || len(keys[0]) < 1 {
	// 	log.Println("Url Param 'loglevel' is missing")
	// 	return
	// }

	// level := keys[0]
	// log.Debugf("Level: %s\n", level)
	// m.ActiveConfig().SetLogLevel(level)

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
