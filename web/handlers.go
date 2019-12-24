package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

func homeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	sid, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		page := &HomePage{Name: cname.Value}
		t, err := template.ParseFiles("./template/home.html")
		if err != nil {
			log.Printf("parse files error!\n")
			return
		}

		t.Execute(w, page)
		return
	}

	if len(cname.Value) != 0 && len(sid.Value) != 0 {
		http.Redirect(w, r, "/userhome", http.StatusFound)
		return
	}

}

func userHomeHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cname, err1 := r.Cookie("username")
	_, err2 := r.Cookie("session")
	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	fname := r.FormValue("username")

	var page *UserPage
	if len(cname.Value) != 0 {
		page = &UserPage{Name: cname.Value}
	} else if len(fname) != 0 {
		page = &UserPage{Name: fname}
	}

	t, err := template.ParseFiles("./template/userhome.html")
	if err != nil {
		log.Printf("pase file error!\n")
		return
	}
	t.Execute(w, page)
}

func apiHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Method != http.MethodPost {
		rec, _ := json.Marshal(ErrorRequestNotRecognized)
		io.WriteString(w, string(rec))
		return
	}

	res, _ := ioutil.ReadAll(r.Body)
	apibody := &ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		rec, _ := json.Marshal(ErrorRequestBodyParseFailed)
		io.WriteString(w, string(rec))
		return
	}

	request(apibody, w, r)
	defer r.Body.Close()
}

func proxyHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url, _ := url.Parse("http://localhost:9000/")
	proxy := httputil.NewSingleHostReverseProxy(url)
	log.Printf("proxy serveHTTP 9000 !\n")
	proxy.ServeHTTP(w, r)
}
