package dbconnect

import (
	"log"

	"gopkg.in/mgo.v2"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
)

/*Global Mongodb package*/
var MongoSession *mgo.Session
var MongoDatabase *mgo.Database
var LoadConfigMongoDB = config.DataConfig.Mongo["users_system"]
var err error

func init() {
	MongoSession, err = mgo.Dial(LoadConfigMongoDB.Host)
	if err != nil {
		log.Fatal("Failed to start the Mongo session")
	}
	//load database name
	if MongoDatabase == nil {
		SetDatabaseMongoDB()
	}
}

/*Mongodb set database*/
func SetDatabaseMongoDB() {
	MongoDatabase = MongoSession.DB(LoadConfigMongoDB.Name)
}

/*Mongodb connect clone Session*/
func GetMongoSessionClone() *mgo.Session {
	return MongoSession.Clone()
}

/*Mongodb connect copy Session*/
func GetMongoSessionCopy() *mgo.Session {
	return MongoSession.Copy()
}

/*Get Collection*/
func GetMongoCollection(name string) *mgo.Collection {
	return MongoDatabase.C(name)
}

