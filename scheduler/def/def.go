package def

var (
	INTERVAL = "* * */1 * * ?" //每小时执行一次
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
