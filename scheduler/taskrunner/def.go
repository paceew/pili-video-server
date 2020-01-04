package taskrunner

const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE  = "e"
	CLOSE             = "c"

	VIDEOS_PATH   = "./videos/original/"
	VIDEOS_PATH2  = "./videos/720p/"
	VIDEOS_PATH3  = "./videos/480p/"
	VIDEOS_PATH4  = "./videos/360p/"
	VIDEOS_FORMAT = ".mp4"
	ICON_PATH     = "./icon/"

	INTERVAL_DEL = 3600 // 删除视频时间间隔3600s

	INTERVAL_FM = 1800 // 视频转码时间间隔1800s
)

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error
