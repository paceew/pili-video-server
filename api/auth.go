package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pace/sample/api/dbops"
	"github.com/pace/sample/api/def"
	"github.com/pace/sample/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_USERNAME = "X-User-Name"

func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}

	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	log.Printf("uname:%v\n", uname)
	session.ReSetSession(sid, uname)

	r.Header.Add(HEADER_FIELD_USERNAME, uname)
	log.Printf("Header uname : %v", r.Header.Get(HEADER_FIELD_USERNAME))
	return true
}

func ValidateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	pname := p.ByName("user_name")
	if len(uname) == 0 || len(pname) == 0 || pname != uname {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	return true
}

func ValidateUserPwd(w http.ResponseWriter, pwd string, username string) bool {
	upwd, err := dbops.GetUserCredential(username)
	if err != nil || len(upwd) == 0 || pwd != upwd {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return false
	}

	return true
}

func ValidateLogin(w http.ResponseWriter, r *http.Request) bool {
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	if len(uname) == 0 {
		sendErrorResponse(w, def.ErrorNotAuthUser)
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
