package model

import (
	//"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	//m "github.com/dhf0820/cernerFhir/pkg/model"
	//"os"
)

//TODO: Not sure we watn to initialize everything for Testing. No reason to start listening
func setupTest(token string) *AuthSession {
	godotenv.Load("env_test")

	InitializeAll("")
	asId, err := Login("test", "password")
	if err != nil {
		log.Panicf("Could not log in")
	}

	// as := AuthSession{}
	// fmt.Printf("\nCreating Session\n")
	// if token == "" {
	// 	as.Token = "go-test"
	// } else {
	// 	as.Token = token
	// }
	// err := as.CreateSession()
	// if err != nil {
	// 	fmt.Printf("CreateSession Failed: %s\n", err.Error())
	// }
	as, err := ValidateSession(asId)
	if err != nil {
		log.Panicf("Could not validate new session")
	}
	return as
	// session, err := ValidateAuth("go-test")
	// if err != nil {
	// 	fmt.Printf("%s\n", err.Error())
	// }
	// return session
}

// func setupAD() *AccessDetails {
// 	godotenv.Load("env_test")
// 	InitializeAll("")
// 	tokenStr, err := Login("dhf", "password")
// 	if err != nil {
// 		log.Panicf("Login failed")
// 	}
// 	token, err := VerifyTokenString(tokenStr)
// 	if err != nil {
// 		log.Panic("TestToken did not verify")
// 	}
// 	err = TokenValid(token)
// 	if err!= nil {
// 		log.Panic("TestToken is invalid")
// 	}
// 	ad, err := GetTokenMetaData(token)
// 	if err != nil {
// 		log.Panic("TestGetMetaData error")
// 	}
// 	//ad.TokenStr = tokenStr

// 	return ad
// }
