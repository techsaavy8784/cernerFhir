package model

import (
	"fmt"
	"testing"

	//"github.com/davecgh/go-spew/spew"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteAuthSession(t *testing.T) {
	as := setupTest("")
	Convey("Delete AuthSession", t, func() {
		So(as.ID, ShouldNotEqual, primitive.NilObjectID)
		//session := as
		session, err := ValidateAuth("test")
		// fmt.Printf("validated session: %s\n", session.SessionID)
		// fmt.Printf("ooriginal session: %s\n", as.SessionID)
		So(err, ShouldBeNil)
		So(session, ShouldNotBeNil)
		err = session.Delete()
		So(err, ShouldBeNil)
		fmt.Printf("Validating a deleted session\n")
		_, err = ValidateAuth("test")
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual,"Not Authorized")
	})
}

func TestGetAuthSession(t *testing.T) {
	as := setupTest("")
	Convey("GetAuthSession", t, func() {
		So(as.ID, ShouldNotEqual, primitive.NilObjectID)
	
		session, err :=GetSessionForUserID(as.UserID)
		So(err, ShouldBeNil)
		So(session, ShouldNotBeNil)
		So(as.PatSessionId, ShouldEqual, session.PatSessionId)
		updAS, err := as.UpdatePatSessionId()
		So(err, ShouldBeNil)
		So(updAS.PatSessionId, ShouldNotEqual, session.PatSessionId)
	})
}

// func TestCreateSession(t *testing.T) {
// 	as := setupTest("")
// 	Convey("Delete AuthSession", t, func() {
// 		So(as.ID, ShouldNotEqual, primitive.NilObjectID)
// 		session := *as
// 		err := as.CreateSession()
// 		So(err, ShouldBeNil)
// 		s, err := ValidateAuth(as.Token )
// 		So(err, ShouldBeNil)
// 		So(s, ShouldNotBeNil)
// 		So(s.SessionID, ShouldEqual, s.SessionID)
// 		So(as, ShouldNotEqual, session.SessionID)
// 	})
// }