package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//允许访问所有域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	m.r.ServeHTTP(w, r)
}

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
	router.GET("/icon/:vid_id", proxyHandler)

	router.ServeFiles("/statics/*filepath", http.Dir("./template"))

	return router
}

func main() {
	r := RegistHandler()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8080", mh)
}
