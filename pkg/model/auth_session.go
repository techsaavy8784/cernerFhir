package model

import (
	"context"
	"errors"
	"fmt"

	//"os"
	"strings"
	"time"

	//"github.com/dgrijalva/jwt-go"
	uuid "github.com/aidarkhanov/nanoid/v2"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	//"github.com/google/uuid"
	//"github.com/davecgh/go-spew/spew"
	"github.com/dhf0820/cernerFhir/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type Session struct {
// 	Token     string `json:"token"`
// 	CacheName string `json:"cacheName"`
// }

type AuthSession struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Status       Status             `json:"status" bson:"status"`
	PatSessionId string             `json:"pat_session_id" bson:"pat_session_id"`
	DocSessionId string             `json:"doc_session_id" bson:"doc_session_id"`
	EncSessionId string             `json:"enc_session_id" bson:"enc_session_id"`
	SessionID    string             `json:"session_id" bson:"session_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserName     string             `json:"user_name" bson:"user_name"`
	CurrentPatId string             `json:"current_pat_id" bson:"current_pat_id"` //Keeps the current patient. If changes, start a new session, Delete old
	ExpireAt     int64              `json:"expireAt" bson:"expire_at"`
	AccessedAt   *time.Time         `json:"accessed_at" bson:"accessed_at"`
}

type Status struct {
	Diagnostic string `json:"diag" bson:"diag"`
	Reference  string `json:"ref" bson:"ref"`
	Patient    string `json:"pat" bson:"pat"`
	Encounter  string `json:"enc" bson:"enc"`
}

func ValidateSession(id string) (*AuthSession, error) {
	//log.Infof("ValidateSession:49 - [%s]", id)
	if strings.Trim(id, " ") == "" {
		log.Error("auth_session:41 - session is blank")
		return nil, fmt.Errorf("401|Unauthorized")
	}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("ValidateSession:56 -- IDFromHex Error: %s", err.Error())
	}
	filter := bson.D{{"_id", ID}}

	collection, _ := storage.GetCollection("sessions")
	as := &AuthSession{}
	err = collection.FindOne(context.TODO(), filter).Decode(as)
	if err != nil {
		log.Errorf("ValidateSession:64 -- Find Session for ID %s returned ERROR: %s", id, err.Error())
		//fmt.Printf("as: %s\n", spew.Sdump(as))
		return nil, fmt.Errorf("not Authorized")
	}
	//fmt.Printf("\n\nSession as received from DB: %s\n", spew.Sdump(as))
	as.SessionID = id

	tnow := time.Now().UTC().Unix()
	//log.Infof("ValidateSession:70 - Time now: %d  expireTime: %d", tnow, as.ExpireAt)
	if tnow > as.ExpireAt {
		return nil, errors.New("NotLoggedIn")
	}

	as.UpdateSessionAccessedAt()
	//fmt.Printf("\n\nValidate is returning session: %s\n", spew.Sdump(as))
	return as, nil
}

//func (as *AuthSession) NewSessionID() error {
// 	id, err := uuid.New()
// 	if err!= nil {
// 		return fmt.Errorf("Cound not generate uuid: %s\n", err.Error())
// 	}
// 	err = DeleteAllCasheForSession(as.SessionID)
// 	as.SessionID = id
// 	filter := bson.M{"_id": as.ID}

// 	update := bson.M{"$set": bson.M{"session_id": as.SessionID}}
// 	collection, _ := storage.GetCollection("sessions")
// 	_, err = collection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		msg := fmt.Sprintf("Update SessionID failed: %s", err.Error)
// 		log.Error(msg)
// 		return errors.New(msg)
// 	}

// }

func (as *AuthSession) CreateSession(userId primitive.ObjectID) error { // SessionID is provided
	//as.SessionID = uuid.New()
	if as.SessionID != "" {
		as.Delete()
	}

	id, err := uuid.New()
	if err != nil {
		return fmt.Errorf("auth_session:95 -- Could not generate uuid: %s\n", err.Error())
	}
	//fmt.Printf("CreateSession:100 -- cheking if session exists: %s\n", spew.Sdump(as))
	// as, err = ValidateAuth(as.Token)
	// if err == nil {
	// 	log.Infof("Session already exists for %s\n", as.Token)

	// 	as.UpdateSessionID()
	// 	return nil //errors.New("Session already exsts")
	// } else {
	// 	msg := fmt.Sprintf("auth_session:77 -- err: %s", err.Error())
	// 	log.Error(msg)
	// 	return errors.New(msg)
	// }
	// if as == nil {
	// 	log.Errorf("auth_session:76 -- as is nil returned from")
	// }
	as.UserID = userId
	as.SessionID = id
	//log.Infof("Creating Session: %s\n", spew.Sdump(as))
	err = as.Insert()
	if err != nil {
		return fmt.Errorf("auth_session:117 -- Insert failed: %s", err.Error())
	}
	// filter := bson.D{{"token", as.Token}}
	// collection, _ := storage.GetCollection("sessions")

	// err = collection.FindOne(context.TODO(), filter).Decode(&as)
	// if err != nil {
	// 	fmt.Printf("Create:82 - FindFilter: %s - Err:%s\n", as.Token, err.Error())
	// }
	//fmt.Printf("Right after Insert: %s\n", spew.Sdump(as))
	return nil
}

func (as *AuthSession) Delete() error {
	startTime := time.Now()
	collection, _ := storage.GetCollection("sessions")
	filter := bson.D{{"sessionid", as.SessionID}}
	//log.Debugf("    bson filter delete: %v\n", filter)
	deleteResult, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Errorf("!     137 -- DeleteSession for Dession %s failed: %v", as.SessionID, err)
		return err
	}
	log.Infof("@@@!!!   140 -- Deleted %d Sessions for session: %v in %s", deleteResult.DeletedCount, as.SessionID, time.Since(startTime))
	return nil
}

func CreateSessionForUser(user *User) (*AuthSession, error) {
	userID := user.ID

	filter := bson.D{{"user_id", userID}}
	collection, _ := storage.GetCollection("sessions")
	as := &AuthSession{}
	err := collection.FindOne(context.TODO(), filter).Decode(as) // See if the user already has a session
	if err == nil {                                              // The user has a session, keep using it
		as.UpdateSessionAccessedAt() // Extend the current session
		return as, nil
	}
	// Create a new Session
	as.UserID = userID
	as.UserName = user.UserName
	err = as.Insert()
	if err != nil {
		msg := fmt.Sprintf("insert Session error: %s", err.Error())
		log.Error(msg)
		return nil, errors.New(msg)
	}
	return as, nil
}

func GetSessionForUserID(userID primitive.ObjectID) (*AuthSession, error) {
	filter := bson.D{{"user_id", userID}}
	collection, _ := storage.GetCollection("sessions")
	as := &AuthSession{}
	err := collection.FindOne(context.TODO(), filter).Decode(as) // See if the user already has a session
	return as, err
}

func ValidateAuth(token string) (*AuthSession, error) {
	//TODO: Actually validate the session as a valid JWT. Right now just using

	if strings.Trim(token, " ") == "" {
		log.Error("auth_session:187 - token is blank")
		return nil, fmt.Errorf("401|Unauthorized")
	}
	oId, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		msg := "ValidateAuth192 -- Invalid ID"
		log.Error(msg)
		return nil, errors.New(msg)
	}
	filter := bson.D{{"_id", oId}}

	collection, _ := storage.GetCollection("sessions")
	var as AuthSession
	err = collection.FindOne(context.TODO(), filter).Decode(&as)
	if err != nil {
		log.Debugf("auth_session:202 -- Find Session for token %s returned ERROR: %v", token, err)
		//fmt.Printf("as: %s\n", spew.Sdump(as))
		return nil, fmt.Errorf("not Authorized")
	}
	tnow := time.Now().UTC().Unix()
	if tnow > as.ExpireAt {
		return nil, errors.New("notLoggedIn")
	}
	//log.Debugf("auth_session210 -- Validate Found Session: %s", spew.Sdump(as))
	//log.Debugf("Found Session: %s update ExpireAt: %s  \n\n", as.ID.String(), as.ExpireAt.String())
	//spew.Dump(as)
	as.UpdateSessionAccessedAt()
	return &as, nil
}

func (as *AuthSession) Insert() error {

	// fmt.Printf("On entry into Insert:%s\n", spew.Sdump(as))

	as.ExpireAt = as.CalculateExpireTime()
	tn := time.Now().UTC()
	as.AccessedAt = &tn
	collection, _ := storage.GetCollection("sessions")

	insertResult, err := collection.InsertOne(context.TODO(), as)
	if err == nil {
		as.ID = insertResult.InsertedID.(primitive.ObjectID)
		//log.Debugf("   Set ID: %s\n", as.ID.String())
	} else {
		log.Errorf("Insert:231 -- Error: %s\n", err.Error())
		return err
	}

	return nil
}

func (as *AuthSession) UpdateSessionAccessedAt() {
	//fmt.Printf("AuthSession.UpdateSessionAccessedAt:242 -- %s\n", spew.Sdump(as))
	filter := bson.M{"_id": as.ID}
	//oldTime := as.ExpireAt
	as.ExpireAt = as.CalculateExpireTime()
	//log.Infof("UpdateSessionAccessAt:242 - OldExpireTime = %d   New ExpireTime = %d", oldTime, as.ExpireAt)
	tn := time.Now().UTC()
	as.AccessedAt = &tn
	update := bson.M{"$set": bson.M{"expire_at": as.ExpireAt, "accessed_at": as.AccessedAt}}

	collection, _ := storage.GetCollection("sessions")
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Errorf("auth_session:248 -- Update error ignored: %s\n\n", err)
	}
	asUpd, _ := GetSessionForUserID(as.UserID)
	as.AccessedAt = asUpd.AccessedAt
	as.ExpireAt = asUpd.ExpireAt
	//as = asUpd
	//log.Debugf("Matched: %d  -- modified: %d for ID: %s\n", res.MatchedCount, res.ModifiedCount, as.ID.String())
}

func (as *AuthSession) Update(update bson.M) (*AuthSession, error) {

	//fmt.Printf("AuthSession.Update: 274 -- as: %s\n", spew.Sdump(as))
	collection, _ := storage.GetCollection("sessions")
	//fmt.Printf("LIne 258\n")
	filter := bson.M{"_id": as.ID}
	//fmt.Printf("Filter: %v\n", filter)
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Errorf("AuthSession.Update:272 error %s", err)
		return nil, err
	}
	//log.Debugf("AuthSession.Update:275 -- Matched: %d  -- modified: %d for ID: %s", res.MatchedCount, res.ModifiedCount, as.ID.String())

	asUpd, err := GetSessionForUserID(as.UserID)
	as = asUpd
	return asUpd, err
}

func (as *AuthSession) UpdateDiagStatus(status string) (*AuthSession, error) {
	fmt.Printf("AuthSession.UpdateDiagStatus:292\n")
	as.Status.Diagnostic = status
	update := bson.M{"$set": bson.M{"status": as.Status}}
	asUpd, err := as.Update(update)
	if err != nil {
		err = fmt.Errorf("UpdateStatus:294 -- error: %s", err.Error())
		log.Error(err.Error())
		return nil, err
	}
	return asUpd, nil
}

func (as *AuthSession) UpdatePatStatus(status string) (*AuthSession, error) {
	fmt.Printf("AuthSession.UpdatePatStatus:302")
	as.Status.Patient = status
	update := bson.M{"$set": bson.M{"status": as.Status}}
	asUpd, err := as.Update(update)
	if err != nil {
		err = fmt.Errorf("UpdateStatus:294 -- error: %s", err.Error())
		log.Error(err.Error())
		return nil, err
	}
	return asUpd, nil
}

func (as *AuthSession) UpdateRefStatus(status string) (*AuthSession, error) {
	//fmt.Printf("AuthSession.UpdateStatus:316")
	as.Status.Reference = status
	update := bson.M{"$set": bson.M{"status": as.Status}}
	asUpd, err := as.Update(update)
	if err != nil {
		err = fmt.Errorf("UpdateStatus:322 -- error: %s", err.Error())
		log.Error(err.Error())
		return nil, err
	}
	return asUpd, nil
}

func (as *AuthSession) UpdateEncStatus(status string) (*AuthSession, error) {
	fmt.Printf("AuthSession.UpdateEncStatus:329")

	as.Status.Encounter = status
	update := bson.M{"$set": bson.M{"status": as.Status}}
	asUpd, err := as.Update(update)
	if err != nil {
		err = fmt.Errorf("UpdateStatus:335 -- error: %s", err.Error())
		log.Error(err.Error())
		return nil, err
	}
	return asUpd, nil
}

func (as *AuthSession) UpdateEncSessionId() (*AuthSession, error) {
	fmt.Printf("AuthSession.UpEncSessionId:348 --Entry: %s\n", spew.Sdump(as))

	id, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf("AuthSession.UpdateEncSessionId:352 -- Could not generate Enc uuid: %s", err.Error())
	}
	update := bson.M{"$set": bson.M{"enc_session_id": id}}
	if err != nil {
		return nil, fmt.Errorf("AuthSession.UpdatEncSessionId:291 -- Cound not set EncSessionID uuid: %s", err.Error())
	}
	fmt.Printf("AuthSession.UpdatEncSessionId:293 -- %s\n", spew.Sdump(as))
	asUpd, err := as.Update(update)
	as = asUpd
	return asUpd, err
}

func (as *AuthSession) UpdatePatSessionId() (*AuthSession, error) {
	fmt.Printf("AuthSession.UpdatePatSessionId:366 --Entry: %s\n", spew.Sdump(as))

	id, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf("AuthSession.UpdatePatSessionId:287 -- Cound not generate Pat uuid: %s", err.Error())
	}
	update := bson.M{"$set": bson.M{"pat_session_id": id}}
	if err != nil {
		return nil, fmt.Errorf("AuthSession.UpdatePatSessionId:291 -- Cound not set PatSessionID uuid: %s", err.Error())
	}
	fmt.Printf("AuthSession.UpdatePatSessionId:293 -- %s\n", spew.Sdump(as))
	asUpd, err := as.Update(update)
	as = asUpd
	return asUpd, err
}

func (as *AuthSession) UpdateDocSessionId() (*AuthSession, error) {
	fmt.Printf("AuthSession.UpdateDocSessionId:383 --Entry: %s\n", spew.Sdump(as))
	id, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf("auth_session.UpdateDocId:302 -- Cound not generate Doc uuid: %s", err.Error())
	}
	update := bson.M{"$set": bson.M{"doc_session_id": id}}
	if err != nil {
		return nil, fmt.Errorf("AuthSession.UpdateDocSessionId:306 -- Cound not set DocSessionID uuid: %s", err.Error())
	}
	return as.Update(update)
}

func (as *AuthSession) UpdateSessionID() (*AuthSession, error) {
	id, err := uuid.New()
	if err != nil {
		return nil, fmt.Errorf("AuthSession.UpdateSessionId:314 -- Cound not generate uuid: %s", err.Error())
	}
	update := bson.M{"$set": bson.M{"session_id": id}}
	return as.Update(update)

	// collection, _ := storage.GetCollection("sessions")
	// res, err := collection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	// 	log.Errorf(" Update error %s", err)
	// 	return err
	// }
	// log.Debugf("auth_session:265 -- Matched: %d  -- modified: %d for ID: %s", res.MatchedCount, res.ModifiedCount, as.ID.String())
	//return nil
}

func (as *AuthSession) CalculateExpireTime() int64 {
	loc, _ := time.LoadLocation("UTC")
	addlTime := time.Hour * 2
	ExpireAt := time.Now().In(loc).Add(addlTime).Unix()
	return ExpireAt
}

func (as *AuthSession) GetDocumentStatus() string {
	latest, _ := GetSessionForUserID(as.UserID)
	if latest.Status.Diagnostic == "filling" || latest.Status.Reference == "filling" {
		return "filling"
	}
	return "done"
}

func (as *AuthSession) GetDiagReptStatus() string {
	latest, _ := GetSessionForUserID(as.UserID)
	return latest.Status.Diagnostic
}

func (as *AuthSession) GeReptRefStatus() string {
	latest, _ := GetSessionForUserID(as.UserID)
	return latest.Status.Reference
}
func (as *AuthSession) GetPatientStatus() string {
	latest, _ := GetSessionForUserID(as.UserID)
	return latest.Status.Patient
}

func (as *AuthSession) GetEncounterStatus() string {
	latest, _ := GetSessionForUserID(as.UserID)
	return latest.Status.Encounter
}
