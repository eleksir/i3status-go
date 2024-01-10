package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/hjson/hjson-go"
)

type MyConfig struct {

	// Default text color
	Color string `json:"color,omitempty"`

	// Default background color
	Background string `json:"background,omitempty"`

	// LA Plugin
	LA bool `json:"la,omitempty"`

	Mem struct {
		Enabled  bool `json:"enabled,omitempty"`
		ShowSwap bool `json:"show_swap,omitempty"`
	} `json:"mem,omitempty"`

	Clock struct {
		Enabled bool   `json:"enabled,omitempty"`
		Color   string `json:"color,omitempty"`

		LeftClick struct {
			Enabled bool     `json:"enabled,omitempty"`
			Cmd     []string `json:"cmd,omitempty"`
		} `json:"left_click,omitempty"`

		RightClick struct {
			Enabled bool     `json:"enabled,omitempty"`
			Cmd     []string `json:"cmd,omitempty"`
		} `json:"right_click,omitempty"`
	} `json:"clock,omitempty"`

	Battery struct {
		Enabled bool `json:"enabled,omitempty"`
		Color   struct {
			Full        string `json:"full,omitempty"`
			Empty       string `json:"empty,omitempty"`
			AlmostFull  string `json:"almost_full,omitempty"`
			AlmostEmpty string `json:"almost_empty,omitempty"`
		} `json:"color,omitempty"`
	} `json:"battery,omitempty"`

	CPUTemp struct {
		Enabled bool     `json:"enabled,omitempty"`
		File    []string `json:"file,omitempty"`
	} `json:"cpu_temp,omitempty"`

	CapsLock struct {
		Enabled    bool   `json:"enabled,omitempty"`
		Background string `json:"background,omitempty"`
		Color      string `json:"color,omitempty"`
	} `json:"capslock,omitempty"`

	Vpn struct {
		Enabled        bool   `json:"enabled,omitempty"`
		StatusFile     string `json:"statusFile,omitempty"`
		MtimeThreshold int    `json:"mtime_threshold,omitempty"`
		DownColor      string `json:"down_color,omitempty"`
		UpColor        string `json:"up_color,omitempty"`

		TCPCheck struct {
			Enabled bool   `json:"enabled,omitempty"`
			Host    string `json:"host,omitempty"`
			Port    int    `json:"port,omitempty"`
			Timeout int    `json:"timeout,omitempty"`
		} `json:"tcp_check,omitempty"`
	} `json:"vpn,omitempty"`

	SimpleVolumePa struct {
		Enabled        bool     `json:"enabled,omitempty"`
		Symbol         string   `json:"symbol,omitempty"`
		Step           int      `json:"step,omitempty"`
		RightClickCmd  []string `json:"right_click_cmd,omitempty"`
		WheelUp        int      `json:"wheel_up,omitempty"`
		WheelDown      int      `json:"wheel_down,omitempty"`
		MaxVolumeLimit int      `json:"max_volume_limit,omitempty"`
	} `json:"simple-volume-pa,omitempty"`

	NetIf struct {
		Enabled   bool   `json:"enabled,omitempty"`
		DownColor string `json:"down_color,omitempty"`
		UpColor   string `json:"up_color,omitempty"`

		If []struct {
			Name string `json:"name,omitempty"`
			Dir  string `json:"dir,omitempty"`
		} `json:"if,omitempty"`
	} `json:"net-if,omitempty"`

	Cron struct {
		Enabled  bool   `json:"enabled,omitempty"`
		TimeZone string `json:"timezone,omitempty"`

		Tasks []struct {
			Time string   `json:"time,omitempty"`
			Cmd  []string `json:"cmd,omitempty"`
		} `json:"tasks,omitempty"`
	} `json:"cron,omitempty"`

	AppButtons bool `json:"app_buttons,omitempty"`

	Apps []struct {
		FullText            string   `json:"full_text,omitempty"`
		Name                string   `json:"name,omitempty"`
		Cmd                 string   `json:"cmd,omitempty"`
		Args                []string `json:"args,omitempty"`
		Instance            string   `json:"instance,omitempty"`
		Class               string   `json:"class,omitempty"`
		Color               string   `json:"color,omitempty"`
		Background          string   `json:"background,omitempty"`
		Border              string   `json:"border,omitempty"`
		BorderActive        string   `json:"border_active,omitempty"`
		Separator           bool     `json:"separator,omitempty"`
		SeparatorBlockWidth int      `json:"separator_block_width,omitempty"`
	} `json:"apps,omitempty"`
}

// Conf global variable with application config.
var Conf MyConfig

// readConf reads and validates confg if config does not exist, it puts default config to the same dir where i3 config
// is located.
func readConf() (MyConfig, error) {
	var (
		path   string
		config MyConfig
		err    error
		buf    []byte
	)

	path, err = LocateConfFile()

	if err != nil {
		return config, err
	}

	fileInfo, err := os.Stat(path)

	// Предполагаем, что файла либо нет, либо мы не можем его прочитать, второе надо бы логгировать, но пока забьём.
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			buf = DefaultConfig

			if err := os.WriteFile(path, buf, 0644); err != nil {
				return config, err
			}
		} else {
			return config, err
		}
	} else {
		// Конфиг-файл длинноват для конфига.
		if fileInfo.Size() > 65535 {
			err := fmt.Errorf("config file %s is too long for config", path) //nolint: goerr113

			return config, err
		}

		buf, err = os.ReadFile(path)

		// Не удалось прочитать
		if err != nil {
			err = fmt.Errorf("unable to read config file %s: %w", path, err)

			return config, err
		}
	}

	// According to docs, hjson seems can parse "quirky" json, but parses it to the map.
	// We interested in struct as output product: so we parse config to intermediate map then marshal it to json and
	// then produced json unmarshal to struct. Not very effective way, but in app lifetime it happens only once.
	var (
		sampleConfig MyConfig
		tmp          map[string]interface{}
	)

	err = hjson.Unmarshal(buf, &tmp)

	if err != nil {
		err := fmt.Errorf("unable to parse config file %s: %w", path, err)

		return config, err
	}

	tmpJSON, err := json.Marshal(tmp)

	if err != nil {
		err := fmt.Errorf("unable to parse config file %s: %w", path, err)

		return config, err
	}

	if err := json.Unmarshal(tmpJSON, &sampleConfig); err != nil {
		err := fmt.Errorf("unable to parse config file %s: %w", path, err)

		return config, err
	}

	// We're done with marshal-unmarshal config data, it is time to validate config.
	if sampleConfig.Color == "" {
		sampleConfig.Color = "#3e78fd"
	}

	if sampleConfig.Background == "" {
		sampleConfig.Background = "#edeceb"
	}

	// sampleConfig.La will be false if not set in config
	// sampleConfig.Mem will be false if not set in config
	// sampleConfig.Clock.Enabled will be false if not set in config

	if sampleConfig.Clock.Color == "" {
		sampleConfig.Clock.Color = "#666666"
	}

	// sampleConfig.Clock.LeftClick.Enabled will be false if not set in config

	if len(sampleConfig.Clock.LeftClick.Cmd) == 0 {
		sampleConfig.Clock.LeftClick.Cmd = append(sampleConfig.Clock.LeftClick.Cmd, "true")
	}

	// sampleConfig.Clock.RightClick.Enabled will be false if not set in config

	if len(sampleConfig.Clock.RightClick.Cmd) == 0 {
		sampleConfig.Clock.RightClick.Cmd = append(sampleConfig.Clock.RightClick.Cmd, "true")
	}

	// sampleConfig.battery.Enabled will be false if not set in config
	// sampleConfig.Battery.Color.Full will be empty string if not set
	// sampleConfig.Battery.Color.Empty will be empty string if not set
	// sampleConfig.Battery.Color.AlmostFull will be empty string if not set
	// sampleConfig.Battery.Color.AlmostEmpty will be empty string if not set

	// sampleConfig.CpuTemp.Enabled will be false if not set in config

	// No files configured - disable plugin
	if len(sampleConfig.CPUTemp.File) == 0 {
		sampleConfig.CPUTemp.Enabled = false
	}

	// sampleConfig.CapsLock.Enabled will false if not set in config

	if sampleConfig.CapsLock.Color == "" {
		sampleConfig.CapsLock.Color = "#3e78fd"
	}

	if sampleConfig.CapsLock.Background == "" {
		sampleConfig.CapsLock.Background = "#edeceb"
	}

	// sampleConfig.Vpn.Enabled will false if not set in config

	// No status file - disable plugin
	if sampleConfig.Vpn.StatusFile == "" {
		sampleConfig.Vpn.Enabled = false
	}

	// Check status file at least once per 3 seconds
	if sampleConfig.Vpn.MtimeThreshold < 3 {
		if sampleConfig.Vpn.Enabled {
			log.Printf("vpn.mtime_threshold not set, using 3")
		}

		sampleConfig.Vpn.MtimeThreshold = 3 //nocritic: wsl
	}

	// sampleConfig.Vpn.DownColor will be empty string if no value set in config
	// sampleConfig.Vpn.UpColor will be empty string if no value set in config
	// sampleConfig.Vpn.TcpCheck.Enabled will false if not set in config

	// Disable plugin if value inadequate
	if sampleConfig.Vpn.TCPCheck.Port > 65535 {
		if sampleConfig.Vpn.TCPCheck.Enabled {
			log.Printf("vpn.tcp_check.port > 65535, disabling vpn.tcp_check.enabled")
		}

		sampleConfig.Vpn.TCPCheck.Enabled = false
	}

	// Disable plugin if value indaquate
	if sampleConfig.Vpn.TCPCheck.Host == "" {
		if sampleConfig.Vpn.TCPCheck.Enabled {
			log.Printf("vpn.tcp_check.host not set, disabling vpn.tcp_check.enabled")
		}

		sampleConfig.Vpn.TCPCheck.Enabled = false
	}

	if sampleConfig.Vpn.TCPCheck.Timeout < 3 {
		if sampleConfig.Vpn.TCPCheck.Enabled {
			log.Printf("vpn.tcp_check.timeout < 3, using 3")
		}

		sampleConfig.Vpn.TCPCheck.Timeout = 3
	}

	// TODO: appbuttons, etc..
	// sampleConfig.SimpleVolumePa.Enabled will false if not set in config

	if sampleConfig.SimpleVolumePa.Symbol == "" {
		sampleConfig.SimpleVolumePa.Symbol = `🔊`
	}

	if sampleConfig.SimpleVolumePa.Step == 0 {
		sampleConfig.SimpleVolumePa.Step = 5
	}

	if sampleConfig.SimpleVolumePa.WheelUp == 0 {
		sampleConfig.SimpleVolumePa.WheelUp = 4
	}

	if sampleConfig.SimpleVolumePa.WheelDown == 0 {
		sampleConfig.SimpleVolumePa.WheelDown = 5
	}

	if sampleConfig.SimpleVolumePa.MaxVolumeLimit == 0 {
		sampleConfig.SimpleVolumePa.MaxVolumeLimit = 100
	}

	if len(sampleConfig.SimpleVolumePa.RightClickCmd) > 0 {
		if sampleConfig.SimpleVolumePa.RightClickCmd[0] == "" {
			sampleConfig.SimpleVolumePa.RightClickCmd[0] = "true"
		}
	} else {
		sampleConfig.SimpleVolumePa.RightClickCmd = append(sampleConfig.SimpleVolumePa.RightClickCmd, "true")
	}

	// sampleConfig.Cron.Enabled will false if not set in config

	if sampleConfig.Cron.TimeZone == "" {
		sampleConfig.Cron.TimeZone = "GMT+0"
	}

	if len(sampleConfig.Cron.Tasks) == 0 {
		sampleConfig.Cron.Enabled = false
	}

	// sampleConfig.NetIf.Enabled will false if not set in config

	if len(sampleConfig.NetIf.If) == 0 {
		sampleConfig.NetIf.Enabled = false
	}

	if sampleConfig.NetIf.DownColor == "" {
		sampleConfig.NetIf.DownColor = "red"
	}

	if sampleConfig.NetIf.UpColor == "" {
		sampleConfig.NetIf.UpColor = "green"
	}

	// sampleConfig.AppButtons will false if not set in config

	if len(sampleConfig.Apps) == 0 {
		sampleConfig.AppButtons = false
	} else {
		for num, app := range sampleConfig.Apps {
			// app.Instance can be missing.
			// app.Class can be missing.

			if app.Name == "" {
				app.Name = fmt.Sprintf("app%d", num)
			}

			// app.Args can be empty slice. In that case command will be run without aruments.

			if app.Cmd == "" {
				app.Cmd = "true"
			}

			if app.Color == "" {
				app.Color = Conf.Color
			}

			if app.Background == "" {
				app.Background = Conf.Background
			}

			if app.Border == "" {
				app.Border = Conf.Color
			}

			if app.BorderActive == "" {
				app.BorderActive = Conf.Color
			}

			if app.FullText == "" {
				app.FullText = fmt.Sprintf(" %d ", num)
			}

			// app.Separator can be omitted, in that case it is false
			// app.SeparatorBlockWidth can be missing
		}
	}

	config = sampleConfig

	return config, nil
}

// LocateConfFile make a try to locate i3 config dir.
func LocateConfFile() (string, error) {
	i3ConfigPath, err := xdg.SearchConfigFile("i3/config")

	if err != nil {
		err = fmt.Errorf("unable to find i3 config: %w", err)

		return "", err
	}

	return filepath.Dir(i3ConfigPath) + "/i3status-go.json", nil
}
