package jsonserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"../util"
)

// CalcJSONMap rp.ParseForm() return rp.Form
func CalcJSONMap(rp *http.Request) map[string]interface{} {
	rp.ParseForm()

	for k, v := range rp.Form {
		fmt.Println(k)
		fmt.Println(v)

		b := []byte(k)
		var f interface{}
		json.Unmarshal(b, &f)

		m := f.(map[string]interface{})

		return m
	}
	return nil
}

// ConvertJSON : do ParseForm, and json.Marshal(rp.Form).
// if error is occured then panic.
func ConvertJSON(rp *http.Request) []byte {
	js := CalcJSONMap(rp)
	buf, err := json.Marshal(js)

	util.IfnilError(err)

	return buf
}

// JSONEchoResponseFunc is writing response method.
// This is echo back method.
// (*wp).Write(json.Marshal(rp.Form))
func JSONEchoResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	buf := ConvertJSON(rp)
	(*wp).Write(buf)
}
