package handler

import (
	"encoding/json"
	"fmt"

	m "github.com/dhf0820/cernerFhir/pkg/model"

	//"gopkg.in/mgo.v2/bson"
	"net/http"
	//"strconv"
	//"strings"
	//"time"
	//"github.com/davecgh/go-spew/spew"
	//"github.com/oleiade/reflections"
)

type HealthResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func WriteHealthResponse(w http.ResponseWriter, status int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	var resp HealthResponse
	resp.Status = status
	resp.Message = message

	switch status {
	case 200:
		w.WriteHeader(http.StatusOK)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
	case 403:
		w.WriteHeader(http.StatusForbidden)
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err

	}
	return nil
}

//Routes processes
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	version := fmt.Sprintf("OK: Version %s", m.Version)
	WriteHealthResponse(w, 200, version)
	fmt.Println(version)
}
