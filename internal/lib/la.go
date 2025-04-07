package lib

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/load"
)

// UpdateLaStats вытаскивает показание LA за последнюю минуту.
func (c *MyConfig) UpdateLaStats() {
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

		l, _ := load.Avg()

		lav := fmt.Sprintf("%.2f", l.Load1)

		if c.Values.La != lav {
			c.Values.La = lav
			c.Channels.UpdateReady <- true
		}
	}
}
