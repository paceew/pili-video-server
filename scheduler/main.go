package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/scheduler/taskrunner"
	"github.com/pili-video-server/scheduler/timertask"
)

func RegistHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video_del/:vid_id", AddDelection)

	router.GET("/video_format/:vid_id", AddUnFormat)

	return router
}
func main() {
	go taskrunner.Start()
	go timertask.Start()
	r := RegistHandler()
	http.ListenAndServe(":9001", r)
}
