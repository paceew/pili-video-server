package timertask

import (
	"log"
	"time"

	"github.com/pili-video-server/scheduler/dbops"
	"github.com/pili-video-server/scheduler/def"
)

type Task struct {
}

//收藏*0.3 + 点赞*0.4 + 评论*0.5 新视频（一天内发布）[提权,总值*1.5]
func calculaHot(data *def.VideoData) float32 {
	hot := float32(data.CollNum)*0.3 + float32(data.LikeNum)*0.4 + float32(data.CommNum)*0.5

	//time string 转时间戳
	loc, _ := time.LoadLocation("Local")
	the_time, err := time.ParseInLocation("2006-01-02 15:04:05", data.Creatime, loc)
	if err != nil {
		log.Printf("time error:%v!\n", err)
		return hot
	}
	unix_time := the_time.Unix()

	if unix_time >= time.Now().AddDate(0, 0, -1).Unix() && unix_time <= time.Now().Unix() {
		hot *= 1.5
	}

	return hot
}

//先读出所有video data 然后根据每个vid去redis读取likeNum，然后计算每个video hot 再写入video data
func (t *Task) TaskCountHot() {
	log.Printf("begin to count hot !\n")
	data, err := dbops.ReadData()
	if err != nil {
		log.Printf("db error:%v\n", err)
		return
	}

	for _, row := range data {
		row.LikeNum, _ = dbops.ReadLikeNum(row.Vid)
		row.Hot = calculaHot(row)
	}

	err = dbops.WriteData(data)
	if err != nil {
		log.Printf("db error :%v\n", err)
		return
	}
}
