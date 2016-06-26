package dbmgr

import (
	"net/http"

	"../jsonserver"
	"../mongoutil"
	"gopkg.in/mgo.v2"
)

type dbController struct {
	MDBElem map[string]map[string]*mgo.Collection
}

// SessionMgr http+json+mongodb.
// require Initialize
type SessionMgr struct {
	mpSession      *mgo.Session
	mpDBController *dbController
}

// Initialize session create
func (s *SessionMgr) Initialize(mongoip string) error {
	s.mpSession = mongoutil.CreateSession(mongoip)
	return nil
}

func (s *SessionMgr) getCollection(dbname, collectionname string) *mgo.Collection {
	if mpDBController.MDBElem[dbname][collectionname] == nil {
		mpDBController.MDBElem[dbname][collectionname] = s.mpSession.DB(dbname).C(collectionname)
	}
	return mpDBController.MDBElem[dbname][collectionname]
}

func mapToString(m map[string]interface{}, key string) string {
	return m[key].(string)
}

// GetFindQuery get db, collection and query from httprequest
func GetFindQuery(rp *http.Request) (string, string, interface{}) {
	js := jsonserver.CalcJSONMap(rp)

	dbName := mapToString(js, "db")
	collectionName := mapToString(js, "collection")
	query := js["query"]

	return dbName, collectionName, query
}

// GetInsertQuery get db, collection and query from httprequest
func GetInsertQuery(rp *http.Request) (string, string, interface{}) {
	return GetFindQuery(rp)
}

// GetUpdateQuery get db, collection, findquery and updatequery from httprequest
func GetUpdateQuery(rp *http.Request) (string, string, interface{}, interface{}) {
	js := jsonserver.CalcJSONMap(rp)

	dbName := mapToString(js, "db")
	collectionName := mapToString(js, "collection")
	findquery := js["findquery"]
	updatequery := js["updatequery"]
	return dbName, collectionName, findquery, updatequery
}

// Find find method.
func (s *SessionMgr) Find(dbname, collectionname string, query interface{}) []byte {
	collection := s.getCollection(dbname, collectionname)
	return mongoutil.Find(collection, query)
}

// Insert insert method.
func (s *SessionMgr) Insert(dbname, collectionname string, query interface{}) {
	collection := s.getCollection(dbname, collectionname)
	mongoutil.Insert(collection, query)
}

// Update update method.
func (s *SessionMgr) Update(dbname, collectionname string, findquery interface{}, updatequery interface{}) {
	collection := s.getCollection(dbname, collectionname)
	mongoutil.Update(collection, findquery, updatequery)
}
