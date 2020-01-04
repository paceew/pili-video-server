package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/api/dbops"
	"github.com/pili-video-server/api/def"
	"github.com/pili-video-server/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_USERNAME = "X-User-Name"

//验证session如果存在则重设session并且设置requset header X-User-Name
func validateUserSession(r *http.Request) bool {
	// log.Printf("validating User Session ...\n")
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		r.Header.Add(HEADER_FIELD_USERNAME, "")
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		log.Printf("session expired!")
		r.Header.Add(HEADER_FIELD_USERNAME, "")
		return false
	}
	log.Printf("uname:%v\n", uname)
	session.ReSetSession(sid, uname)

	r.Header.Add(HEADER_FIELD_USERNAME, uname)
	log.Printf("Header uname : %v", r.Header.Get(HEADER_FIELD_USERNAME))
	return true
}

//根据params的uname和request header 的uname进行验证
func ValidateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	pname := p.ByName("user_name")
	log.Printf("header name : %v,pname: %v!\n", uname, pname)
	if len(uname) == 0 || len(pname) == 0 || pname != uname {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	return true
}

func VailidateAdmin(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	aname := r.Header.Get(HEADER_FIELD_USERNAME)
	pname := p.ByName("admin_name")
	if len(aname) == 0 || len(pname) == 0 || pname != aname {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	if !dbops.IsAdmin(aname) {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	return true
}

func ValidateUserPwd(w http.ResponseWriter, pwd string, username string) bool {
	upwd, err := dbops.GetUserCredential(username)
	if err != nil || len(upwd) == 0 || pwd != upwd {
		sendErrorResponse(w, def.ErrorUserPassword)
		return false
	}

	return true
}

func ValidateLogin(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, def.ErrorNotLogin)
		return false
	}

	return true
}

func ValidateVideoAnthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	vid := p.ByName("vid_id")
	pname, _ := dbops.GetUserNameByVid(vid)
	if len(uname) == 0 || len(pname) == 0 || pname != uname {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	return true
}

func ValidateCommentAnthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	cid := p.ByName("com_id")
	pname, _ := dbops.GetUserNameByCid(cid)
	if len(uname) == 0 || len(pname) == 0 || pname != uname {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	return true
}
