package main

import (
	"fmt"
	"net/http"

	"./basicserver"
	"./jsonserver"
)

const (
	cUSERNAME = "user"
	cPASSWORD = "password"
	cPORT     = "20000"
)

var responsefuncmap = map[string]basicserver.ResponseFunc{
	"^/static_res/.*": basicserver.SetStaticResponseFunc,
	"^/json_test":     jsonserver.JsonEchoResponseFunc,
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
		responsefuncmap,
		basicserver.DefaultResponseFunc,
	)
}

func main() {
	fmt.Println("Hello World.")

	http.HandleFunc("/", viewHandler)
	http.ListenAndServeTLS(":"+cPORT, "cert.pem", "key.pem", nil)
}
