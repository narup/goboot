// data - Data layer
//
// Copyright (c) 2015-2016 - Puran Singh <puran@phil.us>
//
// All rights reserved.

package data

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/narup/gmgo"

	"database/sql"

	"github.com/phil-inc/plib/core/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mainDb gmgo.Db
var cacheDb *redis.DB
var PostgresDb *sql.DB

// ErrNotFound indicates data not found
var ErrNotFound = errors.New("not found").Error()

// BaseData represents the common data used across all the data models.
type BaseData struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty" pson:"id"`
	CreatedDate *time.Time    `json:"createdDate" bson:"createdDate,omitempty" pson:"created_date"`
	UpdatedDate *time.Time    `json:"updatedDate" bson:"updatedDate,omitempty" pson:"updated_date"`
}

// InitData function initializes data representation with common attributes
func (data *BaseData) InitData() {
	data.ID = bson.NewObjectId()

	ts := time.Now().UTC()
	data.CreatedDate = &ts
	data.UpdatedDate = &ts
}

// StringID returns hex id value
func (data *BaseData) StringID() string {
	return data.ID.Hex()
}

// ObjectID returns bson.ObjectId for hex string id
func ObjectID(id string) bson.ObjectId {
	return bson.ObjectIdHex(id)
}

// Session returns the mongodb session copied from
// main session
func Session() *gmgo.DbSession {
	return mainDb.Session()
}

// DB returns the mongodb session copied from
// main session
func DB() *gmgo.Db {
	return &mainDb
}

//Cache returns cache db
func Cache() *redis.DB {
	return cacheDb
}

// SetupRedisDb setup redis database
func SetupRedisDb(url string) error {
	db, err := redis.Connect(url)
	if err != nil {
		return err
	}
	cacheDb = db
	return nil
}

//SetupMongoDb the MongoDB connection
func SetupMongoDb(dbConfig gmgo.DbConfig) (gmgo.Db, error) {
	//setup PhilDB Mongo database connection
	if err := gmgo.Setup(dbConfig); err != nil {
		log.Fatalf("Error setting up data layer : %s %+v.\n", err, dbConfig)
		return gmgo.Db{}, err
	}

	newDb, err := gmgo.Get(dbConfig.DBName)
	if err != nil {
		log.Fatalf("Db connection error : %s.\n", err)
	}

	mainDb = newDb
	return newDb, nil
}

//SetupPostgres - creates connection to Postgres database
func SetupPostgres(url, dbName string) error {
	connURL := fmt.Sprintf("%s/%s?sslmode=disable", url, dbName)

	db, err := sql.Open("postgres", connURL)

	if err != nil {
		return err
	}

	PostgresDb = db
	log.Println("Connected to postgres")
	return nil
}

// Close database connections and release resources
func Close() {
}

// DBRef returns the mgo DBRef for given collection name and object ID
func DBRef(collectionName string, ID bson.ObjectId) *mgo.DBRef {
	return &mgo.DBRef{Collection: collectionName, Id: ID}
}

//RefID returns string HEX id for DBRef
func RefID(dbRef *mgo.DBRef) string {
	if dbRef == nil {
		return ""
	}
	if bsonID, ok := dbRef.Id.(bson.ObjectId); ok {
		return bsonID.Hex()
	}
	return ""
}
