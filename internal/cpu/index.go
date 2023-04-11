package cpu

import (
	"github.com/shirou/gopsutil/v3/cpu"
	logger "github.com/sirupsen/logrus"
	"resource-equalizer/internal/job"
	"resource-equalizer/pkg/config"
	"time"
)

var tasks []int64
var taskChan = make(map[int64]chan bool)
var status = 0 // 0-无动作 1-启动占用中 -1-停止占用中

func StartCpuWatcher() {
	// 初始化配置
	_, err := job.Scheduler.Every(3).Seconds().Do(Watcher)
	if err != nil {
		logger.Error("CPU监控任务添加失败!", err)
	}
}

func Watcher() {
	if config.GetBool("cpu.disable") {
		return
	}
	info, _ := cpu.Percent(5*time.Second, false)
	percent := info[0]

	if percent > 0 {
		// 判断是否高于高位值
		if percent > config.GetFloat64("cpu.high") {
			// 停止1个任务
			stopOneTask()
			return
		} else if percent < config.GetFloat64("cpu.low") {
			// 启动一个任务
			startOneTask()
			return
		} else {
			// 介于low和high之间
			if status == 1 {
				// 启动占用中
				if percent < config.GetFloat64("cpu.target") {
					// 启动一个任务
					startOneTask()
					return
				}
				status = 0
			} else if status == -1 {
				// 停止占用中
				if percent > config.GetFloat64("cpu.target") {
					// 停止1个任务
					stopOneTask()
					return
				}
				status = 0
			}
		}
	}
}

func stopOneTask() {
	if len(tasks) > 0 {
		taskChan[tasks[0]] <- true
		tasks = tasks[1:]
		status = -1
	}
}

func startOneTask() {
	c := make(chan bool)
	id := time.Now().UnixNano()
	tasks = append(tasks, id)
	taskChan[id] = c
	go doTask(c)
	status = 1
}

var count = 10_000

func doTask(c chan bool) {
	for {
		select {
		case <-c:
			return
		default:
			// 占用cpu进行大量计算
			var fib = [3]int{0, 1, 0}
			func() {
				for i := 2; i <= count; i++ {
					fib[2] = fib[0] + fib[1]
					fib[0] = fib[1]
					fib[1] = fib[2]
				}
			}()
		}
	}
}
