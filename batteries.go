package main

import (
	"fmt"
	"math"
	"time"

	"github.com/distatus/battery"
)

// Batt shows battery charge on i3bar.
var Batt = "<big>⚡</big> ??% •"

// UpdateBatteryInfo updates info about battery charge.
func UpdateBatteryInfo() {
	var (
		ch     int
		status string
		Batts  string
	)

	for {
		batteries, err := battery.GetAll()

		if err != nil {
			Batts = "<big>⚡</big> no batt"
		} else {
			var (
				battsInfo string
				chStr     string
			)

			for i, b := range batteries {
				switch b.State.Raw {
				case battery.Charging:
					status = `▲`
				case battery.Discharging:
					status = `▼`
				case battery.Empty:
					status = `✘`
				default:
					status = `•`
				}

				// N.B. there can be case when battery is overcharged and shows >100%. It alse can indicate that
				//      calibration data is out of date and battery should be re-calibrated.
				ch = int(math.Round((b.Full - b.Current) * (100 / b.Full)))

				switch {
				case ch > 85:
					if Conf.Battery.Color.Full == "" {
						chStr = fmt.Sprintf("% 3d%%", ch)
					} else {
						chStr = fmt.Sprintf(
							`<span foreground="%s">% 3d%%</span>`,
							Conf.Battery.Color.Full,
							ch,
						)
					}

				case ch < 85 && ch > 40:
					if Conf.Battery.Color.AlmostFull == "" {
						chStr = fmt.Sprintf("% 3d%%", ch)
					} else {
						chStr = fmt.Sprintf(
							`<span foreground="%s">% 3d%%</span>`,
							Conf.Battery.Color.AlmostFull,
							ch,
						)
					}

				case ch <= 40 && ch >= 10:
					if Conf.Battery.Color.AlmostEmpty == "" {
						chStr = fmt.Sprintf("% 3d%%", ch)
					} else {
						chStr = fmt.Sprintf(
							`<span foreground="%s">% 3d%%</span>`,
							Conf.Battery.Color.AlmostEmpty,
							ch,
						)
					}

				case ch < 10:
					if Conf.Battery.Color.Empty == "" {
						chStr = fmt.Sprintf("% 3d%%", ch)
					} else {
						chStr = fmt.Sprintf(
							`<span foreground="%s">% 3d%%</span>`,
							Conf.Battery.Color.Empty,
							ch,
						)
					}
				}

				battsInfo += fmt.Sprintf("<big>⚡</big>B%d %s %s", i, chStr, status)
			}

			if battsInfo != "" {
				Batts = battsInfo
			}
		}

		if Batt == Batts {
			Batt = Batts
			UpdateReady <- true
		}

		time.Sleep(5 * time.Second)
	}
}
