package basicserver

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ResponseFunc func(wp *http.ResponseWriter, rp *http.Request)

type Page struct {
	Title string
}

func DefaultResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	fmt.Println("DEFAULT")

	files, err := ioutil.ReadDir("static_res")
	file_arr_str := ""
	for i := 0; i < len(files); i++ {
		file_arr_str += "\"" + files[i].Name() + "\","
	}

	page := Page{"DefaultPage"}
	tmpl, err := template.ParseFiles("default.html")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(*wp, page)
	if err != nil {
		panic(err)
	}
}

func SetStaticResponseFunc(wp *http.ResponseWriter, rp *http.Request) {
	fmt.Println("StaticRes")
	fname := rp.URL.Path[1:len(rp.URL.Path)]
	if strings.HasSuffix(rp.URL.Path, ".css") {
		(*wp).Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(rp.URL.Path, ".js") {
		(*wp).Header().Set("Content-Type", "text/javascript")
	} else if strings.HasSuffix(rp.URL.Path, ".ts") {
		(*wp).Header().Set("Content-Type", "application/octet-stream")
	} else if !strings.Contains(rp.URL.Path, ".") {
		(*wp).Header().Set("Content-Type", "text/html")
		//		fname += ".html" // .htmは救えない
	}

	fp, err := os.Open(fname)
	if err != nil {
		curDir, _ := filepath.Abs(".")
		fmt.Println(curDir)
		fmt.Println("os.OpenError " + fname)
		panic(err)
	}
	defer fp.Close()

	buf := make([]byte, 1024*1024*10)
	for {
		n, _ := fp.Read(buf)
		if n == 0 {
			break
		}
		(*wp).Write(buf[:n])
	}
}
func CheckAuth(r *http.Request, inputusername string, inputpassword string) bool {
	username, password, ok := r.BasicAuth()
	if ok == false {
		return false
	}
	return username == inputusername && password == inputpassword
}

func ViewHandler(
	w *http.ResponseWriter,
	r *http.Request,
	checkAuthFunc func(*http.Request) bool,
	authFunc func(*http.ResponseWriter),
	responseFuncMap map[string]ResponseFunc,
	defaultFunc ResponseFunc,
) {
	if checkAuthFunc(r) == false {
		authFunc(w)
		return
	}
	//	r.ParseForm()

	url_bytes := []byte(r.URL.Path)
	match_flag := false
	for k, f := range responseFuncMap {
		matched, err := regexp.Match(k, url_bytes)
		if err != nil {
			panic(err)
		}
		if matched {
			// debug print
			fmt.Println("match")
			fmt.Println(k)
			fmt.Println(r.URL.Path)

			// url to function
			f(w, r)
			match_flag = true
			break
		}
	}
	if !match_flag {
		defaultFunc(w, r)
	}
}
