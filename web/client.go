package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func request(body *ApiBody, w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error

	switch body.Method {
	case http.MethodGet:
		req, _ := http.NewRequest("GET", body.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("error : %v!\n", err)
			return
		}
	case http.MethodPost:
		req, _ := http.NewRequest("POST", body.Url, bytes.NewBuffer([]byte(body.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("error : %v!\n", err)
			return
		}
	case http.MethodDelete:
		req, _ := http.NewRequest("DELETE", body.Url, nil)
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("error : %v!\n", err)
			return
		}
	case http.MethodPut:
		req, _ := http.NewRequest("PUT", body.Url, bytes.NewBuffer([]byte(body.ReqBody)))
		req.Header = r.Header
		resp, err = httpClient.Do(req)
		if err != nil {
			log.Printf("error : %v!\n", err)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Bad api request")
		return
	}
	response(w, resp)
}

func response(w http.ResponseWriter, r *http.Response) {
	res, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rec, _ := json.Marshal(ErrorInternalFaults)
		w.WriteHeader(500)
		io.WriteString(w, string(rec))
		return
	}

	w.WriteHeader(r.StatusCode)
	io.WriteString(w, string(res))
}
