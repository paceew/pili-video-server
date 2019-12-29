package timertask

import (
	"github.com/pili-video-server/scheduler/def"
	"github.com/robfig/cron"
)

func Start() {
	c := cron.New()
	c.AddFunc(def.INTERVAL, TaskCountHot)

	c.Start()
}
