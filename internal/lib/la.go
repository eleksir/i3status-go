package lib

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/load"
)

// UpdateLaStats вытаскивает показание LA за последнюю минуту.
func (c *MyConfig) UpdateLaStats() {
	for {
		l, _ := load.Avg()

		lav := fmt.Sprintf("%.2f", l.Load1)

		if c.Values.La != lav {
			c.Values.La = lav
			c.Channels.UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}
