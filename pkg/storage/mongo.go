package storage

import (
	"context"
	"errors"
	"fmt"
	"os"

	//"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client       *mongo.Client
	DatabaseName string
	url          string
}

var db MongoDB

func Open(url string) (*MongoDB, error) {
	if url == "" {
		//mt.Printf("No url provided\n")
		exists := false
		url, exists = os.LookupEnv("MONGODB")
		if !exists{
			log.Fatal("mongodb environment missing!")
		}
	}
	log.Infof("Using mongodb: %s", url)
	databaseName, exists := os.LookupEnv("DATABASE")
	if !exists {
		log.Fatal("database environment missing!")
	}
	//fmt.Printf("mongo: url: %s database: %s\n", url, databaseName)

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err == nil {
		err = client.Ping(context.TODO(), nil)
		//fmt.Printf("Ping server to actually open\n")
	}
	if err != nil {
		fmt.Printf("Ping failed: %s\n", err)
		log.Fatal(err)
	}
	// res, err := client.ListDatabases(context.Background(), bson.D{})
	// if err != nil{
	// 	log.Errorf("Error Listing Databases: %s", err.Error())
	// }
	// fmt.Printf("Databases: %s\n", spew.Sdump(res.Databases))
	// myDB := client.Database("fhir")
	// listCollections(myDB)
	db.Client = client
	db.url = url
	db.DatabaseName = databaseName
	//fmt.Printf("Returning from Mongo Open\n")
	return &db, err
}

// func listCollections(db *mongo.Database){
	    
//     // use a filter to only select capped collections
// 	fmt.Printf("Collections in: %s\n", db.Name())
//     result, err := db.ListCollectionNames(
//         context.TODO(), bson.D{})
//        // bson.D{{"options.capped", false}})

//     if err != nil {
//         log.Fatal(err)
//     }
// 	fmt.Println("List of Collections")
	
//     for _, coll := range result {
//         fmt.Println(coll)
//     }
// }

func Current() (*MongoDB, error) {
	if db.Client != nil {
		return &db, nil
	}
	_, err := Open("")
	return &db, err
}

func (db *MongoDB) Close() error {
	err := db.Client.Disconnect(context.TODO())
	return err
}

func GetCollection(collection string) (mongo.Collection, error) {
	db, err := Current() //"mongodb://admin:Sacj0nhat1@cat.vertisoft.com:27017")
	if err != nil {
		log.Fatal(err)
		//return nil, err
	}
	client := db.Client
	coll := client.Database(db.DatabaseName).Collection(collection)
	return *coll, nil
}

func GetSession() (mongo.Session, error) {
	db, err := Current() // get the current mongo connection
	if err != nil {
		log.Fatal(err)
	}
	client := db.Client
	return client.StartSession()
}

func StartTransaction(session mongo.Session) error {
	return session.StartTransaction()
}

func Client() *mongo.Client {
	db, err := Current() // get the current mongo connection
	if err != nil {
		log.Fatal(err)
	}
	return db.Client
}

func IsDup(err error) bool {
	var e mongo.WriteException
    if errors.As(err, &e) {
        for _, we := range e.WriteErrors {
            if we.Code == 11000 {
				//fmt.Printf("ErrCode: %d\n", we.Code)
                return true
            } else {
				fmt.Printf("ErrorCode: %d\n", we.Code)
				return false
			}
        }
    }
    return false
}