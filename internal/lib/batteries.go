package lib

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/distatus/battery"
)

// UpdateBatteryInfo updates info about battery charge.
func (c *MyConfig) UpdateBatteryInfo() {
	var (
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

		var (
			batteries []*battery.Battery = []*battery.Battery{}
			Batts     string
			ch        int
			status    string
		)

		if c.Battery.UseSysfs {
			for _, file := range c.Battery.SysfsFiles {
				var batt battery.Battery

				_, err := os.Stat(file)

				// If file exists.
				if !os.IsNotExist(err) {
					var (
						myCharge string
						myStatus string
					)

					b, err := os.ReadFile(file)

					if err != nil {
						log.Printf("Unable to read file %s: %s", file, err)
						continue
					}

					myCharge = strings.Trim(string(b), "\n")

					statusFile := filepath.Dir(file) + "/status"

					// If file exists.
					if !os.IsNotExist(err) {
						b, err := os.ReadFile(statusFile)

						if err != nil {
							log.Printf("Unable to read file %s: %s", statusFile, err)
							continue
						}

						myStatus = strings.Trim(string(b), "\n")
					} else {
						continue
					}

					// Make stub battery class :) .
					switch myStatus {
					case "Charging":
						batt.State.Raw = battery.Charging
					case "Discharging":
						batt.State.Raw = battery.Discharging
					case "Empty":
						batt.State.Raw = battery.Empty
					case "Full":
						batt.State.Raw = battery.Full
					default:
						batt.State.Raw = battery.Unknown
					}

					ch, err := strconv.Atoi(myCharge)

					if err != nil {
						log.Printf("Unable to convert string from file %s to integer", file)
						continue
					}

					batt.Current = float64(ch)

					batteries = append(batteries, &batt)
				}
			}
		} else {
			batteries, _ = battery.GetAll()
			// In theory, this module should give for each entry its separate err, but in practice it gives
			// one single err for all entries, so we cannot detemine whist exatly entry errored.
		}

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
			if c.Battery.UseSysfs {
				ch = int(b.Current)
			} else {
				ch = int(math.Round((b.Full - (b.Full - b.Current)) * (100 / b.Full)))
			}

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
