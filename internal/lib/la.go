package lib

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/load"
)

// UpdateLaStats вытаскивает показание LA за последнюю минуту.
func (c *MyConfig) UpdateLaStats() {
	ticker := time.NewTicker(time.Second * 3)

	for range ticker.C {
		l, _ := load.Avg()

		lav := fmt.Sprintf("%.2f", l.Load1)

		if c.Values.La != lav {
			c.Values.La = lav
			c.Channels.UpdateReady <- true
		}
	}
}
