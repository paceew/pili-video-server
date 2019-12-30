package timertask

import (
	"time"

	"github.com/pili-video-server/scheduler/def"
)

type Worker struct {
	ticker *time.Ticker
	task   *Task
}

func NewWorker(interval time.Duration, t *Task) *Worker {
	return &Worker{
		ticker: time.NewTicker(interval * time.Second),
		task:   t,
	}
}

func (w *Worker) startWorker() {
	for {
		select {
		case <-w.ticker.C:
			w.task.TaskCountHot()
		}
	}
}

func Start() {
	t := &Task{}
	w := NewWorker(def.INTERVAL, t)
	go w.startWorker()
}
