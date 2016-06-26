package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"./basicserver"
	"./jsonserver"
	"./mongoutil"
)

const (
	cUSERNAME = "user"
	cPASSWORD = "password"
	cPORT     = "20000"
	cMONGOIP  = "localhost"

	cMONGODBNAME         = "anim"
	cMONGOCOLLECTIONNAME = "fileinfo"
)

var sResponseFuncMap = map[string]basicserver.ResponseFunc{
	"^/static_res/.*": basicserver.SetStaticResponseFunc,
	"^/json_test":     jsonserver.JSONEchoResponseFunc,
	"^/insert":        insertResponseFunc,
	"^/find":          findResponseFunc,
	"^/update":        updateResponseFunc,
}

//var sDBMgr dbmgr.DBMgr

var sSession *mgo.Session

//var sTestCollection *mgo.Collection
var sCollectionMap map[string]*mgo.Collection

func getCollection(collectionName string) *mgo.Collection {
	if sCollectionMap[collectionName] == nil {
		sCollectionMap[collectionName] = sSession.DB(cMONGODBNAME).C(collectionName)
	}
	return sCollectionMap[collectionName]
}

func insertResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	js := jsonserver.CalcJSONMap(rp)

	collectionNameI := js["collection"]
	insertQuery := js["query"]
	collectionName := collectionNameI.(string)
	collection := getCollection(collectionName)

	mongoutil.Insert(collection, insertQuery)
	(*wp).Write([]byte("success"))
}

func findResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	js := jsonserver.CalcJSONMap(rp)

	collectionNameI := js["collection"]
	findquery := js["query"]
	collectionName := collectionNameI.(string)
	collection := getCollection(collectionName)

	res := mongoutil.Find(collection, findquery)
	(*wp).Write(res)
}

func updateResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	m := jsonserver.CalcJSONMap(rp)
	var collectionNameI, findquery, updatequery interface{}

	collectionNameI = m["collection"]
	findquery = m["findquery"]
	updatequery = m["updatequery"]
	collectionName := collectionNameI.(string)

	collection := getCollection(collectionName)

	if findquery != nil {
		mongoutil.Update(collection, findquery, updatequery)
		(*wp).Write([]byte("success"))
	} else {
		(*wp).Write([]byte("failure"))
	}
}

func authFunc(w *http.ResponseWriter) {
	(*w).Header().Set("WWW-Authenticate", "Basic realm='input your id and password")
	(*w).WriteHeader(401)
	(*w).Write([]byte("401 Unauthorized\n"))
}

func checkAuthFunc(r *http.Request) bool {
	return basicserver.CheckAuth(r, cUSERNAME, cPASSWORD)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	basicserver.ViewHandler(
		&w,
		r,
		checkAuthFunc,
		authFunc,
		sResponseFuncMap,
		basicserver.DefaultResponseFunc,
	)
}

func main() {
	fmt.Println("Hello World.")
	sSession = mongoutil.CreateSession(cMONGOIP)
	//	sTestCollection = sSession.DB(cMONGODBNAME).C(cMONGOCOLLECTIONNAME)

	//sDBMgr.Initialize(cMONGOIP)

	http.HandleFunc("/", viewHandler)
	http.ListenAndServeTLS(":"+cPORT, "cert.pem", "key.pem", nil)
}
