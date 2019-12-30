package utils

import (
	"log"
	"net/http"

	"github.com/pili-video-server/streamserver/def"
)

func SendFormatVideoRequest(vid string) {
	url := "http://" + def.SCHEDULER_ADDR + "video_format/" + vid
	_, err := http.Get(url)
	if err != nil {
		log.Printf("internal error:%v!\n", err)
		return
	}
}
