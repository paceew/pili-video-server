package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/streamserver/def"
	"github.com/pili-video-server/streamserver/utils"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, err := template.ParseFiles("./streamserver/upload.html")
	if err != nil {
		log.Printf("parse files error:%v!\n", err)
		return
	}
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid_id")
	var vl string
	mat := p.ByName("mat")
	switch mat {
	case "org":
		vl = def.VIDEO_DIR_ORG + vid
	case "720":
		vl = def.VIDEO_DIR_720 + vid + def.VIDEO_FROMAT
	case "480":
		vl = def.VIDEO_DIR_480 + vid + def.VIDEO_FROMAT
	case "360":
		vl = def.VIDEO_DIR_360 + vid + def.VIDEO_FROMAT
	}
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

	//上传视频
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
	err = ioutil.WriteFile(def.VIDEO_DIR_ORG+vid, data, 0666)
	if err != nil {
		log.Printf("io write error : %v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	//上传封面
	icon, _, err := r.FormFile("icon")
	if err != nil {
		log.Printf("icon error :%v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	dataIc, err := ioutil.ReadAll(icon)
	if err != nil {
		log.Printf("dataIc error :%v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	err = ioutil.WriteFile(def.ICON_DIR+vid, dataIc, 0666)
	if err != nil {
		log.Printf("io write error : %v\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	//发送转码请求到scheduler
	go utils.SendFormatVideoRequest(vid)
	sendNormalResponse(w, "upload ok !", 201)
}

func GetIcon(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid_id")
	il := def.ICON_DIR + vid
	icon, err := os.Open(il)
	if err != nil {
		log.Printf("open icon error:%v!\n", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	w.Header().Set("Conten-Type", "image/png")
	http.ServeContent(w, r, "", time.Now(), icon)

	defer icon.Close()
}
