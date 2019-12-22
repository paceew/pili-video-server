package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegistHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", homeHandler)
	router.POST("/", homeHandler)

	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)

	//api 透传
	router.POST("/api", apiHandler)

	//proxy 转发
	router.POST("/video/:vid_id", proxyHandler)
	router.GET("/video/:vid_id", proxyHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./template"))

	return router
}

func main() {
	r := RegistHandler()
	http.ListenAndServe(":8080", r)
}
