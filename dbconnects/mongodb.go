package dbconnect

import (
	"crypto/tls"
	"log"
	"net"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"bitbucket.org/2tgroup/ciwp-api-users/config"
)

/*Global Mongodb package*/
var MongoSession *mgo.Session
var MongoDatabase *mgo.Database
var LoadConfigMongoDB = config.DataConfig.Mongo["default"]
var err error

func init() {

	MongoSession, err = mgo.Dial(LoadConfigMongoDB.Host)

	if err != nil {

		tlsConfig := &tls.Config{}

		tlsConfig.InsecureSkipVerify = true

		dialInfo, errDial := mgo.ParseURL(LoadConfigMongoDB.Host)

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, errx := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, errx
		}

		if errDial != nil {
			log.Fatal("Can not Parse URL : ", errDial)
		}

		MongoSession, err = mgo.DialWithInfo(dialInfo)

	}

	if err != nil {
		log.Fatal("Failed to start the Mongo session")
	}

	//load database name
	if MongoDatabase == nil {
		SetDatabaseMongoDB(LoadConfigMongoDB.Name)
	}

}

/*SetDatabaseMongoDB set database*/
func SetDatabaseMongoDB(dbName string) {
	MongoDatabase = MongoSession.DB(dbName)
}

/*GetMongoSessionClone connect clone Session*/
func GetMongoSessionClone() *mgo.Session {
	return MongoSession.Clone()
}

/*GetMongoSessionCopy copy session add more resource connect*/
func GetMongoSessionCopy() *mgo.Session {
	return MongoSession.Copy()
}

/*GetMongoCollection get current collection*/
func GetMongoCollection(name string) *mgo.Collection {
	return MongoDatabase.C(name)
}

/*MongodbToBson convert to query to bson*/
func MongodbToBson(query interface{}) (bsonQ interface{}) {
	byteData, _ := bson.Marshal(query)
	bson.Unmarshal(byteData, &bsonQ)
	return bsonQ
}
