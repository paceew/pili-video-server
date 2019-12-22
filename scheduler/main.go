package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/scheduler/taskrunner"
)

func RegistHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/video_del/:vid_id", AddDelection)

	return router
}
func main() {
	go taskrunner.Start()
	r := RegistHandler()
	http.ListenAndServe(":9001", r)
}
