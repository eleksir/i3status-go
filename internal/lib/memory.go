package lib

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/mem"
)

// Mem struct with mem stats.
type Mem struct {
	Usedpct uint64
	Shared  uint64
	Swap    uint64
}

// Memory stores string with mem stats for i3bar.
var Memory Mem

// UpdateMemStats parses mem info stats.
func (c MyConfig) UpdateMemStats() {
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
			if Memory.Usedpct != uint64(v.UsedPercent) || Memory.Shared != v.Shared/1024/1024 ||
				Memory.Swap != v.SwapTotal-v.SwapFree {
				Memory.Usedpct = uint64(v.UsedPercent)
				Memory.Shared = v.Shared / 1024 / 1024
				Memory.Swap = sw.Used
				c.UpdateReady <- true
			}
		} else {
			if Memory.Usedpct != uint64(v.UsedPercent) || Memory.Shared != v.Shared/1024/1024 {
				Memory.Usedpct = uint64(v.UsedPercent)
				Memory.Shared = v.Shared / 1024 / 1024
				c.UpdateReady <- true
			}
		}

		time.Sleep(3 * time.Second)
	}
}
