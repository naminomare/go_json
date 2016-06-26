package main

import (
	"fmt"
	"net/http"

	"./basicserver"
	"./dbmgr"
	"./jsonserver"
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

var sDBMgr dbmgr.SessionMgr

func insertResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	dbname, collectionname, query := dbmgr.GetInsertQuery(rp)
	sDBMgr.Insert(dbname, collectionname, query)
	(*wp).Write([]byte("success"))
}

func findResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	dbname, collectionname, query := dbmgr.GetFindQuery(rp)
	res := sDBMgr.Find(dbname, collectionname, query)
	(*wp).Write(res)
}

func updateResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	dbname, collectionname, findquery, updatequery := dbmgr.GetUpdateQuery(rp)
	sDBMgr.Update(dbname, collectionname, findquery, updatequery)
	(*wp).Write([]byte("success"))
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
	// sSession = mongoutil.CreateSession(cMONGOIP)
	//	sTestCollection = sSession.DB(cMONGODBNAME).C(cMONGOCOLLECTIONNAME)

	sDBMgr.Initialize(cMONGOIP)

	http.HandleFunc("/", viewHandler)
	http.ListenAndServeTLS(":"+cPORT, "cert.pem", "key.pem", nil)
}
