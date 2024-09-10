package model

import (
	"testing"
	//"time"

	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateToken(t *testing.T) {
	Convey("CreateToken ", t, func() {
		token, err := CreateToken("dhfId")
		So(err, ShouldBeNil)
		So(token, ShouldNotBeNil)
		//log.Printf("Token returned: %s", spew.Sdump(token))

		// accessData, err := GetTokenMetaData(token)
		// So(err, ShouldBeNil)
		// log.Printf("AccessData: %s", spew.Sdump(accessData))
		// So(accessData, ShouldNotBeNil)
		// expiresAt := accessData.ExpiresAt
		// at := SetTokenSession(token, "ABCD")
		// time.Sleep(time.Second)
		// exp_token := UpdateTokenExpire(at)
		// log.Printf("Updated Token returned: %s", spew.Sdump(exp_token))
	
		// ad, err := GetTokenMetaData(exp_token)
		// So(err, ShouldBeNil)
		// So(ad, ShouldNotBeNil)
		// So(ad.ExpiresAt, ShouldNotEqual, expiresAt)
		// So(ad.SessionId, ShouldEqual, "ABCD")

	})
}

func TestLogin(t *testing.T) {
	godotenv.Load("env_test")
	InitializeAll("")
	Convey("CreateLogin ", t, func() {
		token, err := Login("dhf","password")
		So(err, ShouldBeNil)
		So(token, ShouldNotBeNil)
		as, err := ValidateAuth(token)
		So(err, ShouldBeNil)
		So(as, ShouldNotBeNil)
		patSessionId := as.PatSessionId
		updAs, err := as.UpdatePatSessionId()
		So(err, ShouldBeNil)
		So(updAs.PatSessionId, ShouldNotEqual, patSessionId)
		log.Printf("as: %s", spew.Sdump(updAs))

		// err = TokenValid(token)
		// So(err, ShouldBeNil)
		// accessData, err := GetTokenMetaData(token)
		// So(err, ShouldBeNil)
		// log.Printf("AccessData: %s", spew.Sdump(accessData))
		// So(accessData, ShouldNotBeNil)
	})
}

