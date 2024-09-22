package lib

import (
	"fmt"
	"math"
	"time"

	"github.com/distatus/battery"
)

// UpdateBatteryInfo updates info about battery charge.
func (c MyConfig) UpdateBatteryInfo() {
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
							"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.Color,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
							Conf.Battery.ChargeColor.Full,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					}

				case ch < 85 && ch > 40:
					if Conf.Battery.ChargeColor.AlmostFull == "" {
						chStr = fmt.Sprintf(
							"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.Color,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
							Conf.Battery.ChargeColor.AlmostFull,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					}

				case ch <= 40 && ch >= 10:
					if Conf.Battery.ChargeColor.AlmostEmpty == "" {
						chStr = fmt.Sprintf(
							"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.Color,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
							Conf.Battery.ChargeColor.AlmostEmpty,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					}

				case ch < 10:
					if Conf.Battery.ChargeColor.Empty == "" {
						chStr = fmt.Sprintf(
							"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
							Conf.Battery.Color,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					} else {
						chStr = fmt.Sprintf(
							`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
							Conf.Battery.ChargeColor.Empty,
							Conf.Battery.Background,
							Conf.Battery.Font,
							Conf.Battery.FontSize,
							ch,
						)
					}
				}

				battsInfo += fmt.Sprintf(
					"<span color='%s' background='%s' font='%s' size='%s'>%s</span><span color='%s' background='%s' font='%s' size='%s'>B%d </span>%s<span color='%s' background='%s' font='%s' size='%s'> %s</span>",
					Conf.Battery.Color,
					Conf.Battery.Background,
					Conf.Battery.SymbolFont,
					Conf.Battery.SymbolFontSize,
					Conf.Battery.Symbol,
					Conf.Battery.Color,
					Conf.Battery.Background,
					Conf.Battery.Font,
					Conf.Battery.FontSize,
					i,
					chStr,
					Conf.Battery.Color,
					Conf.Battery.Background,
					Conf.Battery.Font,
					Conf.Battery.FontSize,
					status,
				)
			}

			if battsInfo != "" {
				Batts = battsInfo
			}
		}

		if c.BatteryString == Batts {
			c.BatteryString = Batts
			c.UpdateReady <- true
		}

		time.Sleep(5 * time.Second)
	}
}
