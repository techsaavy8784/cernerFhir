package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	m "github.com/dhf0820/cernerFhir/pkg/model"
)

//####################################### Structures #######################################
//GenericResponse struct the resultant message being returned
type GenericResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//####################################### Response Functions #######################################
func WriteGenericResponse(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// switch status {
	// case 200:
	// 	w.WriteHeader(http.StatusOK)
	// case 400:
	// 	w.WriteHeader(http.StatusBadRequest)
	// case 401:
	// 	w.WriteHeader(http.StatusUnauthorized)
	// case 403:
	// 	w.WriteHeader(http.StatusForbidden)
	// case 404:
	// 	w.WriteHeader(http.StatusNotFound)
	// case 405:
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// case 500:
	// 	w.WriteHeader(http.StatusInternalServerError)
	// case 501:
	// 	w.WriteHeader(http.StatusNotImplemented)
	// case 503:
	// 	w.WriteHeader(http.StatusServiceUnavailable)
	// }
	resp := GenericResponse{Status: status, Message: message}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Errorf("Error marshaling json: %s", err.Error())
		return err
	}

	return nil
}

//####################################### Route Handlers #######################################
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	version := fmt.Sprintf("OK: Version %s", m.Version)
	WriteGenericResponse(w, 200, version)
	fmt.Println(version)
}

// func healthcheck(w http.ResponseWriter, r *http.Request) {
// 	session := ApiSession.Copy()
// 	defer session.Close()

// 	status, message := CheckHealth(session)
// 	WriteGenericResponse(w, status, message)
// }

//func HeaderStatus(code int)

//HandleFhirError extracts the acutal error code and message from err. It send the message to
//the genericResponse Writer providing the proper code and message. The result is a usable api message
func HandleFhirError(from string, w http.ResponseWriter, err error) {
	log.Infof("FHIR Error Handler: %v\n", err)
	code, message := extractErrorDetails(err.Error())
	err = WriteGenericResponse(w, code, message)
	if err != nil {
		log.Errorf("Error writing FHIR ERROR response: %s", err)
	}
	log.Debugf("%s failed with code: %d  msg: %s", from, code, message)
	return
}

func extractErrorDetails(result string) (int, string) {
	s := strings.Split(result, "|")
	var statusCode int
	// if statusCode, err = strconv.ParseInt(s[0], 10, 64); err == nil {
	statusCode, err := strconv.Atoi(s[0])
	if err != nil {
		log.Warnf("extractErrorDetails error: %v\n", err)
		statusCode = 500
	}
	if len(s) > 1 {
		return statusCode, s[1]
	} else {
		return statusCode, result
	}
}

// func ValidateSession(from string, w http.ResponseWriter, token string) *m.AuthSession {
// 	as, err := m.ValidateAuth(token)

// 	return as
// }

func SetTokenCookie(w http.ResponseWriter, token string) *http.Cookie {
	timeout := 60 * m.LoginExpiresAfter()
	tokenCookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   int(timeout),
		HttpOnly: true,
	}
	http.SetCookie(w, tokenCookie)
	return tokenCookie
}
