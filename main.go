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
	BorderTop    int    `json:"border_top,omitempty"`
	BorderRight  int    `json:"border_right,omitempty"`
	BorderBottom int    `json:"border_bottom,omitempty"`
	BorderLeft   int    `json:"border_left,omitempty"`
	// measured either in pixels or in characters, so either int or string, let's make it string :)
	MinWidth            string `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Separator           bool   `json:"separator,omitempty"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
}

// UpdateReady channel that "refreshes" i3bar (generates stdout json line).
var UpdateReady = make(chan bool)

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

	if err != nil {
		log.Panicln(err)
	}

	go Spawner()
	go ParseStdin()
	go CleanZombies()

	if Conf.AppButtons {
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
	if Conf.LA {
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

			if Conf.AppButtons {
				for _, app := range Conf.Apps {
					var b I3BarOutBlock

					b.FullText = app.FullText
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

					j = append(j, b)
				}
			}

			if Conf.CPUTemp.Enabled {
				var b I3BarOutBlock
				b.FullText = fmt.Sprintf("CPU: %d°", CPUTemperature)
				j = append(j, b)
			}

			if Conf.Mem.Enabled {
				var b I3BarOutBlock

				if Conf.Mem.ShowSwap {
					b.FullText = fmt.Sprintf("M:%d%% SHM:%dM SW:%dM", Memory.Usedpct, Memory.Shared, Memory.Swap)
				} else {
					b.FullText = fmt.Sprintf("M:%d%% SHM:%dM", Memory.Usedpct, Memory.Shared)
				}

				j = append(j, b)
			}

			if Conf.LA {
				var b I3BarOutBlock
				b.FullText = fmt.Sprintf("La: %s", La)
				j = append(j, b)
			}

			if Conf.NetIf.Enabled {
				var b I3BarOutBlock
				b.FullText = IfStatus
				b.Markup = "pango"
				j = append(j, b)
			}

			if Conf.Vpn.Enabled {
				var b I3BarOutBlock
				b.FullText = VPNStatus
				b.Markup = "pango"
				j = append(j, b)
			}

			if Conf.Battery.Enabled {
				var b I3BarOutBlock
				b.FullText = Batt
				j = append(j, b)
			}

			if Conf.SimpleVolumePa.Enabled {
				var b I3BarOutBlock
				b.Name = "simple-volume-pa"
				b.FullText = SoundVolume
				b.Markup = "pango"
				j = append(j, b)
			}

			if Conf.Clock.Enabled {
				var b I3BarOutBlock
				b.Name = `wallclock`
				b.FullText = Clock
				b.Markup = "pango"
				b.Color = Conf.Clock.Color
				j = append(j, b)
			}

			if PrintOutput && len(j) > 0 {
				PrintToI3bar(j)
			}
		}
	}
}

// PrintToI3bar prints info to stdout according to ipc docs (https://i3wm.org/docs/i3bar-protocol.html)
func PrintToI3bar(message []I3BarOutBlock) {
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
