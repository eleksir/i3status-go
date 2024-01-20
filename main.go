package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"
)

// I3BarOutBlock is structure element for I3BarOut, it represents i3bar output json block format.
type I3BarOutBlock struct {
	FullText string `json:"full_text"`
	// ShortText will be shown if not enough room for FullText, threshold width defined in MinWidth
	ShortText    string `json:"short_text,omitempty"`
	Color        string `json:"color,omitempty"`
	Background   string `json:"background,omitempty"`
	Border       string `json:"border,omitempty"`
	BorderTop    int    `json:"border_top"`
	BorderRight  int    `json:"border_right"`
	BorderBottom int    `json:"border_bottom"`
	BorderLeft   int    `json:"border_left"`
	// measured either in pixels or in characters, so either int or string, let's make it string :)
	MinWidth            string `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Separator           bool   `json:"separator"`
	SeparatorBlockWidth int    `json:"separator_block_width"`
	Markup              string `json:"markup,omitempty"`
}

// UpdateReady channel that "refreshes" i3bar (generates stdout json line).
var UpdateReady = make(chan bool)
var MsgChan = make(chan []I3BarOutBlock, 64)

//go:embed i3status-go-example.json
var EmbeddedDefaultConfig embed.FS
var DefaultConfig []byte

func init() {
	DefaultConfig, _ = EmbeddedDefaultConfig.ReadFile("i3status-go-example.json")
}

// Program entry point.
func main() {
	var err error

	Conf, err = readConf()

	// TODO: Проставить дефолтные значения для глобальных переменных модулей.
	Batt = fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
		Conf.Battery.Color,
		Conf.Battery.Background,
		Conf.Battery.SymbolFont,
		Conf.Battery.SymbolFontSize,
		Conf.Battery.Symbol,
	)

	Batt += fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'> ??%% •</span>",
		Conf.Battery.Color,
		Conf.Battery.Background,
		Conf.Battery.Font,
		Conf.Battery.FontSize,
	)

	SoundVolume = fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>🔊</span>",
		Conf.SimpleVolumePa.Color,
		Conf.SimpleVolumePa.Background,
		Conf.SimpleVolumePa.SymbolFont,
		Conf.SimpleVolumePa.SymbolFontSize,
	)

	SoundVolume += fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>:0%%</span>",
		Conf.SimpleVolumePa.Color,
		Conf.SimpleVolumePa.Background,
		Conf.SimpleVolumePa.Font,
		Conf.SimpleVolumePa.FontSize,
	)

	Clock = fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>Thu, 1 Jan 1970   1:00</span>",
		Conf.Clock.Color,
		Conf.Clock.Background,
		Conf.Clock.Font,
		Conf.Clock.FontSize,
	)

	if err != nil {
		log.Panicln(err)
	}

	go Spawner()
	go ParseStdin()
	go CleanZombies()
	go SVPAHandler()
	go PrintToI3bar()

	if Conf.AppButtons.Enabled {
		go UpdateI3WinList()
	}

	// Kick signal handler
	go sigHandler()
	signal.Notify(sigChan,
		syscall.SIGUSR1,
		syscall.SIGUSR2)

	// Populate memory stats
	if Conf.Mem.Enabled {
		go UpdateMemStats()
	}

	// Populate LA stats
	if Conf.LA.Enabled {
		go UpdateLaStats()
	}

	// Populate CPUTemperature
	if Conf.CPUTemp.Enabled {
		go UpdateCPUTemperature()
	}

	// Populate Clock
	if Conf.Clock.Enabled {
		go UpdateClock()
	}

	if Conf.Battery.Enabled {
		go UpdateBatteryInfo()
	}

	if Conf.SimpleVolumePa.Enabled {
		go UpdateVolumeInfo()
	}

	if Conf.NetIf.Enabled {
		go UpdateIfStatus()
	}

	if Conf.Cron.Enabled {
		go RunCron()
	}

	if Conf.Vpn.Enabled {
		go UpdateVPNStatus()
	}

	/*
		I3bar documentation pretends that message protocol must be valid json. In practice, we only have to print valid
		header, empty json array and (potentially infinite) json lines (line that is valid json by itself) that is
		actually json arrays. We do not need to *close* this json at all.
		Gracefully closed json required when i3bar initiates our program to stop|quit, this (should) happens just before
		i3bar itself terminating. So we don't care.
	*/

	// Print header and one empty message and wait for updates
	fmt.Printf(
		"{\"version\": 1, \"stop_signal\": %d, \"cont_signal\": %d, \"click_events\": true}\n",
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	fmt.Println("[ [],")

	for {
		if <-UpdateReady {
			// Actually build json struct, marshal it and print result
			var j []I3BarOutBlock

			// TODO: впилить для первого и последнего батона separator
			if Conf.AppButtons.Enabled {
				for num, app := range Conf.Apps {
					var b I3BarOutBlock

					if num == 0 {
						b.FullText = fmt.Sprintf(
							"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
							Conf.AppButtons.Separator.Left.Color,
							Conf.AppButtons.Separator.Left.Background,
							Conf.AppButtons.Separator.Left.Font,
							Conf.AppButtons.Separator.Left.FontSize,
							Conf.AppButtons.Separator.Left.Symbol,
						)
					}

					b.FullText = app.FullText

					if num == len(Conf.Apps) {
						b.FullText = fmt.Sprintf(
							"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
							Conf.AppButtons.Separator.Right.Color,
							Conf.AppButtons.Separator.Right.Background,
							Conf.AppButtons.Separator.Right.Font,
							Conf.AppButtons.Separator.Right.FontSize,
							Conf.AppButtons.Separator.Right.Symbol,
						)
					}

					b.Background = app.Background
					b.Color = app.Color
					b.Instance = app.Instance
					b.Markup = `pango`
					b.Separator = app.Separator
					b.SeparatorBlockWidth = app.SeparatorBlockWidth
					b.Name = app.Name

					if HasWindows(app.Class, app.Instance) {
						b.Border = app.BorderActive
					} else {
						b.Border = app.Border
					}

					b.BorderTop = 1
					b.BorderRight = 1
					b.BorderBottom = 1
					b.BorderLeft = 1

					j = append(j, b)
				}
			}

			if Conf.CPUTemp.Enabled {
				var b I3BarOutBlock

				b.Color = Conf.CPUTemp.Color
				b.Background = Conf.CPUTemp.Background

				if Conf.CPUTemp.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.CPUTemp.Separator.Left.Color,
						Conf.CPUTemp.Separator.Left.Background,
						Conf.CPUTemp.Separator.Left.Font,
						Conf.CPUTemp.Separator.Left.FontSize,
						Conf.CPUTemp.Separator.Left.Symbol,
					)
				}

				b.FullText += fmt.Sprintf(
					"<span color='%s' background='%s' font='%s' size='%s'>CPU: %d°</span>",
					Conf.CPUTemp.Color,
					Conf.CPUTemp.Background,
					Conf.CPUTemp.Font,
					Conf.CPUTemp.FontSize,
					CPUTemperature,
				)

				if Conf.CPUTemp.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.CPUTemp.Separator.Right.Color,
						Conf.CPUTemp.Separator.Right.Background,
						Conf.CPUTemp.Separator.Right.Font,
						Conf.CPUTemp.Separator.Right.FontSize,
						Conf.CPUTemp.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.Mem.Enabled {
				var b I3BarOutBlock

				b.Color = Conf.Mem.Color
				b.Background = Conf.Mem.Background

				if Conf.Mem.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Mem.Separator.Left.Color,
						Conf.Mem.Separator.Left.Background,
						Conf.Mem.Separator.Left.Font,
						Conf.Mem.Separator.Left.FontSize,
						Conf.Mem.Separator.Left.Symbol,
					)
				}

				if Conf.Mem.ShowSwap {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>M:%d%% SHM:%dM SW:%dM</span>",
						Conf.Mem.Color,
						Conf.Mem.Background,
						Conf.Mem.Font,
						Conf.Mem.FontSize,
						Memory.Usedpct,
						Memory.Shared,
						Memory.Swap,
					)
				} else {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>M:%d%% SHM:%dM</span>",
						Conf.Mem.Color,
						Conf.Mem.Background,
						Conf.Mem.Font,
						Conf.Mem.FontSize,
						Memory.Usedpct,
						Memory.Shared,
					)
				}

				if Conf.Mem.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Mem.Separator.Right.Color,
						Conf.Mem.Separator.Right.Background,
						Conf.Mem.Separator.Right.Font,
						Conf.Mem.Separator.Right.FontSize,
						Conf.Mem.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.LA.Enabled {
				var b I3BarOutBlock

				b.Color = Conf.LA.Color
				b.Background = Conf.LA.Background

				if Conf.LA.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.LA.Separator.Left.Color,
						Conf.LA.Separator.Left.Background,
						Conf.LA.Separator.Left.Font,
						Conf.LA.Separator.Left.FontSize,
						Conf.LA.Separator.Left.Symbol,
					)
				}

				b.FullText += fmt.Sprintf(
					"<span color='%s' background='%s' font='%s' size='%s'>LA:%s</span>",
					Conf.LA.Color,
					Conf.LA.Background,
					Conf.LA.Font,
					Conf.LA.FontSize,
					La,
				)

				if Conf.LA.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.LA.Separator.Right.Color,
						Conf.LA.Separator.Right.Background,
						Conf.LA.Separator.Right.Font,
						Conf.LA.Separator.Right.FontSize,
						Conf.LA.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.NetIf.Enabled {
				var b I3BarOutBlock

				b.Color = Conf.NetIf.Color
				b.Background = Conf.NetIf.Background

				if Conf.NetIf.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.NetIf.Separator.Left.Color,
						Conf.NetIf.Separator.Left.Background,
						Conf.NetIf.Separator.Left.Font,
						Conf.NetIf.Separator.Left.FontSize,
						Conf.NetIf.Separator.Left.Symbol,
					)
				}

				b.FullText += fmt.Sprintf(
					"<span color='%s' background='%s' font='%s' size= '%s'>%s</span>",
					Conf.NetIf.Color,
					Conf.NetIf.Background,
					Conf.NetIf.Font,
					Conf.NetIf.FontSize,
					IfStatus,
				)

				if Conf.NetIf.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.NetIf.Separator.Right.Color,
						Conf.NetIf.Separator.Right.Background,
						Conf.NetIf.Separator.Right.Font,
						Conf.NetIf.Separator.Right.FontSize,
						Conf.NetIf.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.Vpn.Enabled {
				var b I3BarOutBlock

				b.Color = Conf.Vpn.Color
				b.Background = Conf.Vpn.Background

				if Conf.Vpn.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Vpn.Separator.Left.Color,
						Conf.Vpn.Separator.Left.Background,
						Conf.Vpn.Separator.Left.Font,
						Conf.Vpn.Separator.Left.FontSize,
						Conf.Vpn.Separator.Left.Symbol,
					)
				}

				b.FullText += fmt.Sprintf(
					"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
					Conf.Vpn.Color,
					Conf.Vpn.Background,
					Conf.Vpn.Font,
					Conf.Vpn.FontSize,
					VPNStatus,
				)

				if Conf.Vpn.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Vpn.Separator.Right.Color,
						Conf.Vpn.Separator.Right.Background,
						Conf.Vpn.Separator.Right.Font,
						Conf.Vpn.Separator.Right.FontSize,
						Conf.Vpn.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.Battery.Enabled {
				var b I3BarOutBlock

				b.Color = Conf.Battery.Color
				b.Background = Conf.Battery.Background

				if Conf.Battery.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Battery.Separator.Left.Color,
						Conf.Battery.Separator.Left.Background,
						Conf.Battery.Separator.Left.Font,
						Conf.Battery.Separator.Left.FontSize,
						Conf.Battery.Separator.Left.Symbol,
					)
				}

				b.FullText += Batt

				if Conf.Battery.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Battery.Separator.Right.Color,
						Conf.Battery.Separator.Right.Background,
						Conf.Battery.Separator.Right.Font,
						Conf.Battery.Separator.Right.FontSize,
						Conf.Battery.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.SimpleVolumePa.Enabled {
				var b I3BarOutBlock

				b.Name = "simple-volume-pa"
				b.Color = Conf.SimpleVolumePa.Color
				b.Background = Conf.SimpleVolumePa.Background

				if Conf.SimpleVolumePa.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.SimpleVolumePa.Separator.Left.Color,
						Conf.SimpleVolumePa.Separator.Left.Background,
						Conf.SimpleVolumePa.Separator.Left.Font,
						Conf.SimpleVolumePa.Separator.Left.FontSize,
						Conf.SimpleVolumePa.Separator.Left.Symbol,
					)
				}

				// Pango format is already applied in plugin src.
				b.FullText += SoundVolume

				if Conf.SimpleVolumePa.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.SimpleVolumePa.Separator.Right.Color,
						Conf.SimpleVolumePa.Separator.Right.Background,
						Conf.SimpleVolumePa.Separator.Right.Font,
						Conf.SimpleVolumePa.Separator.Right.FontSize,
						Conf.SimpleVolumePa.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if Conf.Clock.Enabled {
				var b I3BarOutBlock

				b.Name = `wallclock`
				b.Color = Conf.Clock.Color
				b.Background = Conf.Clock.Background

				if Conf.Clock.Separator.Left.Enabled {
					b.FullText = fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Clock.Separator.Left.Color,
						Conf.Clock.Separator.Left.Background,
						Conf.Clock.Separator.Left.Font,
						Conf.Clock.Separator.Left.FontSize,
						Conf.Clock.Separator.Left.Symbol,
					)
				}

				b.FullText += fmt.Sprintf(
					"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
					Conf.Clock.Color,
					Conf.Clock.Background,
					Conf.Clock.Font,
					Conf.Clock.FontSize,
					Clock,
				)

				if Conf.Clock.Separator.Right.Enabled {
					b.FullText += fmt.Sprintf(
						"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
						Conf.Clock.Separator.Right.Color,
						Conf.Clock.Separator.Right.Background,
						Conf.Clock.Separator.Right.Font,
						Conf.Clock.Separator.Right.FontSize,
						Conf.Clock.Separator.Right.Symbol,
					)
				}

				b.Markup = "pango"
				b.Separator = false

				j = append(j, b)
			}

			if PrintOutput && len(j) > 0 {
				MsgChan <- j
			}
		}
	}
}

// PrintToI3bar prints info to stdout according to ipc docs (https://i3wm.org/docs/i3bar-protocol.html)
func PrintToI3bar() {
	for message := range MsgChan {
		// we do not need to html-encode output, json.Marshal does this forcefully, so invent our own Marshal
		buf := new(bytes.Buffer)
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(message)

		if err != nil {
			log.Printf("Unable to json-encode message, %s\n", err)
		}

		fmt.Println(strings.TrimSuffix(buf.String(), "\n") + ",")
	}
}
