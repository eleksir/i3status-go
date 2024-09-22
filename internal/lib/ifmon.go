package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// IfStatus global var conaining cumulative status of all monitored net intefaces.
var IfStatus string

// UpdateIfStatus updates network interfaces status for i3bar.
func (c MyConfig) UpdateIfStatus() {
	for {
		var statusSum string

		for _, item := range Conf.NetIf.If {
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
					StatusStr := fmt.Sprintf("<span foreground=\"%s\">⍋</span>", Conf.NetIf.UpColor)
					operstate = []byte(StatusStr)
				case "down":
					statusStr := fmt.Sprintf("<span foreground=\"%s\">⍒</span>", Conf.NetIf.DownColor)
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

		if c.IfStatus != statusSum {
			c.IfStatus = statusSum
			c.UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}
