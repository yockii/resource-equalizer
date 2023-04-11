package main

import (
	"github.com/panjf2000/ants/v2"
	"resource-equalizer/internal/cpu"
	"resource-equalizer/internal/job"
	"resource-equalizer/internal/mem"
	"resource-equalizer/pkg/config"
)

func init() {
	config.Set("logger.level", "info")
	config.InitialLogger()
}

func main() {
	defer ants.Release()

	cpu.StartCpuWatcher()
	mem.StartMemWatcher()
	defer mem.Close()

	job.Start()
}
