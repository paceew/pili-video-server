package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/streamserver/def"
)

type middleWareHandler struct {
	r  *httprouter.Router
	cl *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//middle ware handler
	if !m.cl.GetConn() {
		sendErrorResponse(w, def.ErrorConnectLimit)
		return
	}

	//允许访问所有域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	m.r.ServeHTTP(w, r)
	//releaseConnect
	defer m.cl.ReleaseConn()

}

func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.cl = NewConnLimiter(cc)
	return m
}

func RegisterHandler() *httprouter.Router {
	router := httprouter.New()

	router.GET("/test", testPageHandler)

	router.GET("/video/:vid_id/:mat", streamHandler)
	router.POST("/video/:vid_id", uploadHandler)

	router.GET("/icon/:vid_id", GetIcon)

	return router
}

func main() {
	r := RegisterHandler()
	mh := NewMiddleWareHandler(r, 5)
	http.ListenAndServe(":9000", mh)
}
