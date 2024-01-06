package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

// ClickEvent struct as described at https://i3wm.org/docs/i3bar-protocol.html
type ClickEvent struct {
	Name      string   `json:"name,omitempty"`
	Instance  string   `json:"instance,omitempty"`
	Button    int      `json:"button"`
	Modifiers []string `json:"modifiers"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	RelativeX int      `json:"relative_x"`
	RelativeY int      `json:"relative_y"`
	OutputX   int      `json:"output_x"`
	OutputY   int      `json:"output_y"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
}

// ParseStdin tries to parse text that i3bar prints to our stdin. Currently - it is mouse click events on different
// area names of i3bar.
func ParseStdin() {
	reader := bufio.NewReader(os.Stdin)

	// De-facto it is jsonl, except first line is garbage. Also, first symbol in each strint garbage too.
	firstelem := true

	for {
		var e ClickEvent

		buf, err := reader.ReadString('\n')

		if err != nil {
			if !firstelem {
				log.Printf("Unable to read from stdin: %s", err)
			}

			break
		}

		if len(buf) == 0 {
			log.Print("String from stdin has zero length, skipping")

			continue
		}

		for i := 0; i < len(buf); i++ {
			if strings.HasPrefix(buf, "{") {
				break
			}

			buf = buf[1:]
		}

		if len(buf) == 0 {
			log.Print("String from stdin has zero length, skipping")

			continue
		}

		err = json.Unmarshal([]byte(buf), &e)

		if err != nil {
			if firstelem {
				firstelem = false
			} else {
				log.Printf("Unable to decode json line from stdin: %s, buf: %s", err, buf)
			}

			continue
		}

		// Just in case.
		if firstelem {
			firstelem = false
		}

		// Clock clicks parse.
		if e.Name == "wallclock" {
			if Conf.Clock.LeftClick.Enabled && e.Button == 1 {
				RunChan <- Conf.Clock.LeftClick.Cmd
			} else if Conf.Clock.RightClick.Enabled && e.Button == 3 {
				RunChan <- Conf.Clock.RightClick.Cmd
			}

			continue
		}

		if e.Name == "simple-volume-pa" {
			if Conf.SimpleVolumePa.Enabled {
				if e.Button == 3 {
					RunChan <- Conf.SimpleVolumePa.RightClickCmd

					continue
				}

				vol, err := pa.Volume()

				if err != nil {
					if err := PaReinit(); err != nil {
						log.Printf("Unable to get pulseaudio volume: %s", err)
					} else {
						vol, err = pa.Volume()

						if err != nil {
							log.Printf("Unable to get volume pulseaudio server behaves weirdly: %s", err)
						}
					}
				}

				switch e.Button {
				case Conf.SimpleVolumePa.WheelUp:
					vol += float32(Conf.SimpleVolumePa.Step) / 100

					if vol > (float32(Conf.SimpleVolumePa.MaxVolumeLimit) / 100) {
						vol = float32(Conf.SimpleVolumePa.MaxVolumeLimit) / 100
					}

					if err := pa.SetVolume(vol); err != nil {
						log.Printf("Unable to set pulseaudio volume: %s", err)
					}

				case Conf.SimpleVolumePa.WheelDown:
					vol -= float32(Conf.SimpleVolumePa.Step) / 100

					if vol < 0 {
						vol = 0
					}

					if err := pa.SetVolume(vol); err != nil {
						log.Printf("Unable to set pulseaudio volume: %s", err)
					}
				}
			}

			continue
		}

		if !Conf.AppButtons {
			continue
		}

		for _, app := range Conf.Apps {
			switch {
			case app.Name != "" && app.Instance != "":
				if app.Name == e.Name && app.Instance == e.Instance {
					prg := append([]string{}, app.Cmd)
					prg = append(prg, app.Args...)

					RunChan <- prg

					break
				}
			case app.Name != "" && e.Name == app.Name:
				prg := append([]string{}, app.Cmd)
				prg = append(prg, app.Args...)

				RunChan <- prg

			case app.Instance != "" && e.Instance == app.Instance:
				prg := append([]string{}, app.Cmd)
				prg = append(prg, app.Args...)

				RunChan <- prg
			}
		}
	}
}
