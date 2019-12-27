package taskrunner

import (
	"time"
)

type Worker struct {
	ticker *time.Ticker
	runner *Runner
}

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
	wDel := NewWorker(5, rDel)
	go wDel.startWorker()
}
