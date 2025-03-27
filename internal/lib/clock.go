package lib

import (
	"fmt"
	"time"
)

// UpdateClock get and updates (on i3bar) info about system clock.
func (c *MyConfig) UpdateClock() {
	ticker := time.NewTicker(time.Second * 1)

	for range ticker.C {
		currentTime := time.Now()
		hours, minutes, _ := currentTime.Clock()
		year, month, day := currentTime.Date()
		dow := currentTime.Weekday()
		rmonth := [12]string{"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"}
		rdow := [7]string{"Вс", "Пн", "Вт", "Ср", "Чт", "Пт", "Сб"}

		myclock := fmt.Sprintf("     %s, %d %s %d  % 2d:%02d  ", rdow[dow], day, rmonth[month-1], year, hours, minutes)

		if myclock != c.Values.ClockTime {
			c.Values.ClockTime = myclock
			c.Channels.UpdateReady <- true
		}
	}
}
