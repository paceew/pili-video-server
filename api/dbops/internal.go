package dbops

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/pili-video-server/api/def"
)

// func InserSession(sid string, ttl int64, username string) error {
// 	ttlstr := strconv.FormatInt(ttl, 10)
// 	stmtIns, err := dbConn.Prepare("INSERT INTO sessions (id, login_name, TTL) VALUES (?, ?, ?)")
// 	if err != nil {
// 		return err
// 	}

// 	_, err = stmtIns.Exec(sid, username, ttlstr)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmtIns.Close()
// 	return nil
// }

// func RetrieveSession(sid string) (*def.SimpleSession, error) {
// 	session := &def.SimpleSession{}
// 	stmtOut, err := dbConn.Prepare("SELECT TTL, login_name FROM sessions WHERE id=?")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var ttl string
// 	var username string
// 	err = stmtOut.QueryRow(sid).Scan(&ttl, &username)
// 	if err != nil && err != sql.ErrNoRows {
// 		return nil, err
// 	}

// 	if ttlInt, err := strconv.ParseInt(ttl, 10, 64); err == nil {
// 		session.TTL = ttlInt
// 		session.Username = username
// 	} else {
// 		return nil, err
// 	}
// 	defer stmtOut.Close()
// 	return session, nil
// }

// func RetrieveAllSession() (*sync.Map, error) {
// 	sessionMap := &sync.Map{}
// 	stmtOut, err := dbConn.Prepare("SELECT id, login_name, TTL FROM sessions")
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := stmtOut.Query()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for rows.Next() {
// 		var id int64
// 		var ttl string
// 		var username string
// 		err = rows.Scan(&id, &username, &ttl)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if ttlInt, err1 := strconv.ParseInt(ttl, 10, 64); err1 == nil {
// 			session := &def.SimpleSession{Username: username, TTL: ttlInt}
// 			sessionMap.Store(id, session)
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	return sessionMap, nil
// }

// func DeleteSession(sid string) error {
// 	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE id= ?")
// 	if err != nil {
// 		return err
// 	}

// 	_, err = stmtOut.Exec(sid)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmtOut.Close()
// 	return nil
// }

func GetUserNameByVid(vid string) (string, error) {
	stmtout, err := dbConn.Prepare("SELECT users.username From users, video_info WHERE video_info.author_id = users.id WHERE video_info.id = ?")
	if err != nil {
		return "", err
	}
	var name string
	err = stmtout.QueryRow(vid).Scan(&name)
	if err != nil {
		return "", err
	}

	defer stmtout.Close()
	return name, err
}

func GetUserNameByCid(cid string) (string, error) {
	stmtout, err := dbConn.Prepare("SELECT users.username From users, comments WHERE comments.author_id = users.id WHERE comments.id = ?")
	if err != nil {
		return "", err
	}
	var name string
	err = stmtout.QueryRow(cid).Scan(&name)
	if err != nil {
		return "", err
	}

	defer stmtout.Close()
	return name, err
}

func GetModIdByName(modName string) (int, error) {
	stmtOut, err := dbConn.Prepare("SELECT id FROM modulars WHERE name = ?")
	if err != nil {
		return 0, err
	}
	var id int
	err = stmtOut.QueryRow(modName).Scan(&id)
	if err != nil {
		return 0, err
	}
	defer stmtOut.Close()
	return id, err
}

func LoadMessageFromDB(uname string) {
	uid, _ := GetUserId(uname)
	stmtout, err := dbConn.Prepare("SELECT count(user_id) FROM private_messages WHERE user_id = ? AND status = 1 GROUP BY user_id ")
	if err != nil {
		return
	}

	var num int
	err = stmtout.QueryRow(uid).Scan(&num)
	if err != nil {
		return
	}

	conn := Pool.Get()
	if conn == nil {
		log.Printf("load message from db error!\n")
	}
	defer conn.Close()

	_, err = conn.Do("SET", string(uid)+"_messnum", string(num))
	if err != nil {
		return
	}

	loadMessage(uname)
}

func loadMessage(uname string) {
	//具体消息写入redis
	uid, _ := GetUserId(uname)
	stmtout1, err := dbConn.Prepare(`SELECT username,is_sender,message_content,send_time,status FROM private_messages,users
	 WHERE friend_id = users.id AND user_id = ? AND status != 3 ORDER BY friend_id asc,send_time desc`)
	if err != nil {
		log.Printf("load message db prepare error:%v!\n", err)
		return
	}
	defer stmtout1.Close()

	conn := Pool.Get()
	if conn == nil {
		log.Printf("load message from db error!\n")
	}
	defer conn.Close()

	rows, _ := stmtout1.Query(uid)
	for rows.Next() {
		var friendName, content, sendTime, sendName string
		var isSender, status int
		err = rows.Scan(&friendName, &isSender, &content, &sendTime, &status)
		if err != nil {
			log.Printf("rows error:%v\n", err)
			return
		}
		//如果是1则表示发送者是本人，0则是对方
		if isSender == 1 {
			sendName = uname
		} else {
			sendName = friendName
		}
		messageDetail := &def.MessageDetail{SendName: sendName, Content: content, SendTime: sendTime, Status: status}
		if res, err := json.Marshal(messageDetail); err != nil {
			log.Printf("json marshal error:%v!\n", err)
			return
		} else {
			_, err = conn.Do("RPUSH", "message_list_"+friendName, res)
			if err != nil {
				log.Printf("someting error:%v!\n", err)
				return
			}
		}
	}

	//从db读取消息列表
	stmtout2, err := dbConn.Prepare(`SELECT username FROM (SELECT * FROM private_messages ORDER BY id DESC)p,users 
	WHERE friend_id = users.id AND user_id = ? AND status != 3 GROUP BY p.user_id,username`)
	if err != nil {
		log.Printf("load message db prepare error:%v!\n", err)
		return
	}
	defer stmtout2.Close()

	rows, _ = stmtout2.Query(uid)
	var res []*def.Message
	for rows.Next() {
		var friendName string
		err := rows.Scan(&friendName)
		if err != nil {
			log.Printf("rows error:%v\n", err)
			return
		}
		//从队列取出每个私信的第一个消息,写入消息列表redis
		messageDetailByte, _ := redis.Bytes(conn.Do("LPOP", "message_list_"+friendName))
		_, _ = conn.Do("LPUSH", "message_list_"+friendName, messageDetailByte)
		messageDetail := &def.MessageDetail{}
		if err := json.Unmarshal(messageDetailByte, messageDetail); err != nil {
			log.Printf("unmarshal error;%v!\n", err)
			return
		}
		//判断第一条消息是否未读(设置未读条数)
		var num int = 0
		if messageDetail.Status == 1 {
			num = 1
		}
		message := &def.Message{FriendName: friendName, Num: num, FristMessage: messageDetail}
		res = append(res, message)
	}

	messageList := &def.MessageList{Messages: res}
	if temp, err := json.Marshal(messageList); err != nil {
		log.Printf("json marshal error:%v!\n", err)
		return
	} else {
		_, err = conn.Do("SET", "messagelist", temp)
		if err != nil {
			log.Printf("someting error:%v!\n", err)
			return
		}
	}

}
