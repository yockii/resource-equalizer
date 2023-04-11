package mem

import (
	"github.com/shirou/gopsutil/v3/mem"
	logger "github.com/sirupsen/logrus"
	"resource-equalizer/internal/job"
	"resource-equalizer/pkg/config"
	"time"

	"arena"
)

var status = 0 // 0-无动作 1-启动占用中 -1-停止占用中
// var memStructs []interface{}
var memStructs []*arena.Arena

func StartMemWatcher() {
	// 初始化配置
	_, err := job.Scheduler.Every(3).Seconds().Do(Watcher)
	if err != nil {
		logger.Error("mem监控任务添加失败!", err)
	}
}

func Watcher() {
	if config.GetBool("mem.disable") {
		return
	}

	v, _ := mem.VirtualMemory()

	//fmt.Printf("内存信息: %f, 数据量: %d \n", v.UsedPercent, len(memStructs))

	percent := v.UsedPercent
	if percent > 0 {
		// 判断是否高于高位值
		if percent > config.GetFloat64("mem.high") {
			// 停止1个任务
			stopOneTask()
			return
		} else if percent < config.GetFloat64("mem.low") {
			// 启动一个任务
			startOneTask()
			return
		} else {
			// 介于low和high之间
			if status == 1 {
				// 启动占用中
				if percent < config.GetFloat64("mem.target") {
					// 启动一个任务
					startOneTask()
					return
				}
				status = 0
			} else if status == -1 {
				// 停止占用中
				if percent > config.GetFloat64("mem.target") {
					// 停止1个任务
					stopOneTask()
					return
				}
				status = 0
			}
		}
	}
}

var count = 1
var innerCount = 10_000_000

//var count = 10_000_000

func stopOneTask() {
	//if len(memStructs) >= count {
	//	memStructs = memStructs[count:]
	//} else {
	//	memStructs = make([]interface{}, 0)
	//}

	if len(memStructs) > 0 {
		m := memStructs[0]
		m.Free()
		memStructs = memStructs[1:]
	}

	status = -1
}

func startOneTask() {
	//for i := 0; i < count; i++ {
	//	memStructs = append(memStructs, struct {
	//		t string
	//	}{
	//		time.Now().String(),
	//	})
	//}

	m := arena.NewArena()
	s := arena.MakeSlice[string](m, innerCount, innerCount)
	for i := 0; i < innerCount; i++ {
		s = append(s, time.Now().String())
	}
	memStructs = append(memStructs, m)

	status = 1
}

func Close() {
	for _, m := range memStructs {
		m.Free()
	}
}
