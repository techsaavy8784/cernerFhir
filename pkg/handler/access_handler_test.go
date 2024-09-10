package handler

import (
	//"bytes"
	//http "net/http"
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	//"time"
	//"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	//fhir "github.com/dhf0820/cernerFhir/fhirongo"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	//h "github.com/dhf0820/cernerFhir/pkg/handler"
	//m "github.com/dhf0820/cernerFhir/pkg/model"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogin(t *testing.T) {
	godotenv.Load("./env_test")
	m.InitializeAll("")
	Convey("Subject: GetPatient returns the Specific Patient", t, func() {
		Convey("Given an ID: 12742397", func() {
			req := httptest.NewRequest("POST", "/api/rest/v1/login?user_name=dhf&password=password", nil)
			respData := LoginResponse{}
			resp := httptest.NewRecorder()

			NewRouter().ServeHTTP(resp, req)
			So(resp.Code, ShouldEqual, 200)

			b, _ := ioutil.ReadAll(resp.Body)

			//log.Debugf("b: %s\n", string(b))
			//defer resp.Body.Close()

			err := json.Unmarshal(b, &respData)
			So(err, ShouldBeNil)
			So(respData, ShouldNotBeNil)
			sessionId := respData.SessionId

			So(sessionId, ShouldNotEqual, "")
			session, err := m.ValidateSession(sessionId)
			//token, err :=m.VerifyTokenString(tokenString)
			So(err, ShouldBeNil)
			So(session, ShouldNotBeNil)
			//log.Debugf("Token: %s", spew.Sdump(token))
			// ad, err := m.GetTokenMetaData(token)
			// So(err, ShouldBeNil)
			// So(ad.SessionId, ShouldNotBeEmpty)
			//log.Debugf("Session: %s", ad.SessionId)
			//log.Debugf("Token: %s", tokenString)

			log.Debugf("SessionId: %s", session.SessionID)

		})
	})

}
