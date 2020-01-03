package dbops

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/pili-video-server/api/def"
	"github.com/pili-video-server/api/utils"
)

func AddUser(userName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (username, pwd) VALUES (?, ?)")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmtIns.Exec(userName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(userName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE username = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(userName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	defer stmtOut.Close()
	return pwd, nil
}

func GetUser(userName string) (*def.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE username = ?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	var id int
	var pwd string
	err = stmtOut.QueryRow(userName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	result := &def.User{Id: id, Pwd: pwd, Username: userName}
	defer stmtOut.Close()
	return result, nil
}

func GetUserId(userName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT id FROM users WHERE username = ?")
	if err != nil {
		log.Printf("get user id db prepare error!\n")
		return "", err
	}

	var id string
	err = stmtOut.QueryRow(userName).Scan(&id)
	if err != nil {
		log.Printf("query row error:%v!\n", err)
		return "", err
	}

	defer stmtOut.Close()
	return id, nil
}

func ModifyUserInfo(userName string) error {
	//没有过多信息
	return nil
}

func ModifyUserPwd(userName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("UPDATE users SET pwd = ? WHERE username = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	_, err = stmtIns.Exec(pwd, userName)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func AddNewVideo(aid string, vname string, mid int, itd string) (string, error) {
	vid, _ := utils.NewUUID()

	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")

	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info_copy
	(id, author_id, name, modular, create_time) VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("insert db prepare error\n")
		return "", err
	}

	_, err = stmtIns.Exec(vid, aid, vname, mid, ctime)
	if err != nil {
		return "", err
	}

	_ = AddIntroduction(vid, itd)

	// vInfo := &def.VideoInfo{Id: vid, AuthorId: aid, Name: vname, DisplayCtime: ctime}

	defer stmtIns.Close()
	return vid, nil
}

func AddIntroduction(vid string, content string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO introduction VALUES (?,?)")
	if err != nil {
		log.Printf("insert db prepare error\n")
		return err
	}

	_, err = stmtIns.Exec(vid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	if err != nil {
		log.Printf("delete db prepare error!\n")
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}

	defer stmtDel.Close()
	return nil
}

//根据vid获取简介
func GetIntroduction(vid string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT content FROM introduction WHERE vid = ?")
	if err != nil {
		log.Printf("get db prepare error!\n")
		return "", err
	}

	var content string
	err = stmtOut.QueryRow(vid).Scan(&content)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		return "", nil
	}

	defer stmtOut.Close()
	return content, nil
}

func GetVideoInfo(vid string) (*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.name, users.nickname, video_info.create_time, modulars.name, like_number, collect_number, comment_number, content  
	 FROM video_info, modulars, users, introduction  WHERE video_info.id = ? AND video_info.modular = modulars.id AND video_info.author_id = users.username AND video_info.id = introduction.vid`)
	if err != nil {
		log.Printf("get db prepare error!\n")
		return nil, err
	}

	var id, name, displayCTime, modular, authorName, icon, urlOriginal, url720p, url480p, url360p, itd string
	var likeNum, collectNum, commentNum int
	err = stmtOut.QueryRow(vid).Scan(&id, &name, &authorName, &displayCTime, &modular, &likeNum, &collectNum, &commentNum, &itd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	icon = def.STREAMSERVER_ADDR + "icon/" + vid
	urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "/org"
	url720p = def.STREAMSERVER_ADDR + "video/" + vid + "/720p"
	url480p = def.STREAMSERVER_ADDR + "video/" + vid + "/480p"
	url360p = def.STREAMSERVER_ADDR + "video/" + vid + "/360p"

	res := &def.VideoInfo{Id: id, Name: name, AuthorName: authorName, DisplayCtime: displayCTime, Modular: modular, LikeNum: likeNum, CollectNum: collectNum, CommentNum: commentNum, Icon: icon, UrlOriginal: urlOriginal, Url720p: url720p, Url480p: url480p, Url360p: url360p, Introduction: itd}
	defer stmtOut.Close()
	return res, nil
}

func VideoSearch(key string, from, n int) ([]*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.name, users.nickname, video_info.create_time, modulars.name, like_number, collect_number, comment_number,content  
	FROM video_info, modulars, users, introduction  WHERE video_info.author_id = users.username AND video_info.modular = modulars.id AND introduction.vid = video_info.id AND video_info.name like ? LIMIT ?,?`)
	if err != nil {
		log.Printf("get db prepare error!\n")
		return nil, err
	}

	keyStr := "%" + key + "%"

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(keyStr, from, n)
	if err != nil {
		log.Printf("query err\n", err)
		return nil, err
	}
	for rows.Next() {
		var vid, vname, displayTime, modular, authorName, icon, urlOriginal, url720p, url480p, url360p, itd string
		var likeNum, collectNum, commentNum int
		err := rows.Scan(&vid, &vname, &authorName, &displayTime, &modular, &likeNum, &collectNum, &commentNum, &itd)
		if err != nil {
			return nil, err
		}
		icon = def.STREAMSERVER_ADDR + "icon/" + vid
		urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "/org"
		url720p = def.STREAMSERVER_ADDR + "video/" + vid + "/720p"
		url480p = def.STREAMSERVER_ADDR + "video/" + vid + "/480p"
		url360p = def.STREAMSERVER_ADDR + "video/" + vid + "/360p"
		log.Printf("icon:%v\n", icon)

		videoInfo := &def.VideoInfo{Id: vid, AuthorName: authorName, Name: vname, DisplayCtime: displayTime, Modular: modular, LikeNum: likeNum, CollectNum: collectNum, CommentNum: commentNum, Icon: icon, UrlOriginal: urlOriginal, Url720p: url720p, Url480p: url480p, Url360p: url360p, Introduction: itd}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func ListVideoInfo(uname string, from, n int) ([]*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, users.nickname, video_info.name, video_info.create_time, modulars.name, like_number, collect_number, comment_number, content
	FROM video_info,modulars,users, introduction
	WHERE users.username=? AND video_info.modular = modulars.id AND users.username = video_info.author_id AND introduction.vid = video_info.id
	ORDER BY video_info.create_time DESC LIMIT ?,?`)
	if err != nil {
		log.Printf("list db prepare error!\n")
		return nil, err
	}

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(uname, from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid, vname, displayTime, modular, authorName, icon, urlOriginal, url720p, url480p, url360p, itd string
		var likeNum, collectNum, commentNum int
		err := rows.Scan(&vid, &authorName, &vname, &displayTime, &modular, &likeNum, &collectNum, &commentNum, &itd)
		if err != nil {
			return nil, err
		}
		icon = def.STREAMSERVER_ADDR + "icon/" + vid
		urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "/org"
		url720p = def.STREAMSERVER_ADDR + "video/" + vid + "/720p"
		url480p = def.STREAMSERVER_ADDR + "video/" + vid + "/480p"
		url360p = def.STREAMSERVER_ADDR + "video/" + vid + "/360p"

		videoInfo := &def.VideoInfo{Id: vid, AuthorName: authorName, Name: vname, DisplayCtime: displayTime, Modular: modular, LikeNum: likeNum, CollectNum: collectNum, CommentNum: commentNum, Icon: icon, UrlOriginal: urlOriginal, Url720p: url720p, Url480p: url480p, Url360p: url360p, Introduction: itd}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func ListVideoInfoNoExam(uname string, from, n int) ([]*def.VideoInfoNoExam, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info_copy.id, users.nickname, video_info_copy.name, video_info_copy.create_time, modulars.name, content, status
	FROM video_info_copy,modulars,users, introduction
	WHERE users.username=? AND video_info_copy.modular = modulars.id AND users.username = video_info_copy.author_id AND introduction.vid = video_info_copy.id
	ORDER BY video_info_copy.create_time DESC LIMIT ?,?`)
	if err != nil {
		log.Printf("list db prepare error!\n")
		return nil, err
	}

	var res []*def.VideoInfoNoExam
	rows, err := stmtOut.Query(uname, from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid, vname, displayTime, modular, authorName, icon, urlOriginal, itd string
		var status int
		err := rows.Scan(&vid, &authorName, &vname, &displayTime, &modular, &itd, &status)
		if err != nil {
			return nil, err
		}
		icon = def.STREAMSERVER_ADDR + "icon/" + vid
		urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "/org"

		videoInfo := &def.VideoInfoNoExam{Id: vid, AuthorName: authorName, Name: vname, DisplayCtime: displayTime, Modular: modular, Icon: icon, UrlOriginal: urlOriginal, Introduction: itd, Status: status}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func RankVideoInfo(from, n int) ([]*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, users.nickname, video_info.name, video_info.create_time, modulars.name, like_number, collect_number, comment_number, content
	FROM video_info,modulars,users,introduction
	WHERE  video_info.modular = modulars.id AND introduction.vid = video_info.id AND users.username = video_info.author_id
	ORDER BY video_info.hot desc LIMIT ?,?`)
	if err != nil {
		log.Printf("list db prepare error!\n")
		return nil, err
	}

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid, vname, displayTime, modular, authorName, icon, urlOriginal, url720p, url480p, url360p, itd string
		var likeNum, collectNum, commentNum int
		err := rows.Scan(&vid, &authorName, &vname, &displayTime, &modular, &likeNum, &collectNum, &commentNum, &itd)
		if err != nil {
			return nil, err
		}
		icon = def.STREAMSERVER_ADDR + "icon/" + vid
		urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "/org"
		url720p = def.STREAMSERVER_ADDR + "video/" + vid + "/720p"
		url480p = def.STREAMSERVER_ADDR + "video/" + vid + "/480p"
		url360p = def.STREAMSERVER_ADDR + "video/" + vid + "/360p"
		videoInfo := &def.VideoInfo{Id: vid, AuthorName: authorName, Name: vname, DisplayCtime: displayTime, Modular: modular, LikeNum: likeNum, CollectNum: collectNum, CommentNum: commentNum, Icon: icon, UrlOriginal: urlOriginal, Url720p: url720p, Url480p: url480p, Url360p: url360p, Introduction: itd}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func ListVideoInfoMod(mod string, from, n int, mode string) ([]*def.VideoInfo, error) {
	var stmtOut *sql.Stmt
	var err error
	if mode == "time" {
		stmtOut, err = dbConn.Prepare(`SELECT video_info.id, users.nickname, video_info.name, video_info.create_time, modulars.name, like_number, collect_number, comment_number, content
		FROM video_info,modulars,users, introduction WHERE modulars.name=? AND modulars.id = video_info.modular AND video_info.author_id = users.username AND introduction.vid = video_info.id
		ORDER BY video_info.create_time DESC LIMIT ?,?`)
		if err != nil {
			log.Printf("list db prepare error!\n")
			return nil, err
		}
	} else if mode == "hot" {
		stmtOut, err = dbConn.Prepare(`SELECT video_info.id, users.nickname, video_info.name, video_info.create_time, modulars.name, like_number, collect_number, comment_number, content
		FROM video_info,modulars,users, introduction WHERE modulars.name=? AND modulars.id = video_info.modular AND video_info.author_id = users.username AND introduction.vid = video_info.id
		ORDER BY video_info.hot desc LIMIT ?,?`)
		if err != nil {
			log.Printf("list db prepare error!\n")
			return nil, err
		}
	}
	// stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.disply_ctime, modulars.name
	// FROM video_info,modulars WHERE modulars.name=?
	// ORDER BY video_info.create_time DESC LIMIT ?,?`)
	// if err != nil {
	// 	log.Printf("list db prepare error!\n")
	// 	return nil, err
	// }

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(mod, from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid, vname, displayTime, modular, authorName, icon, urlOriginal, url720p, url480p, url360p, itd string
		var likeNum, collectNum, commentNum int
		err := rows.Scan(&vid, &authorName, &vname, &displayTime, &modular, &likeNum, &collectNum, &commentNum, &itd)
		if err != nil {
			return nil, err
		}
		icon = def.STREAMSERVER_ADDR + "icon/" + vid
		urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "/org"
		url720p = def.STREAMSERVER_ADDR + "video/" + vid + "/720p"
		url480p = def.STREAMSERVER_ADDR + "video/" + vid + "/480p"
		url360p = def.STREAMSERVER_ADDR + "video/" + vid + "/360p"
		videoInfo := &def.VideoInfo{Id: vid, AuthorName: authorName, Name: vname, DisplayCtime: displayTime, Modular: modular, LikeNum: likeNum, CollectNum: collectNum, CommentNum: commentNum, Icon: icon, UrlOriginal: urlOriginal, Url720p: url720p, Url480p: url480p, Url360p: url360p, Introduction: itd}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

//点赞与取消点赞
func LikeVideo(vid string, uname string) error {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis error!\n")
	}
	defer conn.Close()

	// ustr := string(uid)
	likestr := "like_" + vid
	yes, _ := redis.Bool(conn.Do("sismember", likestr, uname))
	// log.Printf("user:%v the video:%v is like?:%v\n", uname, likestr, yes)
	if !yes {
		_, err = conn.Do("sadd", likestr, uname)
		if err != nil {
			return err
		}
	} else {
		_, err = conn.Do("srem", likestr, uname)
		if err != nil {
			return err
		}
	}

	return nil
}

//获取点赞数
func LikeCount(vid string) (int, error) {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis conn error!\n")
	}
	defer conn.Close()

	likestr := "like_" + vid
	res, err := redis.Int(conn.Do("scard", likestr))
	if err != nil {
		return 0, err
	}

	return res, nil

}

//判断用户是否点赞
func IsLike(vid string, uname string) (bool, error) {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis conn error!\n")
	}
	defer conn.Close()

	likestr := "like_" + vid
	yes, _ := redis.Bool(conn.Do("sismember", likestr, uname))

	return yes, nil
}

func AddNewComment(aid string, vid, content string) error {
	cid, _ := utils.NewUUID()
	stmtIns, err := dbConn.Prepare("INSERT INTO comments(id, author_id, video_id, comment) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Printf("add new comment db prepare error!\n")
		return err
	}

	_, err = stmtIns.Exec(cid, aid, vid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ListComments(vid string, from, n int) ([]*def.CommentInfo, error) {
	var res []*def.CommentInfo
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, video_info.name , users.username, comment 
	FROM comments,video_info,users WHERE users.id=video_info.author_id 
	AND comments.video_id=video_info.id AND video_info.id = ?  
	ORDER BY comments.create_time DESC LIMIT ?,?`)
	if err != nil {
		log.Printf("list comments db prepare error!\n")
		return nil, err
	}

	rows, err := stmtOut.Query(vid, from, n)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var cid, vname, uname, content string
		err := rows.Scan(&cid, &vname, &uname, &content)
		if err != nil {
			return nil, err
		}

		commentInfo := &def.CommentInfo{Id: cid, VideoName: vname, AuthorName: uname, Content: content}
		res = append(res, commentInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

func DeleteComment(cid string) error {
	stmtIns, err := dbConn.Prepare("DELETE FROM comments WHERE id = ?")
	if err != nil {
		log.Printf("delete comment db prepare error!\n")
		return err
	}

	_, err = stmtIns.Exec(cid)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

//管理员视频列表
func ListVideo(from, n int) ([]*def.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info_copy.id, users.nickname, video_info_copy.name, video_info_copy.create_time, modulars.name, content
	FROM video_info_copy,modulars,users,introduction
	WHERE  video_info_copy.modular = modulars.id AND introduction.vid = video_info_copy.id AND users.username = video_info_copy.author_id AND video_info_copy.status = 0
	ORDER BY video_info_copy.create_time desc LIMIT ?,?`)
	if err != nil {
		log.Printf("list db prepare error!\n")
		return nil, err
	}

	var res []*def.VideoInfo
	rows, err := stmtOut.Query(from, n)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var vid, vname, displayTime, modular, authorName, icon, urlOriginal, url720p, url480p, url360p, itd string
		var likeNum, collectNum, commentNum int
		err := rows.Scan(&vid, &authorName, &vname, &displayTime, &modular, &likeNum, &collectNum, &commentNum, &itd)
		if err != nil {
			return nil, err
		}
		icon = def.STREAMSERVER_ADDR + "icon/" + vid
		urlOriginal = def.STREAMSERVER_ADDR + "video/" + vid + "org"
		url720p = def.STREAMSERVER_ADDR + "video/" + vid + "720p"
		url480p = def.STREAMSERVER_ADDR + "video/" + vid + "480p"
		url360p = def.STREAMSERVER_ADDR + "video/" + vid + "360p"
		videoInfo := &def.VideoInfo{Id: vid, AuthorName: authorName, Name: vname, DisplayCtime: displayTime, Modular: modular, LikeNum: likeNum, CollectNum: collectNum, CommentNum: commentNum, Icon: icon, UrlOriginal: urlOriginal, Url720p: url720p, Url480p: url480p, Url360p: url360p, Introduction: itd}
		res = append(res, videoInfo)
	}

	defer stmtOut.Close()
	return res, nil
}

//审核视频
// func ExamVideo
