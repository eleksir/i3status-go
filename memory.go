package main

import (
	"log"
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
		v, err := mem.VirtualMemory()

		if err != nil {
			log.Printf("Unable to get memory statistics: %s", err)
			time.Sleep(1 * time.Second)

			continue
		}

		sw, err := mem.SwapMemory()

		if err != nil {
			log.Printf("Unable to get swap statistics: %s", err)
			time.Sleep(1 * time.Second)

			continue
		}

		if Conf.Mem.ShowSwap {
			if Memory.Usedpct != int64(v.UsedPercent) || Memory.Shared != int64(v.Shared/1024/1024) ||
				Memory.Swap != int64(v.SwapTotal-v.SwapFree) {
				Memory.Usedpct = int64(v.UsedPercent)
				Memory.Shared = int64(v.Shared / 1024 / 1024)
				Memory.Swap = int64(sw.Total - sw.Free)
				UpdateReady <- true
			}
		} else {
			if Memory.Usedpct != int64(v.UsedPercent) || Memory.Shared != int64(v.Shared/1024/1024) {
				Memory.Usedpct = int64(v.UsedPercent)
				Memory.Shared = int64(v.Shared / 1024 / 1024)
				UpdateReady <- true
			}
		}

		time.Sleep(3 * time.Second)
	}
}
