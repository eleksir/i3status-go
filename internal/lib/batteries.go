package lib

import (
	"fmt"
	"math"
	"time"

	"github.com/distatus/battery"
)

// UpdateBatteryInfo updates info about battery charge.
func (c *MyConfig) UpdateBatteryInfo() {
	var (
		ch                 int
		status             string
		Batts              string
		InitialDelay       = 100 * time.Millisecond
		LoopIterationDelay = 5 * time.Second
		Delay              = InitialDelay
		ticker             = time.NewTicker(Delay)
	)

	for range ticker.C {
		if Delay == InitialDelay {
			Delay = LoopIterationDelay
			ticker.Reset(Delay)
		}

		batteries, _ := battery.GetAll()
		// In theory, this module should give for each entry its separate err, but in practice it gives
		// one single err for all entries, so we cannot detemine whist exatly entry errored.

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
			ch = int(math.Round((b.Full - (b.Full - b.Current)) * (100 / b.Full)))

			switch {
			case ch <= 500 && ch >= 84:
				if c.Battery.ChargeColor.Full == "" {
					chStr = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
						c.Battery.Color,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				} else {
					chStr = fmt.Sprintf(
						`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
						c.Battery.ChargeColor.Full,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				}

			case ch < 85 && ch > 40:
				if c.Battery.ChargeColor.AlmostFull == "" {
					chStr = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
						c.Battery.Color,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				} else {
					chStr = fmt.Sprintf(
						`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
						c.Battery.ChargeColor.AlmostFull,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				}

			case ch <= 40 && ch >= 10:
				if c.Battery.ChargeColor.AlmostEmpty == "" {
					chStr = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
						c.Battery.Color,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				} else {
					chStr = fmt.Sprintf(
						`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
						c.Battery.ChargeColor.AlmostEmpty,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				}

			case ch < 10 && ch >= 0:
				if c.Battery.ChargeColor.Empty == "" {
					chStr = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>",
						c.Battery.Color,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				} else {
					chStr = fmt.Sprintf(
						`<span color='%s' background='%s' font='%s' size='%s'>% 3d%%</span>`,
						c.Battery.ChargeColor.Empty,
						c.Battery.Background,
						c.Battery.Font,
						c.Battery.FontSize,
						ch,
					)
				}

			default:
				continue
			}

			battsInfo += fmt.Sprintf(
				"<span color='%s' background='%s' font='%s' size='%s'>%s</span><span color='%s' background='%s' font='%s' size='%s'>B%d </span>%s<span color='%s' background='%s' font='%s' size='%s'> %s</span>",
				c.Battery.Color,
				c.Battery.Background,
				c.Battery.SymbolFont,
				c.Battery.SymbolFontSize,
				c.Battery.Symbol,
				c.Battery.Color,
				c.Battery.Background,
				c.Battery.Font,
				c.Battery.FontSize,
				i,
				chStr,
				c.Battery.Color,
				c.Battery.Background,
				c.Battery.Font,
				c.Battery.FontSize,
				status,
			)
		}

		if battsInfo != "" {
			Batts = battsInfo
		}

		if c.Values.BatteryString != Batts {
			c.Values.BatteryString = Batts
			c.Channels.UpdateReady <- true
		}
	}
}
