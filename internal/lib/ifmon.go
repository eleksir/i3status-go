package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UpdateIfStatus updates network interfaces status for i3bar.
func (c *MyConfig) UpdateIfStatus() {
	var (
		InitialDelay       = 100 * time.Millisecond
		LoopIterationDelay = 3 * time.Second
		Delay              = InitialDelay
		ticker             = time.NewTicker(Delay)
	)

	for range ticker.C {
		if Delay == InitialDelay {
			Delay = LoopIterationDelay
			ticker.Reset(Delay)
		}

		var statusSum string

		for _, item := range c.NetIf.If {
			var (
				name = item.Name
			)

			if name == "" {
				name = filepath.Base(item.Dir)
			}

			operstate, err := os.ReadFile(item.Dir + "/operstate")

			if err != nil {
				operstate = []byte("?")

				log.Printf("Unable to get net if status from file %s: %s", item.Dir+"/operstate", err)
			} else {
				switch strings.TrimSpace(string(operstate)) {
				case "up":
					StatusStr := fmt.Sprintf("<span foreground=\"%s\">⍋</span>", c.NetIf.UpColor)
					operstate = []byte(StatusStr)
				case "down":
					statusStr := fmt.Sprintf("<span foreground=\"%s\">⍒</span>", c.NetIf.DownColor)
					operstate = []byte(statusStr)
				default:
					operstate = []byte(`?`)
				}
			}

			if statusSum != "" {
				statusSum += " "
			}

			statusSum += fmt.Sprintf("%s:%s", name, operstate)
		}

		if c.Values.IfStatus != statusSum {
			c.Values.IfStatus = statusSum
			c.Channels.UpdateReady <- true
		}
	}
}
