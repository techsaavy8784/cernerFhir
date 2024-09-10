package handler

import (
	"net/http"
	//_ "github.com/dhf0820/cernerFhir/docs"
	//h "github.com/dhf0820/cernerFhir/pkg/handler"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{

	Route{
		"HealthCheck",
		"GET",
		"/api/v1/healthcheck",
		HealthCheck,
	},
	Route{
		"Login",
		"POST",
		"/api/rest/v1/login",
		Login,
	},
	Route{
		"FindPatient",
		"GET",
		"/api/rest/v1/patient",
		SearchCaPatient,
	},
	Route{
		"GetPatient",
		"GET",
		"/api/rest/v1/patient/{id}",
		GetPatient,
	},
	Route{
		"SessionPatient",
		"GET",
		"/api/rest/v1/session/{id}/patient",
		SessionPatients,
	},
	Route{
		"SessionDocument",
		"GET",
		"/api/rest/v1/session/{id}/document",
		SessionDocuments,
	},
	Route{
		"SessionEncounter",
		"GET",
		"/api/rest/v1/session{id}/encounter",
		SessionEncounters,
	},
	Route{
		"QueryDocuments",
		"GET",
		"/api/rest/v1/documents",
		QueryDocuments,
	},
	// Route{
	// 	"GetEncounter",
	// 	"GET",
	// 	"/api/rest/v1/encounter/{id}",
	// 	GetEncounter,
	// },
	Route{
		"FindEncounters",
		"GET",
		"/api/rest/v1/encounters",
		FindEncounters,
	},
	Route{
		"GetDocumentImage",
		"GET",
		"/api/rest/v1/image/{id}",
		GetImage,
	},
	Route{
		"AddDocumentsToEMR",
		"POST",
		"/api/rest/v1/documents/emr_save",
		AddEmrDocuments,
	},
	Route{
		"AddPatientToEMR",
		"POST",
		"/api/rest/v1/patient/emr_save",
		AddPatientEMR,
	},
	// Route{
	// 	"PatientDocuments",
	// 	"GET",
	// 	"/api/v1/rest/patient/{id}/documents",
	// 	PatientDocuments,
	// },
	Route{
		"Admin",
		"PUT",
		"/api/rest/v1/admin",
		UpdateEnv,
	},
	// Route{
	// 	"Swagger",
	// 	"GET",
	// 	"/api/v1/swagger/*",
	// 	doc.Handler,
	// },
	// Route{
	// 	"GetDocumentImage",
	// 	"GET",
	// 	"/api/v1/document_image",
	// 	h.GetDocumentImage,
	// },
	// Route{
	// 	"GetDocument",
	// 	"GET",
	// 	"/api/v1/document/{id}",
	// 	h.GetDocument,
	// },
	// Route{
	// 	"FindPatientDocuments",
	// 	"GET",
	// 	"/api/v1/patient_documents",
	// 	h.FindPatientDocuments,
	// },
	// Route{
	// 	"GetDocumentImage",
	// 	"GET",
	// 	"/api/v1/document_image",
	// 	h.GetDocumentImage,
	// },
	// Route{
	// 	"FindEncounters",
	// 	"GET",
	// 	"/api/v1/encounters",
	// 	h.FindEncounters,
	// },
	// Route{
	// 	"GetEncounter",
	// 	"GET",
	// 	"/api/v1/encounter/{id}",
	// 	h.GetEncounter,
	// },
	// Route{
	// 	"GetPatientEncounters",
	// 	"GET",
	// 	"/api/v1/patient_encounters/{id}",
	// 	h.FindPatientEncounters,
	// },
}
