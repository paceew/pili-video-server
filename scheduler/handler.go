package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/scheduler/dbops"
)

func AddDelection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid_id")
	err := dbops.AddDeletion(vid)
	if err != nil {
		log.Printf("add deletion error:%v!\n", err)
		w.WriteHeader(400)
		resp, _ := json.Marshal("db error!")
		io.WriteString(w, string(resp))
		return
	}

	w.WriteHeader(200)
	resp, _ := json.Marshal("ok!")
	io.WriteString(w, string(resp))

}
