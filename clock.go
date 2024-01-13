package main

import (
	"fmt"
	"time"
)

// Clock shows current date and time on i3bar.
var Clock = "Thu, 1 Jan 1970   1:00"

// UpdateClock get and updates (on i3bar) info about system clock.
func UpdateClock() {
	for {
		currentTime := time.Now()
		hours, minutes, _ := currentTime.Clock()
		year, month, day := currentTime.Date()
		dow := currentTime.Weekday()
		rmonth := [12]string{"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"}
		rdow := [7]string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

		myclock := fmt.Sprintf("     %s, %d %s %d  % 2d:%02d  ", rdow[dow], day, rmonth[month-1], year, hours, minutes)

		if myclock != Clock {
			Clock = myclock
			UpdateReady <- true
		}

		time.Sleep(1 * time.Second)
	}
}
