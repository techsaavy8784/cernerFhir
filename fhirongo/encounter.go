package fhirongo

import (
	"encoding/json"
	"fmt"
	"time"
	//"strings"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	log "github.com/sirupsen/logrus"

)

// GetEncounter will return an Encounter for a number (Encounter)
func (c *Connection) GetEncounter(eid string) (*Encounter, error) {
	fmt.Printf("%sEncounter/%s\n", c.BaseURL, eid)
	res, err := c.Query(fmt.Sprintf("Encounter/%s", eid))
	if err != nil {
		return nil, err
	}
	data := Encounter{}
	if err := json.Unmarshal(res, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Connection) NextFhirEncounters(url string) (*EncounterResults, error) {
	//fmt.Printf("Next retrieving : %s\n", url)
	bytes, err := c.GetFhir(url)
	if err != nil {
		msg := fmt.Sprintf("NextFhirEncounter returned error: %s", err.Error())
		log.Errorf("%s", msg)
		return nil, err
	}

	data := EncounterResults{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	return &data, nil
}


// GetEncounter will return Encounters for a patient with id pid
func (c *Connection) GetPatientEncounters(pid string) (*EncounterResults, error) {
	fmt.Printf("%sEncounter?patient=%s\n", c.BaseURL, pid)
	res, err := c.Query(fmt.Sprintf("Encounter?patient=%s", pid))

	if err != nil {
		fmt.Printf("Encounter Query, %s\n", err)
		return nil, err
	}
	//fmt.Printf("\n\n\n@@@ RAW Encounter: %s\n\n\n", pretty.Pretty(b))
	data := EncounterResults{}
	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	// fmt.Printf("\n\n\nUnmarshalled:")
	//spew.Dump(data)
	return &data, nil
}

// GetEncounters will return Encounters for a patient with id pid
func (c *Connection) FindFhirEncounters(query string) (*EncounterResults, error) {
	fmt.Printf("%sEncounter?%s\n", c.BaseURL, query)
	res, err := c.Query(fmt.Sprintf("Encounter?%s", query))
	if err != nil {

		return nil, err
	}
	//fmt.Printf("\n\n\n@@@ RAW Encounter: %s\n\n\n", pretty.Pretty(b))
	data := EncounterResults{}
	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	// fmt.Printf("\n\n\nUnmarshalled:")
	// spew.Dump(data)
	return &data, nil
}

type EncounterResults struct {
	SearchResult
	Entry []struct {
		FullUrl		string		`json:"full_url"`
		//EntryPartial
		Encounter				`json:"resource"`
	} 		`json:"entry"`
}


type PrelimEncounter struct {
	CacheID 			primitive.ObjectID 	`json:"-" bson:"_id"`
	SessionId			string 				`json:"-" bson:"session_id"`
	ResourceType      	string       		`json:"resourceType" bson:"resource_type"`
	ID               	string       		`json:"id" bson:"id"`
	EffectiveDateTime 	time.Time    		`json:"effectiveDateTime" bson:"effective_date_time"`
	RecordedDate      	time.Time    		`json:"recordedDate" bson:"recorded_date"`

	Meta 				MetaData			`json:"meta" bson:"meta"`
	Text              	TextData     		`json:"text" bson:"text"`
	Identifiers       	[]Identifier 		`json:"identifier" bson:"identifier"`
	Status            	string       		`json:"status" bson:"status"`
	Class             	string       		`json:"class" bson:"class"`
	Type              	[]Concept    		`json:"type" bson:"type"`
	Priority			*CodeableConcept			`json:"priority"`
	Subject           	Person       		`json:"subject" bson:"subject"`
	Patient           	Person       		`json:"patient" bson:"patient"`
	Performer         	Person       		`json:"performer" bson:"performer"`
	Recorder          	Person       		`json:"recorder" bson:"recorder"`
	Code              	Code         		`json:"code" bson:"code"`
	Category          	Code         		`json:"category" bson:"category"`
	Reasons           	[]*Reason    		`json:"reason" bson:"reason"`
	Description       	string       		`json:"description" bson:"description"`
	Participant    		[]EncounterParticipant	`json:"participant,omitempty" bson:"participant,omitempty"`
//		Participant     	[]BackboneElement
	Period				Period		   		`json:"period" bson:"period"`
	Location        	[]Location			`json:"location" bson:"location"`
	ServiceProvider 	ServiceProvider		`json:"service_provider" bson:"service_provider"`
}

type Encounter struct {
	CacheID 		  primitive.ObjectID 		`bson:"_id" json:"-"`
	SessionId		  string 					`bson:"session_id" json:"-"`
	Id                *string                   `bson:"id,omitempty" json:"id,omitempty"`
	Meta              *Meta                     `bson:"meta,omitempty" json:"meta,omitempty"`
	Identifier        []Identifier              `bson:"identifier,omitempty" json:"identifier,omitempty"`
	Status            EncounterStatus           `bson:"status" json:"status"`
	StatusHistory     []EncounterStatusHistory  `bson:"statusHistory,omitempty" json:"statusHistory,omitempty"`
	Class             string                    `bson:"class" json:"class"`
	Type              []CodeableConcept         `bson:"type,omitempty" json:"type,omitempty"`
	Priority          *CodeableConcept          `bson:"priority,omitempty" json:"priority,omitempty"`
	Patient           Person	                `bson:"patient,omitempty" json:"patient,omitempty"`
	EpisodeOfCare     []Reference               `bson:"episodeOfCare,omitempty" json:"episodeOfCare,omitempty"`
	IncommingReferral []Reference 				`bson:"incomingReferral,omitempty" json:"incommingReferral,omitempty"`
	Participant       []EncounterParticipant    `bson:"participant,omitempty" json:"participant,omitempty"`
	Appointment       []Reference               `bson:"appointment,omitempty" json:"appointment,omitempty"`
	Period            *Period                   `bson:"period,omitempty" json:"period,omitempty"`
	Length            *Duration                 `bson:"length,omitempty" json:"length,omitempty"`
	Reason 	          []CodeableConcept         `bson:"reason,omitempty" json:"reason,omitempty"`
	Indication		  []Reference 				`bson:"indication,omitempty" json:"indication"`
	Hospitalization   *EncounterHospitalization `bson:"hospitalization,omitempty" json:"hospitalization,omitempty"`
	Location          []EncounterLocation       `bson:"location,omitempty" json:"location,omitempty"`
	ServiceProvider   *Reference                `bson:"serviceProvider,omitempty" json:"serviceProvider,omitempty"`
	PartOf            *Reference                `bson:"partOf,omitempty" json:"partOf,omitempty"`
//Not part of DSTU2 Cerner
	Text              *Narrative                `bson:"text,omitempty" json:"text,omitempty"`
	Subject           Person	                `bson:"subject,omitempty" json:"subject,omitempty"`
	CreatedAt 		  *time.Time 				`bson:"created_at"`

// Not in Cerner or 4
	// Extension         []Extension               `bson:"extension,omitempty" json:"extension,omitempty"`
	// ModifierExtension []Extension               `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	// Subject           Person	                `bson:"subject,omitempty" json:"subject,omitempty"`
	// ImplicitRules     *string                   `bson:"implicitRules,omitempty" json:"implicitRules,omitempty"`
	// Language          *string                   `bson:"language,omitempty" json:"language,omitempty"`
	// ClassHistory      []EncounterClassHistory   `bson:"classHistory,omitempty" json:"classHistory,omitempty"`
	// ServiceType       *CodeableConcept          `bson:"serviceType,omitempty" json:"serviceType,omitempty"`
	// BasedOn           []Reference               `bson:"basedOn,omitempty" json:"basedOn,omitempty"`
	// ReasonCode        []CodeableConcept         `bson:"reasonCode,omitempty" json:"reasonCode,omitempty"`
	// ReasonReference   []Reference               `bson:"reasonReference,omitempty" json:"reasonReference,omitempty"`
	// Diagnosis         []EncounterDiagnosis      `bson:"diagnosis,omitempty" json:"diagnosis,omitempty"`
	// Account           []Reference               `bson:"account,omitempty" json:"account,omitempty"`
	
}

type EncounterStatusHistory struct {
	Id                *string         `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension     `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension     `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Status            string `bson:"status" json:"status"`
	Period            Period          `bson:"period" json:"period"`
}
type EncounterClassHistory struct {
	Id                *string     `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Class             Coding      `bson:"class" json:"class"`
	Period            Period      `bson:"period" json:"period"`
}
type EncounterParticipant struct {
	Id                *string           `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension       `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension       `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Type              []CodeableConcept `bson:"type,omitempty" json:"type,omitempty"`
	Period            *Period           `bson:"period,omitempty" json:"period,omitempty"`
	Individual        *Reference        `bson:"individual,omitempty" json:"individual,omitempty"`
}
type EncounterDiagnosis struct {
	Id                *string          `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension      `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension      `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Condition         Reference        `bson:"condition" json:"condition"`
	Use               *CodeableConcept `bson:"use,omitempty" json:"use,omitempty"`
	Rank              *int             `bson:"rank,omitempty" json:"rank,omitempty"`
}
type EncounterHospitalization struct {
	Id                     *string           `bson:"id,omitempty" json:"id,omitempty"`
	Extension              []Extension       `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension      []Extension       `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	PreAdmissionIdentifier *Identifier       `bson:"preAdmissionIdentifier,omitempty" json:"preAdmissionIdentifier,omitempty"`
	Origin                 *Reference        `bson:"origin,omitempty" json:"origin,omitempty"`
	AdmitSource            *CodeableConcept  `bson:"admitSource,omitempty" json:"admitSource,omitempty"`
	ReAdmission            *CodeableConcept  `bson:"reAdmission,omitempty" json:"reAdmission,omitempty"`
	DietPreference         []CodeableConcept `bson:"dietPreference,omitempty" json:"dietPreference,omitempty"`
	SpecialCourtesy        []CodeableConcept `bson:"specialCourtesy,omitempty" json:"specialCourtesy,omitempty"`
	SpecialArrangement     []CodeableConcept `bson:"specialArrangement,omitempty" json:"specialArrangement,omitempty"`
	Destination            *Reference        `bson:"destination,omitempty" json:"destination,omitempty"`
	DischargeDisposition   *CodeableConcept  `bson:"dischargeDisposition,omitempty" json:"dischargeDisposition,omitempty"`
}
type EncounterLocation struct {
	Id                *string                  `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension              `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension              `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Location          Reference                `bson:"location" json:"location"`
	Status            EncounterLocationStatus  `bson:"status,omitempty" json:"status,omitempty"`
	PhysicalType      *CodeableConcept         `bson:"physicalType,omitempty" json:"physicalType,omitempty"`
	Period            *Period                  `bson:"period,omitempty" json:"period,omitempty"`
}

type EncounterLocationStatus string 

type OtherEncounter Encounter

type EncounterStatus string


type Duration struct {
	Id         *string             `bson:"id,omitempty" json:"id,omitempty"`
	Extension  []Extension         `bson:"extension,omitempty" json:"extension,omitempty"`
	Value      *string             `bson:"value,omitempty" json:"value,omitempty"`
	Comparator *QuantityComparator `bson:"comparator,omitempty" json:"comparator,omitempty"`
	Unit       *string             `bson:"unit,omitempty" json:"unit,omitempty"`
	System     *string             `bson:"system,omitempty" json:"system,omitempty"`
	Code       *string             `bson:"code,omitempty" json:"code,omitempty"`
}


// type EncounterLocationStatus int

// const (
// 	EncounterLocationStatusPlanned EncounterLocationStatus = iota
// 	EncounterLocationStatusActive
// 	EncounterLocationStatusReserved
// 	EncounterLocationStatusCompleted
// )

// func (code EncounterLocationStatus) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(code.Code())
// }
// func (code *EncounterLocationStatus) UnmarshalJSON(json []byte) error {
// 	s := strings.Trim(string(json), "\"")
// 	switch s {
// 	case "planned":
// 		*code = EncounterLocationStatusPlanned
// 	case "active":
// 		*code = EncounterLocationStatusActive
// 	case "reserved":
// 		*code = EncounterLocationStatusReserved
// 	case "completed":
// 		*code = EncounterLocationStatusCompleted
// 	default:
// 		return fmt.Errorf("unknown EncounterLocationStatus code `%s`", s)
// 	}
// 	return nil
// }
// func (code EncounterLocationStatus) String() string {
// 	return code.Code()
// }
// func (code EncounterLocationStatus) Code() string {
// 	switch code {
// 	case EncounterLocationStatusPlanned:
// 		return "planned"
// 	case EncounterLocationStatusActive:
// 		return "active"
// 	case EncounterLocationStatusReserved:
// 		return "reserved"
// 	case EncounterLocationStatusCompleted:
// 		return "completed"
// 	}
// 	return "<unknown>"
// }
// func (code EncounterLocationStatus) Display() string {
// 	switch code {
// 	case EncounterLocationStatusPlanned:
// 		return "Planned"
// 	case EncounterLocationStatusActive:
// 		return "Active"
// 	case EncounterLocationStatusReserved:
// 		return "Reserved"
// 	case EncounterLocationStatusCompleted:
// 	}
// 	return "<unknown>"
// }

// func (code EncounterLocationStatus) Definition() string {
// 	switch code {
// 	case EncounterLocationStatusPlanned:
// 		return "The patient is planned to be moved to this location at some point in the future."
// 	case EncounterLocationStatusActive:
// 		return "The patient is currently at this location, or was between the period specified.\r\rA system may update these records when the patient leaves the location to either reserved, or completed."
// 	case EncounterLocationStatusReserved:
// 		return "This location is held empty for this patient."
// 	case EncounterLocationStatusCompleted:
// 		return "The patient was at this location during the period specified.\r\rNot to be used when the patient is currently at the location."
// 	}
// 	return "<unknown>"
// }



type Meta struct {
	Id          *string     `bson:"id,omitempty" json:"id,omitempty"`
	Extension   []Extension `bson:"extension,omitempty" json:"extension,omitempty"`
	VersionId   *string     `bson:"versionId,omitempty" json:"versionId,omitempty"`
	LastUpdated *string     `bson:"lastUpdated,omitempty" json:"lastUpdated,omitempty"`
	Source      *string     `bson:"source,omitempty" json:"source,omitempty"`
	Profile     []string    `bson:"profile,omitempty" json:"profile,omitempty"`
	Security    []Coding    `bson:"security,omitempty" json:"security,omitempty"`
	Tag         []Coding    `bson:"tag,omitempty" json:"tag,omitempty"`
}





