package job

import (
	"github.com/go-co-op/gocron"
	"time"
)

var Scheduler = gocron.NewScheduler(time.UTC)

// Start 初始化任务系统
func Start() {
	Scheduler.StartBlocking()
}

func StartAsync() {
	Scheduler.StartAsync()
}

// Stop 停止任务系统
func Stop() {
	Scheduler.Stop()
}
