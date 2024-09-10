package handler

import (
	//"fmt"
	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	//"github.com/davecgh/go-spew/spew"
	m "github.com/dhf0820/cernerFhir/pkg/model"
	"github.com/joho/godotenv"

	"os"
)

//TODO: Not sure we watn to initialize everything for Testing. No reason to start listening
func setupTest(token string) *m.AuthSession {
	godotenv.Load("env_test")
	m.InitializeAll(os.Getenv("MONGODB"))
	asId, err := m.Login("test", "password")
	if err != nil {
		log.Panicf("Could not log in")
	}

	as, err := m.ValidateSession(asId)
	if err != nil {
		log.Panicf("Could not validate new session")
	}
	//fmt.Printf("SetupSession:26 %s\n", spew.Sdump(as))

	return as
	// as := m.AuthSession{}
	// if token == "" {
	// 	token = "handlers"
	// }
	// as.Token = token
	// //as.SessionID = "test"
	// err := as.CreateSession()
	// if err != nil {
	// 	fmt.Printf("CreateSession returned err: %s\n", err.Error())
	// 	return nil
	// }
	// fmt.Printf("Created Session:%s\n", spew.Sdump(as))
	// session, err :=m. ValidateAuth(token)
	// if err != nil {
	// 	fmt.Printf("Create:29 ValidateAuth err: %s\n", err.Error())
	// }
	// return session
}

// func setupAD() *m.AccessDetails {
// 	godotenv.Load("env_test")
// 	m.InitializeAll("")
// 	tokenStr, err := m.Login("dhf", "password")
// 	if err != nil {
// 		log.Panicf("Login failed")
// 	}
// 	token, err := m.VerifyTokenString(tokenStr)
// 	if err != nil {
// 		log.Panic("TestToken did not verify")
// 	}
// 	err = m.TokenValid(token)
// 	if err!= nil {
// 		log.Panic("TestToken is invalid")
// 	}
// 	ad, err := m.GetTokenMetaData(token)
// 	if err != nil {
// 		log.Panic("TestGetMetaData error")
// 	}
// 	//ad.TokenStr = tokenStr

// 	return ad
// }
