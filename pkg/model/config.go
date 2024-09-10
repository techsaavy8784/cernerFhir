package model

import (
	"fmt"
	//"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	//"os"
	fhir "github.com/dhf0820/cernerFhir/fhirongo"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"

	//"github.com/jingweno/conf"
	//"crypto/tls"
	//"gopkg.in/mgo.v2"
	//"net"

	"github.com/dhf0820/cernerFhir/pkg/storage"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	router   *mux.Router
	database *storage.MongoDB
	//	emr           	*EMR
	fhirURL        string
	fhirConnection *fhir.Connection
	port           string
	env            *Environment
	source         string
	mrnID          string
	appURL         string
	appName        string
	mode           string
	serverURL      string
	imageURL       string
	recordLimit    string
	logLevel       log.Level
}

var config Config
var name string
var mode string
var appEnv Environment
var port string

func InitializeAll(dbURL string) *Config {
	fmt.Printf("Initializing MongoURL: %s\n", dbURL)
	//fmt.Printf("InitializeAll called\n")
	name, ok := os.LookupEnv("FHIR") // Name of fhir connection in DB Cerner,...
	if !ok {
		log.Errorf("FHIR Environment is not configured")
		log.Fatal("FHIR Environment is not set")
	}
	mode, ok = os.LookupEnv("MODE")
	fmt.Printf("Using run mode : %s\n", mode)
	if !ok {
		mode = "dev" // Force dev if not set
	}
	log.Infof("Call initConfig with DB: %s", dbURL)
	log.Infof("Name: %s    Mode: %s\n", name, mode)
	config = initConfig(name, mode, dbURL)
	port = config.Port()

	log.Infof("Starting %s FHIR interface VERSION %s in %s mode", config.Source(), config.ServerVersion(), mode)
	// var env m.Environment
	return &config
}

func initConfig(name string, mode string, mongoURL string) Config {

	var err error
	var ok bool
	if mongoURL == "" {
		mongoURL, ok = os.LookupEnv("MONGODB")
		if !ok {
			log.Panic("initConfig:80 -- MONGODB is not defined")
		}
	}
	config.database, err = storage.Open(mongoURL)
	if err != nil {
		log.Fatal("Database is not configured:")
	}
	//fmt.Printf("Opened database: %s\n", mongoURL)
	var appEnv Environment
	appEnv.Name = "cerner" //TODO: Remove hardcode environment vals
	appEnv.Type = "fhir"
	appEnv.Mode = mode
	config.mode = mode
	config.appName = appEnv.Name

	//TODO: change type from app to system to match

	err = appEnv.FindOne()
	if err != nil {
		log.Fatalf("app: fhir is not configured: %s", err.Error())
	}
	log.Debugf("env: %v", appEnv)
	config.serverURL = appEnv.Env["baseURL"]
	config.appURL = appEnv.Env["baseURL"]

	log.Debugf("Urls: server: %s,  app: %s", config.serverURL, config.appURL)
	config.appName = appEnv.Env["appName"]
	config.port = appEnv.Env["port"]
	config.SetLogLevel(appEnv.Env["logLevel"])

	//log.Debugf("Config Using PORT: %s", config.port)

	// Initial env findOne is for the app aka system url
	//this is for the service.
	var env Environment
	config.env = &env
	env.Name = name
	env.Type = "fhir"
	err = env.FindOne()
	if err != nil {
		msg := fmt.Sprintf("Config for service/%s is not configured:", name)
		log.Errorf("%s\n", msg)
		log.Fatal(msg)
	}
	config.fhirURL = env.Env["fhir_url"]
	config.imageURL = env.Env["imageURL"]
	config.recordLimit = env.Env["recordLimit"]
	config.source = env.Env["source"]
	config.mrnID = env.Env["mrn_id"]
	//config.port = env.Env["port"]

	err = initializeEMR() //fhir.New(config.fhirURL)
	if err != nil {
		log.Fatal(fmt.Errorf("InitConfig emr %s failed: %v ", config.fhirURL, err.Error()))
	}

	return config
}

func updateEnv() {
	log.Debugf("In updateEnv")
	err := appEnv.FindOne()
	if err != nil {
		log.Fatal("app: fhir is not configured:")
	}
	log.Debugf("env: %v\n", appEnv)
	config.SetLogLevel(appEnv.Env["logLevel"])

	config.serverURL = appEnv.Env["baseURL"]
	config.appURL = appEnv.Env["baseURL"]

	log.Debugf("Urls: S: %s,  a: %s\n", config.serverURL, config.appURL)
	config.appName = appEnv.Env["appName"]
	config.port = appEnv.Env["port"]

	//log.Debugf("Config Using PORT: %s", config.port)
	var env Environment
	config.env = &env
	env.Name = name
	env.Type = "fhir"
	err = env.FindOne()
	if err != nil {
		msg := fmt.Sprintf("Config for service: %s is not configured:", name)
		log.Errorf("%s\n", msg)
		log.Fatal(msg)
	}
	config.fhirURL = env.Env["fhir_url"]
	config.imageURL = env.Env["imageURL"]
	config.recordLimit = env.Env["recordLimit"]
	config.source = env.Env["source"]
	config.mrnID = env.Env["mrn_id"]
}

func ActiveConfig() *Config {
	return &config
}

func FhirPdfUrl() string {
	return ActiveConfig().Env("fhirPdfUrl")
}

func RecordLimit() int64 {
	limit, err := strconv.ParseInt(ActiveConfig().Env("recordLimit"), 10, 64)
	if err != nil {
		limit = 20
	}
	return limit
}

func PageSize() int64 {
	pageSize, err := strconv.ParseInt(ActiveConfig().Env("page_size"), 10, 64)
	if err != nil {
		pageSize = 20
	}
	return pageSize
}

func LinesPerPage() int64 {
	pageSize, err := strconv.ParseInt(ActiveConfig().Env("page_size"), 10, 64)
	if err != nil {
		pageSize = 20
	}
	return pageSize
}

func ExpireCacheAfter() int64 {
	expireAfter, err := strconv.ParseInt(ActiveConfig().Env("cacheExpireAfter"), 10, 64)
	if err != nil {
		expireAfter = 600 //seconds 10 minutes
	}
	return expireAfter
}

func LoginExpiresAfter() int64 {
	expireAfter, err := strconv.ParseInt(ActiveConfig().Env("loginExpireAfter"), 10, 64)
	if err != nil {
		expireAfter = 15 //seconds 10 minutes
	}
	return expireAfter
}

func CacheCheckFreq() int64 {
	checkFreq, err := strconv.ParseInt(ActiveConfig().Env("cacheCheckFreq"), 10, 64)
	if err != nil {
		checkFreq = 3600 // 1 hour
	}
	return checkFreq
}

func initializeEMR() error {
	log.Infof("%s's fhirUrl: %s\n", config.source, config.fhirURL)
	config.fhirConnection = fhir.New(config.fhirURL)
	err := checkEMRConnection(config.fhirConnection)
	return err
}

func checkEMRConnection(c *fhir.Connection) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, "metadata")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json+fhir")
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("initFhir: did not connect to %s", resp.Status)
	}
	return nil

}

func (c *Config) AppURL() string {
	return c.appURL
}

func (c *Config) ImageURL() string {
	return c.imageURL
}

func (c *Config) AppName() string {
	return c.appName
}

func (c *Config) ServerVersion() string {
	return Version
}

func (c *Config) Env(key string) string {
	return c.env.Env[key]
}

func (c *Config) SetRouter(m *mux.Router) {
	c.router = m
}

func (c *Config) Router() *mux.Router {
	return c.router
}

func (c *Config) DBClient() *mongo.Client {
	db := c.database
	return db.Client
}

func (c *Config) FhirURL() string {
	return c.fhirURL
}

func (c *Config) Fhir() *fhir.Connection {
	return c.fhirConnection
}

func Fhir() *fhir.Connection {
	return config.fhirConnection
}

func (c *Config) Port() string {
	return c.port
}

func (c *Config) Source() string {
	return c.source
}

func (c *Config) RecordLimit() string {
	return c.recordLimit
}

func (c *Config) MrnID() string {
	return c.mrnID
}

func (c *Config) Mode() string {
	return c.mode
}

func (c *Config) Database() *storage.MongoDB {
	return c.database
}

func (c *Config) BaseUrl() string {
	return c.appURL
}

func (c *Config) SetLogLevel(level string) {
	//log.Infof("Config Setting log level to %s", level)
	switch level {
	case "debug":
		c.logLevel = log.DebugLevel
	case "warn":
		c.logLevel = log.WarnLevel
	case "error":
		c.logLevel = log.ErrorLevel
	default:
		c.logLevel = log.InfoLevel
	}
	log.SetLevel(c.logLevel)
}
