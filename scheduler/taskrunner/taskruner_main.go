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
			go w.runner.StartAll()
		}
	}
}

func Start() {
	rDel := NewRunner(5, true, VideoClearDispatcher, VideoClearExecutor)
	wDel := NewWorker(INTERVAL, rDel)
	go wDel.startWorker()
}
