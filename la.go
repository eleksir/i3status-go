package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/load"
)

// La переменная, хранящая значение LA за последнюю минуту.
var La = "-1"

// UpdateLaStats вытаскивает показание LA за последнюю минуту.
func UpdateLaStats() {
	for {
		l, _ := load.Avg()

		lav := fmt.Sprintf("%.2f", l.Load1)

		if La != lav {
			La = lav
			UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}
