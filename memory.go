package main

import (
	"time"

	"github.com/shirou/gopsutil/mem"
)

// Mem struct with mem stats.
type Mem struct {
	Usedpct int64
	Shared  int64
	Swap    int64
}

// Memory stores string with mem stats for i3bar.
var Memory Mem

// UpdateMemStats parses mem info stats.
func UpdateMemStats() {
	for {
		v, _ := mem.VirtualMemory()

		if Memory.Usedpct != int64(v.UsedPercent) || Memory.Shared != int64(v.Shared/1024/1024) ||
			Memory.Swap != int64(v.SwapTotal-v.SwapFree) {
			Memory.Usedpct = int64(v.UsedPercent)
			Memory.Shared = int64(v.Shared / 1024 / 1024)
			Memory.Swap = int64(v.SwapTotal - v.SwapFree)
			UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}
