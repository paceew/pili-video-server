package session

import (
	"encoding/json"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/pace/sample/api/dbops"
	"github.com/pace/sample/api/def"
	"github.com/pace/sample/api/utils"
)

// var Pool redis.Pool

// func init() {
// 	Pool = redis.Pool{
// 		MaxIdle:     16,
// 		MaxActive:   32,
// 		IdleTimeout: 120,
// 		Dial: func() (redis.Conn, error) {
// 			return redis.Dial("tcp", "127.0.0.1:6379")
// 		},
// 	}
// }

// func nowInMilli() int64 {
// 	return time.Now().UnixNano() / 1000000
// }

// func DeleteExpiredSession(sid string) {
// 	sessionMap.Delete(sid)
// 	dbops.DeleteSession(sid)
// }

// func LoadSessionsFromDB() {
// 	r, err := dbops.RetrieveAllSession()
// 	if err != nil {
// 		return
// 	}

// 	r.Range(func(k, v interface{}) bool {
// 		session := v.(*def.SimpleSession)
// 		sessionMap.Store(k, session)
// 		return true
// 	})
// }

func GenerateNewSessionId(username string) string {
	id, _ := utils.NewUUID()
	session := &def.SimpleSession{Id: id, Username: username}

	conn := dbops.Pool.Get()
	if conn == nil {
		log.Printf("generate new session error!\n")
	}
	value, err := json.Marshal(session)
	if err != nil {
		log.Printf("marshal session error:%v!\n", err)
	}

	_, err = conn.Do("SET", session.Id, value, "EX", 60)
	if err != nil {
		log.Printf("someting error:%v!\n", err)
	}

	defer conn.Close()
	return id
}

func DeleteExpiredSession(sid string) {
	conn := dbops.Pool.Get()
	if conn == nil {
		log.Printf("generate new session error!\n")
	}
	_, err := conn.Do("DEL", sid)
	if err != nil {
		log.Printf("someting error:%v!\n", err)
	}

	defer conn.Close()
}

func IsSessionExpired(sid string) (string, bool) {
	conn := dbops.Pool.Get()
	if conn == nil {
		log.Printf("generate new session error!\n")
	}
	defer conn.Close()
	ok, _ := redis.Bool(conn.Do("EXISTS", sid))
	if ok {
		res, err := redis.Bytes(conn.Do("GET", sid))
		if err != nil {
			log.Printf("someting error:%v!\n", err)
			return "", true
		}

		session := &def.SimpleSession{}
		if err = json.Unmarshal(res, session); err != nil {
			log.Printf("unmarshal error:%v!\n", err)
			return "", true
		}
		return session.Username, false
	}
	return "", true
}

func ReSetSession(sid string, username string) {
	session := &def.SimpleSession{Id: sid, Username: username}

	conn := dbops.Pool.Get()
	if conn == nil {
		log.Printf("generate new session error!\n")
	}
	defer conn.Close()

	value, err := json.Marshal(session)
	if err != nil {
		log.Printf("marshal session error:%v!\n", err)
	}

	_, err = conn.Do("SET", session.Id, value, "EX", 60)
	if err != nil {
		log.Printf("someting error:%v!\n", err)
	}

}
