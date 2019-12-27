package main

import (
	"log"
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
	//check session
	log.Printf("middleWareHandlerOne!\n")
	validateUserSession(r)
	//允许访问所有域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	//users handeler
	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)
	router.GET("/user/:user_name", GetUserInfo)
	router.DELETE("/user/:user_name", Logout)
	router.PUT("/user/:user_name/pwd/modify", ModifyPwd)
	router.PUT("/user/:user_name", ModifyUserInfo)

	//videos handler
	router.GET("/user/:user_name/videos/:page", ListAllVideosByUser)
	router.DELETE("/user/:user_name/video/:vid_id", DeleteVideo)
	router.POST("/user/:user_name/video", AddNewVideo)
	// router.POST("/user/:user_name/video/:vid_id/itd", AddIntroduction)
	router.GET("/user/:user_name/video/:vid_id", GetVideo)
	router.GET("/user/:user_name/video/:vid_id/itd", GetIntroduction)
	router.GET("/video/:modular/tim/:page", ListAllVideosByModTim)
	router.GET("/video/:modular/hot/:page", ListAllVideosByModHot)

	//videos like handler
	router.POST("/archive/video/:vid_id/like", LikeVideo)
	router.GET("/archive/video/:vid_id/like", LikeCount)
	router.GET("/archive/video/:vid_id/islike", IsLike)

	//comments handler
	router.GET("/videos/:vid_id/comments", ListComments)
	router.POST("/videos/:vid_id/comments", AddNewComment)
	router.DELETE("/videos/:vid_id/comments/:com_id", DeleteComment)

	//messages handler
	router.GET("/user/:user_name/mess_num", GetUnreadMessages)
	router.GET("/user/:user_name/mess", ListUserMessages)
	router.GET("/user/:user_name/mess/:friend_name", GetUserMessage)
	router.POST("/user/:user_name/mess/:friend_name", SendUserMessage)
	router.DELETE("/user/:user_name/mess/:friend_name", DeleteMessages)

	return router
}

func main() {

	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)

}
