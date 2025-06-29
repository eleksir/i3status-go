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

// UpdateMemStats parses mem info stats.
func (c *MyConfig) UpdateMemStats() {
	var (
		InitialDelay       = 100 * time.Millisecond
		LoopIterationDelay = 3 * time.Second
		Delay              = InitialDelay
		ticker             = time.NewTicker(Delay)
	)

	for range ticker.C {
		if Delay == InitialDelay {
			Delay = LoopIterationDelay
			ticker.Reset(Delay)
		}

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

		if c.Mem.ShowSwap {
			if c.Values.Memory.Usedpct != uint64(v.UsedPercent) || c.Values.Memory.Shared != v.Shared/1024/1024 || c.Values.Memory.Swap != v.SwapTotal-v.SwapFree {
				c.Values.Memory.Usedpct = uint64(v.UsedPercent)
				c.Values.Memory.Shared = v.Shared / 1024 / 1024
				c.Values.Memory.Swap = sw.Used / 1024 / 1024
				c.Channels.UpdateReady <- true
			}
		} else {
			if c.Values.Memory.Usedpct != uint64(v.UsedPercent) || c.Values.Memory.Shared != v.Shared/1024/1024 {
				c.Values.Memory.Usedpct = uint64(v.UsedPercent)
				c.Values.Memory.Shared = v.Shared / 1024 / 1024
				c.Channels.UpdateReady <- true
			}
		}
	}
}
