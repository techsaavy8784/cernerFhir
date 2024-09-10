package model

import (
	"context"
	"fmt"

	"time"

	log "github.com/sirupsen/logrus"

	//"github.com/davecgh/go-spew/spew"
	"github.com/dhf0820/cernerFhir/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Environment struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Type      string             `json:"type" bson:"type"`
	Name      string             `json:"name" bson:"name"`
	Mode      string             `json:"mode" bson:"mode"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Env       map[string]string  `json:"env" bson:"env"`
	test      string
}

func (e *Environment) Find() (*mongo.Cursor, error) {
	// create the bson.D from the fields with data in p
	filter := bson.D{{"type", e.Type}, {"name", e.Name}}
	//filter := createFilter(*p)
	cursor, err := Find(filter)
	return cursor, err
}

func (e *Environment) FindOne() error {
	//fmt.Printf("Environment Looking for %s\n", spew.Sdump(e))
	// create the bson.D from the fields with data in p
	filter := bson.D{{"type", e.Type}, {"name", e.Name}}
	//fmt.Printf("Filter to FindOne: %v\n",filter)
	var err error
	collection, _ := storage.GetCollection("environments")
	//envs := []*Environment{}
	//filter = bson.D{}
	err = collection.FindOne(context.TODO(), filter).Decode(e)
	if err != nil {
		log.Errorf("Environment findone failed filter: %v err:%s\n", filter, err.Error())
		log.Fatal(err)
	}
	//fmt.Printf("Environment Found: %s\n", spew.Sdump(e))
	//filter := createFilter(*p)
	//e, err := FindOne(filter)

	return err
}

// func FindOne(filter bson.D) (*mongo.Cursor, error) {
// 	var result *mongo.Cursor
// 	var err error
// 	collection, _ := storage.GetCollection("environments")

// 	result, err = collection.Find(context.TODO(), filter).Decode(&result)

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return result, err
// }

func Find(filter bson.D) (*mongo.Cursor, error) {
	var result *mongo.Cursor
	var err error
	collection, _ := storage.GetCollection("environments")

	result, err = collection.Find(context.TODO(), filter) //.Decode(&result)

	if err != nil {
		log.Fatal(err)
	}

	return result, err
}

// Insert adds one emvironment to the pending. Checks if already exists and if there returns existing.
func (e *Environment) Insert() error {
	//fmt.Printf("adding: %T: %v\n\n", e, e)
	//fmt.Printf("add Name: %s\n", e.Name)
	// _, err := FindByPhone(c.FaxNumber, c.Facility)
	// log.Fatal(err)
	collection, _ := storage.GetCollection("environments")

	insertResult, err := collection.InsertOne(context.TODO(), e)
	if err != nil {
		log.Fatal(err)
	}
	if err == nil {
		e.ID = insertResult.InsertedID.(primitive.ObjectID)

	}
	fmt.Printf("New Environmentt: %v\n", e)
	return err
}
