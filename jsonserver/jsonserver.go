package jsonserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Hoge struct {
	A string `json:A`
	B string `json:B`
}

func JsonEchoResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	rp.ParseForm()
	fmt.Println(rp.Form)

	buf, err := json.Marshal(rp.Form)
	(*wp).Write(buf)

	if err != nil {
		fmt.Println(err)
	}
}
