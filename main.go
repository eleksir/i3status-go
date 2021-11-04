package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

// Config is a structure that contains parsed config file data
type Config struct {
	Color      string `json:"color"`
	Background string `json:"background"`
	La         int    `json:"la"`
	Mem        int    `json:"mem"`

	Clock struct {
		Enabled int    `json:"enabled"`
		Color   string `json:"color"`

		LeftClick struct {
			Enabled int    `json:"enabled"`
			Cmd     string `json:"cmd"`
		} `json:"left_click"`

		RightClick struct {
			Enabled int    `json:"enabled"`
			Cmd     string `json:"cmd"`
		} `json:"right_click"`
	} `json:"clock"`

	Battery struct {
		Enabled      int    `json:"enabled"`
		Driver       string `json:"driver"`
		SysDir       string `json:"sys_dir"`
		UpowerDevice string `json:"upower_device"`
	} `json:"battery"`

	CPUTemp struct {
		Enabled int    `json:"enabled"`
		File    string `json:"file"`
	} `json:"cputemp"`

	CapsLock struct {
		Enabled    int    `json:"enabled"`
		Background string `json:"background"`
		Color      string `json:"color"`
	} `json:"capslock"`

	Vpn struct {
		Enabled        int    `json:"enabled"`
		DownColor      string `json:"down_color"`
		UpColor        string `json:"up_color"`
		StatusFile     string `json:"statusfile"`
		MtimeThreshold int    `json:"mtime_threshold"`
		TCPCheck       struct {
			Enabled int    `json:"enabled"`
			Host    string `json:"host"`
			Port    string `json:"port"`
			Timeout int    `json:"timeout"`
		} `json:"tcp_check"`
	} `json:"vpn"`

	SimpleVolumePa struct {
		Enabled int    `json:"enabled"`
		Symbol  string `json:"symbol"`
	} `json:"simple-volume-pa"`

	NetIf struct {
		Enabled   int    `json:"enabled"`
		DownColor string `json:"down_color"`
		UpColor   string `json:"up_color"`

		If []struct {
			Name string `json:"name"`
			Dir  string `json:"dir"`
		} `json:"if"`
	} `json:"netif"`

	Cron struct {
		Enabled  int    `json:"enabled"`
		TimeZone string `json:"timezone"`

		Tasks []struct {
			Time string   `json:"time"`
			Cmd  []string `json:"cmd"`
		} `json:"tasks"`
	} `json:"cron"`

	AppButtons int `json:"app_buttons"`

	Apps []struct {
		FullText            string `json:"full_text"`
		Name                string `json:"name"`
		Cmd                 string `json:"cmd"`
		Instance            string `json:"instance"`
		Class               string `json:"class"`
		Color               string `json:"color"`
		Border              string `json:"border"`
		BorderActive        string `json:"border_active"`
		Separator           bool   `json:"separator"`
		SeparatorBlockWidth int    `json:"separator_block_width"`
	} `json:"apps"`
}

type Mem struct {
	Usedpct int64
	Shared  int64
	Swap    int64
}

// I3BarOutBlock is structure element for I3BarOut, it represents i3bar output json block format
type I3BarOutBlock struct {
	FullText string `json:"full_text"`
	// ShortText will be shown if not enough room for FullText, threshold width defined in MinWidth
	ShortText    string `json:"short_text,omitempty"`
	Color        string `json:"color,omitempty"`
	Background   string `json:"background,omitempty"`
	Border       int    `json:"border,omitempty"`
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

// Conf is config file structure
var Conf Config

// UpdateReady triggers stdout json output
var UpdateReady = make(chan bool)

// Memory stores statistics about memory
var Memory Mem

// La stores load average for 1 minute
var La int64 = -1

// CPUTemperature stores current cpu temperature for given cpu core
var CPUTemperature int64

// Clock shows current date and time
var Clock = "Thu, 1 Jan 1970   1:00"

// Batt shows current battery status
var Batt = "<big>⚡</big> ??% •"

// Program entry point
func main() {
	loadConfig()

	// Populate memory stats
	if Conf.Mem == 1 {
		go UpdateMemStats()
	}

	// Populate LA stats
	if Conf.La == 1 {
		go UpdateLaStats()
	}

	// Populate CPUTemperature
	if Conf.CPUTemp.Enabled == 1 {
		go UpdateCPUTemperature()
	}

	// Populate Clock
	if Conf.Clock.Enabled == 1 {
		go UpdateClock()
	}

	// Populate Battery stats
	if Conf.Battery.Enabled == 1 {
		go UpdateBattery()
	}

	/*
		I3bar documentation pretends that message protocol must be valid json. In practice, we only have to print valid
		header, empty json array and (potentially infinite) json lines (line that is valid json by itself) that is
		actually json arrays. We do not need to *close* this json at all.
		Gracefully closed json required when i3bar initiates our program to stop|quit, this (should) happens just before
		i3bar itself terminating. So we don't care.
	*/

	// Print header and one empty message and wait for updates
	fmt.Printf("{\"version\": 1, \"click_events\": true}\n")
	fmt.Println("[ [],")

	for {
		if <-UpdateReady {
			// actually build json struct, marshal it and print result
			var j []I3BarOutBlock

			if Conf.CPUTemp.Enabled == 1 {
				var b I3BarOutBlock
				b.FullText = fmt.Sprintf("%d°", CPUTemperature)
				j = append(j, b)
			}

			if Conf.Mem == 1 {
				var b I3BarOutBlock
				b.FullText = fmt.Sprintf("M:%d%% SHM:%dM SW:%dM", Memory.Usedpct, Memory.Shared, Memory.Swap)
				j = append(j, b)
			}

			if Conf.La == 1 {
				var b I3BarOutBlock
				b.FullText = fmt.Sprintf("La: %d", La)
				j = append(j, b)
			}

			if Conf.Clock.Enabled == 1 {
				var b I3BarOutBlock
				b.FullText = Clock
				b.Markup = "pango"
				b.Color = Conf.Clock.Color
				j = append(j, b)
			}

			if len(j) > 0 {
				PrintToI3bar(j)
			}
		}
	}
}

// Fills Config struct with configuration stored in config file
func loadConfig() {
	var MyConfig Config

	f, _ := os.Open("i3status-go.json")
	r := JsonConfigReader.New(f)
	err := json.NewDecoder(r).Decode(&MyConfig)

	if err != nil {
		log.Fatalf("Unable to parse config, %s\n", err)
	}

	Conf = MyConfig
}

// UpdateMemStats Fills Memory statistics
func UpdateMemStats() {
	for {
		v, _ := mem.VirtualMemory()

		if Memory.Usedpct != int64(v.UsedPercent) || Memory.Shared != int64(v.Shared/1024/1024) ||
			Memory.Swap != int64(v.SwapTotal-v.SwapFree) {
			Memory.Usedpct = int64(v.UsedPercent)
			Memory.Shared = int64(v.Shared / 1024 / 1024)
			Memory.Swap = int64(v.SwapTotal - v.SwapFree)
			UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}

// UpdateLaStats Fills Load Average value for last 1 minute period
func UpdateLaStats() {
	for {
		l, _ := load.Avg()

		if La != int64(l.Load1) {
			La = int64(l.Load1)
			UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}

// UpdateCPUTemperature Fills CPU Temperature Measured by kernel
func UpdateCPUTemperature() {
	for {
		file, err := os.Open(Conf.CPUTemp.File)

		if err != nil {
			log.Printf("Unable to open %s: %s", Conf.CPUTemp.File, err)
		} else {
			reader := bufio.NewReader(file)
			line, _, err := reader.ReadLine()

			if err != nil {
				log.Printf("Unable to read from %s: %s", Conf.CPUTemp.File, err)
				err = file.Close()

				if err != nil {
					log.Printf("Unable to close %s: %s", Conf.CPUTemp.File, err)
				}
			} else {
				err = file.Close()

				if err != nil {
					log.Printf("Unable to close %s: %s", Conf.CPUTemp.File, err)
				} else {
					temp, err := strconv.ParseInt(string(line), 10, 32)

					if err != nil {
						log.Printf("Unable to convert string to number from file %s: %s", Conf.CPUTemp.File, err)
					} else {
						if temp > 1000 {
							temp /= 1000
						}

						if CPUTemperature != temp {
							CPUTemperature = temp
							UpdateReady <- true
						}
					}
				}
			}
		}

		time.Sleep(3 * time.Second)
	}
}

// UpdateClock Fills date and time from system clock
func UpdateClock() {
	for {
		currentTime := time.Now()
		hours, minutes, _ := currentTime.Clock()
		year, month, day := currentTime.Date()
		dow := currentTime.Weekday()
		rmonth := [12]string{"Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"}
		rdow := [7]string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

		myclock := fmt.Sprintf("<big>     %s, %d %s %d  % d:%0d  </big>", rdow[dow-1], day, rmonth[month-1], year, hours, minutes)

		if myclock != Clock {
			Clock = myclock
			UpdateReady <- true
		}

		time.Sleep(1 * time.Second)
	}
}

// UpdateBattery Fills battery statistics for power source given in config
func UpdateBattery() {
	for {
		// steal from https://github.com/soumya92/barista/blob/0eb8431fc7bbdc9e36602a9f73a42acae111e958/modules/battery/battery.go#L310
		Batt = "<big>⚡</big> ??% •"
		time.Sleep(3 * time.Second)
	}
}

// PrintToI3bar prints to stdout json string in i3bar protocol (https://i3wm.org/docs/i3bar-protocol.html)
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
