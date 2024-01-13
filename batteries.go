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
			Batts = fmt.Sprintf(
				"<span font='%s' size='%s'>%s</span><span font='%s' size='%s'> no batt</span>",
				Conf.Battery.SymbolFont,
				Conf.Battery.SymbolFontSize,
				Conf.Battery.Symbol,
				Conf.Battery.Font,
				Conf.Battery.FontSize,
			)
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

				// N.B. there can be case when battery is overcharged and shows >100%. It also can indicate that
				//      calibration data is out of date and battery should be re-calibrated.
				ch = int(math.Round((b.Full - b.Current) * (100 / b.Full)))

				switch {
				case ch > 85:
					if Conf.Battery.ChargeColor.Full == "" {
						chStr = fmt.Sprintf(
							"<span font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.Full,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span font='%s' size='%s' foreground='%s'>% 3d%%</span>`,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.Full,
							ch,
						)
					}

				case ch < 85 && ch > 40:
					if Conf.Battery.ChargeColor.AlmostFull == "" {
						chStr = fmt.Sprintf(
							"<span font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.Full,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span font='%s' size='%s' foreground='%s'>% 3d%%</span>`,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.AlmostFull,
							ch,
						)
					}

				case ch <= 40 && ch >= 10:
					if Conf.Battery.ChargeColor.AlmostEmpty == "" {
						chStr = fmt.Sprintf(
							"<span font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.Full,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span font='%s' size='%s' foreground='%s'>% 3d%%</span>`,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.AlmostEmpty,
							ch,
						)
					}

				case ch < 10:
					if Conf.Battery.ChargeColor.Empty == "" {
						chStr = fmt.Sprintf(
							"<span font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.Full,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span font='%s' size='%s' foreground='%s'>% 3d%%</span>`,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							Conf.Battery.ChargeColor.Empty,
							ch,
						)
					}
				}

				battsInfo += fmt.Sprintf(
					"<span font='%s' size='%s'>%s</span><span font='%s' size='%s'>B%d </span>%s<span font='%s' size='%s'> %s</span>",
					Conf.Battery.SymbolFont,
					Conf.Battery.SymbolFontSize,
					Conf.Battery.Symbol,
					Conf.Battery.Font,
					Conf.Battery.FontSize,
					i,
					chStr,
					Conf.Battery.Font,
					Conf.Battery.FontSize,
					status,
				)
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
