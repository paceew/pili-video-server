package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pace/sample/streamserver/def"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, err := template.ParseFiles("/home/pace/go/src/github.com/pace/sample/streamserver/upload.html")
	if err != nil {
		log.Printf("parse files error:%v!\n", err)
		return
	}
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid_id")
	vl := def.VIDEO_DIR + vid
	log.Printf("video url: %v\n", vl)
	video, err := os.Open(vl)
	if err != nil {
		log.Printf("open video error:%v!\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	w.Header().Set("Conten-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, def.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(def.MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("file error :%v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("data error :%v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	vid := p.ByName("vid_id")
	log.Printf("uploading vid : %v\n", vid)
	err = ioutil.WriteFile(def.VIDEO_DIR+vid, data, 0666)
	if err != nil {
		log.Printf("io write error : %v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	//	utils.Proxy(w, r)                     上传转发
	sendNormalResponse(w, "upload ok !", 201)
}
