package def

const (
	INTERVAL = 3600 //每3600s执行一次
)

//避免循环引用，用于timertask
type VideoData struct {
	Vid      string
	LikeNum  int
	CollNum  int
	CommNum  int
	Creatime string
	Hot      float32
}
