package taskrunner

const (
	READY_TO_DISPATCH = "d"
	READY_TO_EXECUTE  = "e"
	CLOSE             = "c"

	VIDEOS_PATH = "./videos/"
	ICON_PATH   = "./icon/"

	INTERVAL = 360 // 时间间隔360s
)

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error
