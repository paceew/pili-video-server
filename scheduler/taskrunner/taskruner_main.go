package taskrunner

import (
	"time"
)

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

//每隔interval 时间执行一次r
func NewWorker(interval time.Duration, r *Runner) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		runner: r,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C:
			w.runner.StartAll()
		}
	}
}

func Start() {
	//tr 删除视频
	rDel := NewRunner(5, true, VideoClearDispatcher, VideoClearExecutor)
	wDel := NewWorker(INTERVAL_DEL, rDel)
	go wDel.startWorker()

	// tr 转换码率
	rFm := NewRunner(1, true, VideoFormatDispatcher, VideoFormatExecutor)
	wFm := NewWorker(INTERVAL_FM, rFm)
	go wFm.startWorker()
}
