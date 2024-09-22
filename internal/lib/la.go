package lib

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/load"
)

// UpdateLaStats вытаскивает показание LA за последнюю минуту.
func (c MyConfig) UpdateLaStats() {
	for {
		l, _ := load.Avg()

		lav := fmt.Sprintf("%.2f", l.Load1)

		if c.La != lav {
			c.La = lav
			c.UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}
