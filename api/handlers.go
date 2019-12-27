package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pili-video-server/api/dbops"
	"github.com/pili-video-server/api/def"
	"github.com/pili-video-server/api/session"
	"github.com/pili-video-server/api/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	if err := dbops.AddUser(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	sup := &def.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(sup); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.UserCredential{}

	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	uname := p.ByName("user_name")
	log.Printf("url name: %v", uname)
	log.Printf("request name : %v", ubody.Username)
	if uname != ubody.Username {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return
	}

	if !ValidateUserPwd(w, ubody.Pwd, ubody.Username) {
		log.Printf("pass word error!")
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	sup := &def.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(sup); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 201)
	}

}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户是否登陆
	if !ValidateLogin(w, r) {
		log.Printf("Unauthorized user\n")
		return
	}

	uname := p.ByName("user_name")
	user, err := dbops.GetUser(uname)
	if err != nil {
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	userInfo := &def.UserInfo{Id: user.Id, Pwd: user.Pwd, Name: user.Username}
	if resp, err := json.Marshal(userInfo); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}

func Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}
	session.DeleteExpiredSession(sid)
	sendNormalResponse(w, "Logout ok !", 200)
	//	io.WriteString(w, "user logout!")
}

func ModifyPwd(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户
	if !ValidateUser(w, r, p) {
		log.Printf("Unauthorized user\n")
		return
	}

	uname := p.ByName("user_name")
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &def.UserModifyPwd{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("unmarshal error!")
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	if !ValidateUserPwd(w, ubody.PTPwd, uname) {
		log.Printf("pass word error!")
		return
	}

	if err := dbops.ModifyUserPwd(uname, ubody.CPwd); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	}

	if resp, err := json.Marshal(ubody); err != nil {
		log.Printf("marshal error!")
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

	//	io.WriteString(w, "Modify password!")
}

func ModifyUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// if !ValidateUser(w, r) {
	// 	log.Printf("Unauthorized user\n")
	// 	return
	// }

	// uname := p.ByName("user_name")
	// res, _ := ioutil.ReadAll(r.Body)
	// ubody := &def.UserModifyInfo{}
	// if err := json.Unmarshal(res, ubody); err != nil {
	// 	sendErrorResponse(w, def.ErrorInternalFaults)
	// 	return
	// }

	// if err := dbops.ModifyUserInfo(uname, ubody.CPwd); err != nil {
	// 	sendErrorResponse(w, def.ErrorInternalFaults)
	// 	return
	// }

	// if resp, err := json.Marshal(ubody); err != nil {
	// 	sendErrorResponse(w, def.ErrorInternalFaults)
	// 	return
	// } else {
	// 	sendNormalResponse(w, string(resp), 200)
	// }

	io.WriteString(w, "modify user info!")
}

func ListAllVideosByUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户是否登陆
	// if !ValidateLogin(w, r) {
	// 	log.Printf("Unauthorized user\n")
	// 	return
	// }
	fromstr := p.ByName("page")
	fromint, err := strconv.Atoi(fromstr)
	if err != nil {
		log.Print("fromstr error:%v!\n", err)
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}

	uname := p.ByName("user_name")
	videoList, err := dbops.ListVideoInfo(uname, fromint, def.PAGE_NUM)
	if err != nil {
		log.Print("list video db error:%v!\n", err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	videos := &def.VideosList{Videos: videoList}
	if resp, err := json.Marshal(videos); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
	//	io.WriteString(w, "List all videos of:"+uname)
}

func ListAllVideosByModTim(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromstr := p.ByName("page")
	fromint, err := strconv.Atoi(fromstr)
	if err != nil {
		log.Print("fromstr error:%v!\n", err)
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}

	mod := p.ByName("modular")
	videoList, err := dbops.ListVideoInfoMod(mod, fromint, def.PAGE_NUM, "time")
	if err != nil {
		log.Print("list video db error!\n")
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	videos := &def.VideosList{Videos: videoList}
	if resp, err := json.Marshal(videos); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func ListAllVideosByModHot(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromstr := p.ByName("page")
	fromint, err := strconv.Atoi(fromstr)
	if err != nil {
		log.Print("fromstr error:%v!\n", err)
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}

	mod := p.ByName("modular")
	videoList, err := dbops.ListVideoInfoMod(mod, fromint, def.PAGE_NUM, "hot")
	if err != nil {
		log.Print("list video db error!\n")
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	videos := &def.VideosList{Videos: videoList}
	if resp, err := json.Marshal(videos); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func GetVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户是否登陆
	// if !ValidateLogin(w, r) {
	// 	log.Printf("Unauthorized user\n")
	// 	return
	// }

	vid := p.ByName("vid_id")
	res, err := dbops.GetVideoInfo(vid)
	if err != nil {
		log.Printf("get video info db error!\n")
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	video_info := &def.VideoInfo{Id: res.Id, Name: res.Name, DisplayCtime: res.DisplayCtime, AuthorId: res.AuthorId}
	if resp, err := json.Marshal(video_info); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}

}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户
	if !ValidateUser(w, r, p) {
		log.Printf("Unauthorized user\n")
		return
	}

	vid := p.ByName("vid_id")
	err := dbops.DeleteVideoInfo(vid)
	log.Printf("vid : %v!\n", vid)
	if err != nil {
		log.Printf("delete db error: %v!\n", err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "", 204)
	//	io.WriteString(w, "Delete a video: "+vid)
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户是否登陆
	if !ValidateLogin(w, r) {
		log.Printf("user not login\n")
		return
	}

	video := &def.NewVideo{}
	res, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(res, video); err != nil {
		log.Print("json unmarshal video info error!\n")
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}

	log.Printf("video name:%v", video.VideoName) ///////////////////////
	aid, err := dbops.GetUserId(video.AuthorName)
	if err != nil {
		log.Printf("aid : %v get user id error : %v\n", aid, err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	vInfo := &def.VideoInfo{}
	vInfo, err = dbops.AddNewVideo(aid, video.VideoName)
	if err != nil {
		log.Printf("aid : %v add new video error:%v\n", aid, err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	} else {
		if resp, err := json.Marshal(vInfo); err != nil {
			sendErrorResponse(w, def.ErrorInternalFaults)
			return
		} else {
			sendNormalResponse(w, string(resp), 200)
		}
	}

	//	io.WriteString(w, "Add a new video: "+vid)
}

func LikeVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证是否登陆
	if !ValidateLogin(w, r) {
		log.Printf("user not longin\n")
		return
	}
	vid := p.ByName("vid_id")
	uname := r.Header.Get(HEADER_FIELD_USERNAME)
	err := dbops.LikeVideo(vid, uname)
	if err != nil {
		log.Printf("redis error :%v", err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	sendNormalResponse(w, "", 204)

}

func LikeCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid_id")

	count, err := dbops.LikeCount(vid)
	if err != nil {
		log.Printf("someting error :%v", err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	if res, err := json.Marshal(count); err != nil {
		log.Printf("json marshal error :%v", err)
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(res), 204)
		return
	}

}

//已点赞返回1，未点赞返回0
func IsLike(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid_id")
	uname := r.Header.Get(HEADER_FIELD_USERNAME)

	yes, err := dbops.IsLike(vid, uname)
	if err != nil {
		log.Printf("someting error :%v\n", err)
		sendErrorResponse(w, def.ErrorDBError)
		return
	}
	if yes {
		sendNormalResponse(w, "1", 200)
	} else {
		sendNormalResponse(w, "0", 200)
	}
}

func ListComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromstr := p.ByName("page")
	fromint, err := strconv.Atoi(fromstr)
	if err != nil {
		log.Print("fromstr error:%v!\n", err)
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}
	vid := p.ByName("vid_id")

	commentInfo, err := dbops.ListComments(vid, fromint, def.PAGE_NUM)
	if err != nil {
		log.Printf("list comment db error!\n")
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	commentList := &def.CommentsList{Comments: commentInfo}
	if resp, err := json.Marshal(commentList); err != nil {
		sendErrorResponse(w, def.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), 200)
	}
}

func AddNewComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户是否登陆
	if !ValidateLogin(w, r) {
		log.Printf("user not longin\n")
		return
	}

	vid := p.ByName("vid_id")
	res, _ := ioutil.ReadAll(r.Body)
	comment := &def.NewComment{}
	if err := json.Unmarshal(res, comment); err != nil {
		log.Printf("add new comment unmarshal error!\n")
		sendErrorResponse(w, def.ErrorRequestBodyPaseFailed)
		return
	}

	aid, err := dbops.GetUserId(comment.AuthorName)
	if err != nil {
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	err = dbops.AddNewComment(aid, vid, comment.Content)
	if err != nil {
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	sendNormalResponse(w, "ok", 200)
}

func DeleteComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//验证用户权限
	if !ValidateVideoAnthor(w, r, p) && !ValidateCommentAnthor(w, r, p) {
		log.Printf("Unauthorized user\n")
		return
	}

	cid := p.ByName("com_id")
	err := dbops.DeleteComment(cid)
	if err != nil {
		sendErrorResponse(w, def.ErrorDBError)
		return
	}

	sendNormalResponse(w, "", 204)
}

func GetUnreadMessages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return
}

func ListUserMessages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return
}

func GetUserMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return
}

func SendUserMessage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return
}

func DeleteMessages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return
}
