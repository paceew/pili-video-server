package dbops

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"github.com/pili-video-server/scheduler/def"
)

func AddDeletion(vid string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO delect_vid(id) VALUES ?")
	if err != nil {
		log.Printf("add deletionID db prepare error: %v\n", err)
		return err
	}

	_, err = stmtIns.Exec(vid)
	if err != nil {
		log.Printf("add vid db error: %v\n", err)
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ReadDeletion(count int) ([]string, error) {
	var ids []string
	stmtOut, err := dbConn.Prepare("SELECT id FROM delect_vid LIMIT ?")
	if err != nil {
		log.Printf("read deletionID db prepare error: %v!\n", err)
		return ids, err
	}

	rows, err := stmtOut.Query(count)
	if err != nil {
		log.Printf("read vid db error: %v\n", err)
		return ids, err
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	defer stmtOut.Close()
	return ids, nil

}

func DeleteDeletion(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM delect_vid WHERE id = ?")
	if err != nil {
		log.Printf("delete vid db prepare error: %v\n", err)
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("delete vid db prepare error: %v\n", err)
		return err
	}

	defer stmtDel.Close()
	return nil

}

//读取视频的vid，likeNum,collNum,CommNum和createTime
func ReadData() ([]*def.VideoData, error) {
	stmtOut, err := dbConn.Prepare("SELECT id,like_number, collect_number, comment_number, create_time FROM video_info ")
	if err != nil {
		log.Printf("read data db error :%v!\n", err)
		return nil, err
	}

	var res []*def.VideoData
	rows, err := stmtOut.Query()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var likeNum, collNum, commNum int
		var time, vid string
		err := rows.Scan(&vid, &likeNum, &collNum, &commNum, time)
		if err != nil {
			return nil, err
		}
		data := &def.VideoData{Vid: vid, LikeNum: likeNum, CollNum: collNum, CommNum: commNum, Creatime: time}
		res = append(res, data)
	}

	defer stmtOut.Close()
	return res, nil
}

//批量写入数据到video
func WriteData(data []*def.VideoData) error {
	//构建批量插入sql
	sqlStr := "INSERT INTO video_info(like_number,hot) VALUES "
	var vals []interface{}

	for _, row := range data {
		sqlStr += "(?,?),"
		vals = append(vals, row.LikeNum, row.Hot)
	}

	sqlStr = sqlStr[0 : len(sqlStr)-2]

	stmtIns, err := dbConn.Prepare(sqlStr)
	if err != nil {
		log.Printf("db error :%v!\n", err)
		return err
	}

	_, err = stmtIns.Exec(vals)
	if err != nil {
		return err
	}

	return nil
}

//根据从redis获取like_Num
func ReadLikeNum(vid string) (int, error) {
	conn := Pool.Get()
	if conn == nil {
		log.Printf("redis error!\n")
	}
	defer conn.Close()

	likestr := "like_" + vid
	Num, _ := redis.Int(conn.Do("scard", likestr))

	return Num, nil
}
