package fhirongo

import (
	"time"
	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

// Address is a physical address
type Address struct {
	Use        string   `json:"use" bson:"use"`
	Line       []string `json:"line" bson:"line"`
	City       string   `json:"city" bson:"city"`
	State      string   `json:"state" bson:"state"`
	PostalCode string   `json:"postalCode" bson:"postal_code"`
	Country    string   `json:"country" bson:"country"`
	Period     Period   `json:"period,omitempty" bson:"period"`
}

// Attachment is a url attachment
type Attachment struct {
	ContentType string       `json:"contentType" bson:"cont_type"`
	Language    string       `json:"language"`
	Data        base64Binary `json:"base64binary"`
	URL         string       `json:"url"`
	Title       string       `json:"title"`
}

type BackboneElement struct {
	Type       []CodeableConcept `json:"codeable_concept" bson:"codeable_concept"`
	Period     Period            `json:"period"`
	Individual Person            `json:"person"`
}

type base64Binary string

type BinaryContent struct {
	Attachment struct {
		ContentType string `bson:"content_type" json:"contentType"`
		URL         string `bson:"url" json:"url"`
	} `bson:"attachment" json:"attachment"`
}

// Bundle is the header for sets of information
type Bundle struct {
	ResourceType string `json:"resourceType" bson:"resource_type"`
	ID           string `json:"id"`
	Type         string `json:"type"`
	Link         []Link `json:"link"`
	Total        int    `json:"total"`
}

//Category the DiagnosticReport Category
type Category struct {
	CodeableConcept
	Text string `json:"text"`
}

type Code string

type CodeableConcept struct {
	Coding []Coding
	Text   string `json:"text"`
}

// Coding is a code and system
type Coding struct {
	System       string `json:"system"`
	Version      string `json:"version"`
	Code         string `json:"code"`
	Display      string `json:"display"`
	UserSelected bool   `json:"userSelected" bson:"user_selected"`
}

// CodeText is a healthcare condition
type CodeText struct {
	Code Note `json:"code"`
}

// Communication is the language people speak
type Communication struct {
	Preferred bool    `json:"preferred"`
	Language  Concept `json:"language"`
}

// Concept is a general concept such as language
type Concept struct {
	Coding []Coding `json:"coding"`
	Text   string   `json:"text"`
}

type ContactPoint struct {
	Id        *string             `bson:"id,omitempty" json:"id,omitempty"`
	Extension []Extension         `bson:"extension,omitempty" json:"extension,omitempty"`
	System    *ContactPointSystem `bson:"system,omitempty" json:"system,omitempty"`
	Value     *string             `bson:"value,omitempty" json:"value,omitempty"`
	Use       *ContactPointUse    `bson:"use,omitempty" json:"use,omitempty"`
	Rank      *int                `bson:"rank,omitempty" json:"rank,omitempty"`
	Period    *Period             `bson:"period,omitempty" json:"period,omitempty"`
}

// Context encounter only initially
type Context struct {
	EncounterRef EncounterReference `json:"encounter"`
}

// Concept is a general concept such as language

type ContactPointSystem int
type ContactPointUse int
type DaysOfWeek int

// DispenseRequest is a dispensing request
type DispenseRequest struct {
	ValidityPeriod Period `json:"validityPeriod" bson:"validity_period"`
}

// DosageInstruction are the medication instructions
type DosageInstruction struct {
	Text            string   `json:"text"`
	AsNeededBoolean bool     `json:"asNeededBoolean" bson:"as_needed_boolean"`
	Route           Concept  `json:"route"`
	Method          Concept  `json:"method"`
	Timing          Timing   `json:"timing"`
	DoseQuantity    Quantity `json:"doseQuantity" bson:"dose_quanatity"`
}

type EncounterContext struct {
	Encounter struct {
		Reference string `json:"reference"`
	} `json:"encounter"`
	Period Period `json:"period"`
}

//EncounterReference of the report
type EncounterReference struct {
	Reference string `json:"reference"`
}

// Entry are the common entry fields
type Entry struct {
	FullURL string `json:"fullUrl"`
	// 	Resource Resource  `json:"resource"`
	Link   []Link     `json:"link"`
	Search SearchMode `json:"search"`
}

// EntryPartial are the common entry fields
type EntryPartial struct {
	FullURL string     `json:"fullUrl"`
	Link    []Link     `json:"link"`
	Search  SearchMode `json:"search"`
}

// Extension is a codified FHIR extension
type Extension struct {
	URL                  string  `json:"url"`
	ValueCodeableConcept Concept `json:"valueCodeableConcept" bson:"value_codeable_concept"`
}

// type EncounterHospitalization struct {
// 	Id                     *string           `bson:"id,omitempty" json:"id,omitempty"`
// 	Extension              []Extension       `bson:"extension,omitempty" json:"extension,omitempty"`
// 	ModifierExtension      []Extension       `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
// 	PreAdmissionIdentifier *Identifier       `bson:"preAdmissionIdentifier,omitempty" json:"preAdmissionIdentifier,omitempty"`
// 	Origin                 *Reference        `bson:"origin,omitempty" json:"origin,omitempty"`
// 	AdmitSource            *CodeableConcept  `bson:"admitSource,omitempty" json:"admitSource,omitempty"`
// 	ReAdmission            *CodeableConcept  `bson:"reAdmission,omitempty" json:"reAdmission,omitempty"`
// 	DietPreference         []CodeableConcept `bson:"dietPreference,omitempty" json:"dietPreference,omitempty"`
// 	SpecialCourtesy        []CodeableConcept `bson:"specialCourtesy,omitempty" json:"specialCourtesy,omitempty"`
// 	SpecialArrangement     []CodeableConcept `bson:"specialArrangement,omitempty" json:"specialArrangement,omitempty"`
// 	Destination            *Reference        `bson:"destination,omitempty" json:"destination,omitempty"`
// 	DischargeDisposition   *CodeableConcept  `bson:"dischargeDisposition,omitempty" json:"dischargeDisposition,omitempty"`
// }

// Identifier can identify things
type Identifier struct {
	Use    string  `json:"use"`
	System string  `json:"system"`
	Value  string  `json:"value"`
	Type   Concept `json:"type"`
	Period Period  `json:"period"`
}

type Individual struct {
	Reference string `json:"reference"`
	Display   string `json:"display"`
}

// Issue is a FHIR issue
type Issue struct {
	Severity string   `json:"severity"`
	Location []string `json:"location"`
	Code     string   `json:"code"`
	Details  Concept  `json:"details"`
}

// Link is a resource link
type Link struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

// type EncounterLocation struct {
// 	Id                *string                  `bson:"id,omitempty" json:"id,omitempty"`
// 	Extension         []Extension              `bson:"extension,omitempty" json:"extension,omitempty"`
// 	ModifierExtension []Extension              `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
// 	Location          Reference                `bson:"location" json:"location"`
// 	Status            *EncounterLocationStatus `bson:"status,omitempty" json:"status,omitempty"`
// 	PhysicalType      *CodeableConcept         `bson:"physicalType,omitempty" json:"physicalType,omitempty"`
// 	Period            *Period                  `bson:"period,omitempty" json:"period,omitempty"`
// }

//type EncounterLocationStatus string

type Location struct {
	Id                     *string                    `bson:"id,omitempty" json:"id,omitempty"`
	Meta                   *Meta                      `bson:"meta,omitempty" json:"meta,omitempty"`
	ImplicitRules          *string                    `bson:"implicitRules,omitempty" json:"implicitRules,omitempty"`
	Language               *string                    `bson:"language,omitempty" json:"language,omitempty"`
	Text                   *Narrative                 `bson:"text,omitempty" json:"text,omitempty"`
	Extension              []Extension                `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension      []Extension                `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Identifier             []Identifier               `bson:"identifier,omitempty" json:"identifier,omitempty"`
	Status                 *LocationStatus            `bson:"status,omitempty" json:"status,omitempty"`
	OperationalStatus      *Coding                    `bson:"operationalStatus,omitempty" json:"operationalStatus,omitempty"`
	Name                   *string                    `bson:"name,omitempty" json:"name,omitempty"`
	Alias                  []string                   `bson:"alias,omitempty" json:"alias,omitempty"`
	Description            *string                    `bson:"description,omitempty" json:"description,omitempty"`
	Mode                   *LocationMode              `bson:"mode,omitempty" json:"mode,omitempty"`
	Type                   []CodeableConcept          `bson:"type,omitempty" json:"type,omitempty"`
	Telecom                []ContactPoint             `bson:"telecom,omitempty" json:"telecom,omitempty"`
	Address                *Address                   `bson:"address,omitempty" json:"address,omitempty"`
	PhysicalType           *CodeableConcept           `bson:"physicalType,omitempty" json:"physicalType,omitempty"`
	Position               *LocationPosition          `bson:"position,omitempty" json:"position,omitempty"`
	ManagingOrganization   *Reference                 `bson:"managingOrganization,omitempty" json:"managingOrganization,omitempty"`
	PartOf                 *Reference                 `bson:"partOf,omitempty" json:"partOf,omitempty"`
	HoursOfOperation       []LocationHoursOfOperation `bson:"hoursOfOperation,omitempty" json:"hoursOfOperation,omitempty"`
	AvailabilityExceptions *string                    `bson:"availabilityExceptions,omitempty" json:"availabilityExceptions,omitempty"`
	Endpoint               []Reference                `bson:"endpoint,omitempty" json:"endpoint,omitempty"`
}

type LocationMode string
type LocationStatus string

type LocationPosition struct {
	Id                *string     `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	Longitude         string      `bson:"longitude" json:"longitude"`
	Latitude          string      `bson:"latitude" json:"latitude"`
	Altitude          *string     `bson:"altitude,omitempty" json:"altitude,omitempty"`
}
type LocationHoursOfOperation struct {
	Id                *string      `bson:"id,omitempty" json:"id,omitempty"`
	Extension         []Extension  `bson:"extension,omitempty" json:"extension,omitempty"`
	ModifierExtension []Extension  `bson:"modifierExtension,omitempty" json:"modifierExtension,omitempty"`
	DaysOfWeek        []DaysOfWeek `bson:"daysOfWeek,omitempty" json:"daysOfWeek,omitempty"`
	AllDay            *bool        `bson:"allDay,omitempty" json:"allDay,omitempty"`
	OpeningTime       *string      `bson:"openingTime,omitempty" json:"openingTime,omitempty"`
	ClosingTime       *string      `bson:"closingTime,omitempty" json:"closingTime,omitempty"`
}
type OtherLocation Location

//MetaData meta field in DocumentReference/DiagnosticReport
type MetaData struct {
	VersionID   string    `json:"versionId" bson:"version_id"`
	LastUpdated time.Time `json:"lastUpdated" bson:"last_updated"`
}

// Name is a persons name
type Name struct {
	Use    string   `json:"use"`
	Family []string `json:"family"`
	Given  []string `json:"given"`
	Suffix []string `json:"suffix"`
	Prefix []string `json:"prefix"`
}

type Narrative struct {
	Id        *string         `bson:"id,omitempty" json:"id,omitempty"`
	Extension []Extension     `bson:"extension,omitempty" json:"extension,omitempty"`
	Status    NarrativeStatus `bson:"status" json:"status"`
	Div       string          `bson:"div" json:"div"`
}

type NarrativeStatus string

// Note is any general note on some other component
type Note struct {
	Text string `json:"text"`
}

type Participant struct {
	Id                *string          `json:"id,omitempty" bson:"id,omitempty"`
	Extension         []Extension      `json:"extension,omitempty" bson:"extension,omitempty"`
	ModifierExtension []Extension      `json:"modifierExtension,omitempty" bson:"modifierExtension,omitempty"`
	Type              *CodeableConcept `json:"type,omitempty" bson:"type,omitempty"`
	Period            Period           `json:"period" bson:"period"`
	Individual        *Reference       `json:"reference,omitempty" bson:"reference,omitempty"`
}

// Period is a period of time
type Period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Person is a human
type Person Thing

// Quantity is a quantity of something
type Quantity struct {
	Value  float32 `json:"value"`
	Unit   string  `json:"unit"`
	Code   string  `json:"code"`
	System string  `json:"system"`
}

// Reaction is a human reaction
type Reaction struct {
	Certainty     string    `json:"certainty"`
	Onset         time.Time `json:"onset"`
	Manifestation []Note    `json:"manifestation"`
	Note          Note      `json:"note"`
}

type Reason struct {
	Text string `json:"text"`
}

type Reference struct {
	Reference  string     `json:"reference"`
	Identifier Identifier `json:"identifier"`
	Display    string     `json:"display"`
}

// Repeat is a time based repeat of something
type Repeat struct {
	Frequency    float32 `json:"frequency"`
	Period       float32 `json:"period"`
	PeriodUnits  string  `json:"periodUnits"`
	BoundsPeriod Period  `json:"boundsPeriod" bson:"bounds_period"`
}

// ResourcePartial are the common resource fields
type ResourcePartial struct {
	ResourceType      string             `json:"resourceType" bson:"resource_type"`
	EffectiveDateTime time.Time          `json:"effectiveDateTime" bson:"effective_dat_time"`
	RecordedDate      time.Time          `json:"recordedDate" bson:"recorded_date"`
	Status            string             `json:"status"`
	ID                string             `json:"id"`
	Subject           Person             `json:"subject"`
	Patient           Person             `json:"patient"`
	Performer         Person             `json:"performer"`
	Recorder          Person             `json:"recorder"`
	Encounter         EncounterReference `json:"Encounter"`
}

// SearchMode is a FHIR search mode
type SearchMode struct {
	Mode string `json:"mode"`
}

// SearchResult is a search result
type SearchResult struct {
	ResourceType string  `json:"resourceType" bson:"resource_type"`
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	Link         []Link  `json:"link"`
	Total        int     `json:"total"`
	Issues       []Issue `json:"issue"`
	//	Entries 		 []Entry `json:"entry"`
}

type ServiceProvider struct {
	Reference string `json:"reference"`
}

type Status string

//TextData is the html text
type TextData struct {
	Status string `json:"status"`
	Div    string `json:"div"`
}

// Telecom is a phone number
type Telecom struct {
	System string `json:"system"`
	Value  string `json:"value"`
	Use    string `json:"use,omitempty"`
	Period Period `json:"period,omitempty"`
}

// Thing is a FHIR thing
type Thing struct {
	Display   string `json:"display"`
	Reference string `json:"reference"`
}
type Timing struct {
	Repeat Repeat `json:"timing"`
}

// type Type struct {

// }
