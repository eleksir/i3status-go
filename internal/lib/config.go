package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/adrg/xdg"
	"github.com/hjson/hjson-go"
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

type Separator struct {
	Left struct {
		Enabled    bool   `json:"enabled,omitempty"`
		Color      string `json:"color,omitempty"`
		Background string `json:"background,omitempty"`
		Symbol     string `json:"symbol,omitempty"`
		Font       string `json:"font,omitempty"`
		FontSize   string `json:"font_size,omitempty"`
	} `json:"left,omitempty"`

	Right struct {
		Enabled    bool   `json:"enabled,omitempty"`
		Color      string `json:"color,omitempty"`
		Background string `json:"background,omitempty"`
		Symbol     string `json:"symbol,omitempty"`
		Font       string `json:"font,omitempty"`
		FontSize   string `json:"font_size,omitempty"`
	} `json:"right,omitempty"`
}

type MyConfig struct {
	UpdateReady    chan bool
	PrintOutput    bool
	MsgChan        chan []I3BarOutBlock
	SigChan        chan os.Signal
	CPUTemperature int64
	BatteryString  string
	ClockTime      string
	IfStatus       string
	VPNStatus      string
	Memory         Mem
	La             string

	// Default text color
	Color string `json:"color,omitempty"`

	// Default background color
	Background string `json:"background,omitempty"`

	// Default text font
	Font string `json:"font,omitempty"`

	// Default text font size
	FontSize string `json:"font_size,omitempty"`

	Separator Separator `json:"separator,omitempty"`

	// LA Plugin
	LA struct {
		Enabled    bool      `json:"enabled,omitempty"`
		Color      string    `json:"color,omitempty"`
		Background string    `json:"background,omitempty"`
		Font       string    `json:"font,omitempty"`
		FontSize   string    `json:"font_size,omitempty"`
		Separator  Separator `json:"separator,omitempty"`
	} `json:"la,omitempty"`

	// Mem plugin
	Mem struct {
		Enabled    bool      `json:"enabled,omitempty"`
		Color      string    `json:"color,omitempty"`
		Background string    `json:"background,omitempty"`
		Font       string    `json:"font,omitempty"`
		FontSize   string    `json:"font_size,omitempty"`
		ShowSwap   bool      `json:"show_swap,omitempty"`
		Separator  Separator `json:"separator,omitempty"`
	} `json:"mem,omitempty"`

	Clock struct {
		Enabled    bool      `json:"enabled,omitempty"`
		Color      string    `json:"color,omitempty"`
		Background string    `json:"background,omitempty"`
		Font       string    `json:"font,omitempty"`
		FontSize   string    `json:"font_size,omitempty"`
		Separator  Separator `json:"separator,omitempty"`

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
		Enabled        bool   `json:"enabled,omitempty"`
		Color          string `json:"color,omitempty"`
		Background     string `json:"background,omitempty"`
		Font           string `json:"font,omitempty"`
		FontSize       string `json:"font_size,omitempty"`
		Symbol         string `json:"symbol,omitempty"`
		SymbolFont     string `json:"symbol_font,omitempty"`
		SymbolFontSize string `json:"symbol_font_size,omitempty"`

		ChargeColor struct {
			Full        string `json:"full,omitempty"`
			Empty       string `json:"empty,omitempty"`
			AlmostFull  string `json:"almost_full,omitempty"`
			AlmostEmpty string `json:"almost_empty,omitempty"`
		} `json:"charge_color,omitempty"`

		Separator Separator `json:"separator,omitempty"`
	} `json:"battery,omitempty"`

	CPUTemp struct {
		Enabled    bool      `json:"enabled,omitempty"`
		Color      string    `json:"color,omitempty"`
		Background string    `json:"background,omitempty"`
		Font       string    `json:"font,omitempty"`
		FontSize   string    `json:"font_size,omitempty"`
		Separator  Separator `json:"separator,omitempty"`

		File []string `json:"file,omitempty"`
	} `json:"cpu_temp,omitempty"`

	Vpn struct {
		Enabled        bool      `json:"enabled,omitempty"`
		Color          string    `json:"color,omitempty"`
		Background     string    `json:"background,omitempty"`
		Font           string    `json:"font,omitempty"`
		FontSize       string    `json:"font_size,omitempty"`
		StatusFile     string    `json:"statusFile,omitempty"`
		MtimeThreshold int       `json:"mtime_threshold,omitempty"`
		DownColor      string    `json:"down_color,omitempty"`
		UpColor        string    `json:"up_color,omitempty"`
		Separator      Separator `json:"separator,omitempty"`

		TCPCheck struct {
			Enabled bool   `json:"enabled,omitempty"`
			Host    string `json:"host,omitempty"`
			Port    int    `json:"port,omitempty"`
			Timeout int    `json:"timeout,omitempty"`
		} `json:"tcp_check,omitempty"`
	} `json:"vpn,omitempty"`

	SimpleVolumePa struct {
		Enabled         bool      `json:"enabled,omitempty"`
		Color           string    `json:"color,omitempty"`
		Background      string    `json:"background,omitempty"`
		Font            string    `json:"font,omitempty"`
		FontSize        string    `json:"font_size,omitempty"`
		Symbol          string    `json:"symbol,omitempty"`
		SymbolFont      string    `json:"symbol_font,omitempty"`
		SymbolFontSize  string    `json:"symbol_font_size,omitempty"`
		DontExitOnLogin bool      `json:"dont_exit_on_login,omitempty"`
		Separator       Separator `json:"separator,omitempty"`

		Step           int      `json:"step,omitempty"`
		RightClickCmd  []string `json:"right_click_cmd,omitempty"`
		WheelUp        int      `json:"wheel_up,omitempty"`
		WheelDown      int      `json:"wheel_down,omitempty"`
		MaxVolumeLimit int      `json:"max_volume_limit,omitempty"`
	} `json:"simple-volume-pa,omitempty"`

	NetIf struct {
		Enabled    bool      `json:"enabled,omitempty"`
		DownColor  string    `json:"down_color,omitempty"`
		UpColor    string    `json:"up_color,omitempty"`
		Color      string    `json:"color,omitempty"`
		Background string    `json:"background,omitempty"`
		Font       string    `json:"font,omitempty"`
		FontSize   string    `json:"font_size,omitempty"`
		Separator  Separator `json:"separator,omitempty"`

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

	AppButtons struct {
		Enabled   bool      `json:"enabled,omitempty"`
		Separator Separator `json:"separator,omitempty"`
	} `json:"app_buttons,omitempty"`

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

// readConf reads and validates config if config does not exist, it puts default config to the same dir where i3 config
// is located.
func ReadConf(DefaultConfig []byte) (MyConfig, error) {
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

	// Assume that either we unable to read file or file does not exit. We should mention second in logs but forget
	// about it for now.
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
		// Config file looks too long for config...
		if fileInfo.Size() > 65535 {
			err := fmt.Errorf("config file %s is too long for config", path) //nolint: goerr113

			return config, err
		}

		buf, err = os.ReadFile(path)

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

	if sampleConfig.Font == "" {
		sampleConfig.Font = "Liberation Mono"
	}

	if sampleConfig.FontSize == "" {
		sampleConfig.FontSize = "medium"
	}

	// Even if separator is disabled fallback values should be filled in.
	if sampleConfig.Separator.Left.Color == "" {
		sampleConfig.Separator.Left.Color = sampleConfig.Color
	}

	if sampleConfig.Separator.Left.Background == "" {
		sampleConfig.Separator.Left.Background = sampleConfig.Background
	}

	if sampleConfig.Separator.Left.Symbol == "" {
		sampleConfig.Separator.Left.Symbol = "|"
	}

	if sampleConfig.Separator.Left.Font == "" {
		sampleConfig.Separator.Left.Font = sampleConfig.Font
	}

	matched, err := regexp.MatchString(
		`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
		sampleConfig.Separator.Left.FontSize,
	)

	if err != nil {
		log.Printf(
			"Unable to set sampleConfig.Separator.Left.FontSize: %s, fallback to %s",
			err,
			sampleConfig.FontSize,
		)

		sampleConfig.Separator.Left.FontSize = sampleConfig.FontSize
	}

	if !matched {
		log.Printf(
			"Unable to set sampleConfig.Separator.Left.FontSize: %s, fallback to %s",
			err,
			sampleConfig.FontSize,
		)

		sampleConfig.Separator.Left.FontSize = sampleConfig.FontSize
	}

	// Even if separator is disabled fallback values should be filled in.
	if sampleConfig.Separator.Right.Color == "" {
		sampleConfig.Separator.Right.Color = sampleConfig.Color
	}

	if sampleConfig.Separator.Right.Background == "" {
		sampleConfig.Separator.Right.Background = sampleConfig.Background
	}

	if sampleConfig.Separator.Right.Symbol == "" {
		sampleConfig.Separator.Right.Symbol = "|"
	}

	if sampleConfig.Separator.Right.Font == "" {
		sampleConfig.Separator.Right.Font = sampleConfig.Font
	}

	matched, err = regexp.MatchString(
		`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
		sampleConfig.Separator.Right.FontSize,
	)

	if err != nil {
		log.Printf(
			"Unable to set sampleConfig.Separator.Right.FontSize: %s, fallback to %s",
			err,
			sampleConfig.FontSize,
		)

		sampleConfig.Separator.Right.FontSize = sampleConfig.FontSize
	}

	if !matched {
		log.Printf("Unable to set sampleConfig.Separator.Right.FontSize: %s, fallback to %s",
			err,
			sampleConfig.FontSize,
		)

		sampleConfig.Separator.Right.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.LA.Color == "" {
		sampleConfig.LA.Color = sampleConfig.Color
	}

	if sampleConfig.LA.Background == "" {
		sampleConfig.LA.Background = sampleConfig.Background
	}

	if sampleConfig.LA.Font == "" {
		sampleConfig.LA.Font = sampleConfig.Font
	}

	if sampleConfig.LA.FontSize == "" {
		sampleConfig.LA.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.LA.Separator.Left.Color == "" {
		sampleConfig.LA.Separator.Left.Color = sampleConfig.Separator.Left.Color
	}

	if sampleConfig.LA.Separator.Left.Background == "" {
		sampleConfig.LA.Separator.Left.Background = sampleConfig.Separator.Left.Background
	}

	if sampleConfig.LA.Separator.Left.Symbol == "" {
		sampleConfig.LA.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
	}

	if sampleConfig.LA.Separator.Left.Font == "" {
		sampleConfig.LA.Separator.Left.Font = sampleConfig.Separator.Left.Font
	}

	if sampleConfig.LA.Separator.Left.FontSize == "" {
		sampleConfig.LA.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.LA.Separator.Left.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.LA.Separator.Left.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.LA.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.LA.Separator.Left.FontSize, fallback to %s",
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.LA.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}
	}

	if sampleConfig.LA.Separator.Right.Color == "" {
		sampleConfig.LA.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.LA.Separator.Right.Background == "" {
		sampleConfig.LA.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.LA.Separator.Right.Font == "" {
		sampleConfig.LA.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.LA.Separator.Right.FontSize == "" {
		sampleConfig.LA.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	}

	if sampleConfig.LA.Separator.Right.Color == "" {
		sampleConfig.LA.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.LA.Separator.Right.Background == "" {
		sampleConfig.LA.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.LA.Separator.Right.Symbol == "" {
		sampleConfig.LA.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.LA.Separator.Right.Font == "" {
		sampleConfig.LA.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.LA.Separator.Right.FontSize == "" {
		sampleConfig.LA.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.LA.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.LA.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.LA.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.LA.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.LA.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	if sampleConfig.Mem.Color == "" {
		sampleConfig.Mem.Color = sampleConfig.Color
	}

	if sampleConfig.Mem.Background == "" {
		sampleConfig.Mem.Background = sampleConfig.Background
	}

	if sampleConfig.Mem.Font == "" {
		sampleConfig.Mem.Font = sampleConfig.Font
	}

	if sampleConfig.Mem.FontSize == "" {
		sampleConfig.Mem.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.Mem.Separator.Left.Color == "" {
		sampleConfig.Mem.Separator.Left.Color = sampleConfig.Separator.Left.Color
	}

	if sampleConfig.Mem.Separator.Left.Background == "" {
		sampleConfig.Mem.Separator.Left.Background = sampleConfig.Separator.Left.Background
	}

	if sampleConfig.Mem.Separator.Left.Symbol == "" {
		sampleConfig.Mem.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
	}

	if sampleConfig.Mem.Separator.Left.Font == "" {
		sampleConfig.Mem.Separator.Left.Font = sampleConfig.Separator.Left.Font
	}

	if sampleConfig.Mem.Separator.Left.FontSize == "" {
		sampleConfig.Mem.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.Mem.Separator.Left.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.Mem.Separator.Left.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.Mem.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.Mem.Separator.Left.FontSize, fallback to %s",
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.Mem.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}
	}

	if sampleConfig.Mem.Separator.Right.Color == "" {
		sampleConfig.Mem.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.Mem.Separator.Right.Background == "" {
		sampleConfig.Mem.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.Mem.Separator.Right.Symbol == "" {
		sampleConfig.Mem.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.Mem.Separator.Right.Font == "" {
		sampleConfig.Mem.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.Mem.Separator.Right.FontSize == "" {
		sampleConfig.Mem.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.Mem.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.Mem.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Mem.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.Mem.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Mem.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	if sampleConfig.Clock.Color == "" {
		sampleConfig.Clock.Color = sampleConfig.Color
	}

	if sampleConfig.Clock.Background == "" {
		sampleConfig.Clock.Background = sampleConfig.Background
	}

	if sampleConfig.Clock.Font == "" {
		sampleConfig.Clock.Font = sampleConfig.Font
	}

	if sampleConfig.Clock.FontSize == "" {
		sampleConfig.Clock.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.Clock.Separator.Left.Enabled {
		if sampleConfig.Clock.Separator.Left.Color == "" {
			sampleConfig.Clock.Separator.Left.Color = sampleConfig.Separator.Left.Color
		}

		if sampleConfig.Clock.Separator.Left.Background == "" {
			sampleConfig.Clock.Separator.Left.Background = sampleConfig.Separator.Left.Background
		}

		if sampleConfig.Clock.Separator.Left.Symbol == "" {
			sampleConfig.Clock.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
		}

		if sampleConfig.Clock.Separator.Left.Font == "" {
			sampleConfig.Clock.Separator.Left.Font = sampleConfig.Separator.Left.Font
		}

		if sampleConfig.Clock.Separator.Left.FontSize == "" {
			sampleConfig.Clock.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		} else {
			matched, err := regexp.MatchString(
				`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
				sampleConfig.Clock.Separator.Left.FontSize,
			)

			if err != nil {
				log.Printf(
					"Unable to set sampleConfig.Clock.Separator.Left.FontSize: %s, fallback to %s",
					err,
					sampleConfig.Separator.Left.FontSize,
				)

				sampleConfig.Clock.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
			}

			if !matched {
				log.Printf(
					"Unable to set sampleConfig.Clock.Separator.Left.FontSize, fallback to %s",
					sampleConfig.Separator.Left.FontSize,
				)

				sampleConfig.Clock.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
			}
		}
	}

	if sampleConfig.Clock.Separator.Right.Color == "" {
		sampleConfig.Clock.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.Clock.Separator.Right.Background == "" {
		sampleConfig.Clock.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.Clock.Separator.Right.Symbol == "" {
		sampleConfig.Clock.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.Clock.Separator.Right.Font == "" {
		sampleConfig.Clock.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.Clock.Separator.Right.FontSize == "" {
		sampleConfig.Clock.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.Clock.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.Clock.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Clock.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.Clock.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Clock.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	// sampleConfig.Clock.LeftClick.Enabled will be false if not set in config

	if len(sampleConfig.Clock.LeftClick.Cmd) == 0 {
		sampleConfig.Clock.LeftClick.Cmd = append(sampleConfig.Clock.LeftClick.Cmd, "true")
	}

	// sampleConfig.Clock.RightClick.Enabled will be false if not set in config

	if len(sampleConfig.Clock.RightClick.Cmd) == 0 {
		sampleConfig.Clock.RightClick.Cmd = append(sampleConfig.Clock.RightClick.Cmd, "true")
	}

	// sampleConfig.Battery.Enabled will be false if not set in config
	if sampleConfig.Battery.Color == "" {
		sampleConfig.Battery.Color = sampleConfig.Color
	}

	if sampleConfig.Battery.Background == "" {
		sampleConfig.Battery.Background = sampleConfig.Background
	}

	if sampleConfig.Battery.Font == "" {
		sampleConfig.Battery.Font = sampleConfig.Font
	}

	if sampleConfig.Battery.FontSize == "" {
		sampleConfig.Battery.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.Battery.Symbol == "" {
		sampleConfig.Battery.Symbol = "⚡"
	}

	if sampleConfig.Battery.SymbolFont == "" {
		sampleConfig.Battery.SymbolFont = sampleConfig.Battery.Font
	}

	if sampleConfig.Battery.SymbolFontSize == "" {
		sampleConfig.Battery.SymbolFontSize = sampleConfig.Battery.FontSize
	}

	// sampleConfig.Battery.ChargeColor.Full will be empty string if not set
	// sampleConfig.Battery.ChargeColor.Empty will be empty string if not set
	// sampleConfig.Battery.ChargeColor.AlmostFull will be empty string if not set
	// sampleConfig.Battery.ChargeColor.AlmostEmpty will be empty string if not set
	if sampleConfig.Battery.Separator.Left.Color == "" {
		sampleConfig.Battery.Separator.Left.Color = sampleConfig.Separator.Left.Color
	}

	if sampleConfig.Battery.Separator.Left.Background == "" {
		sampleConfig.Battery.Separator.Left.Background = sampleConfig.Separator.Left.Background
	}

	if sampleConfig.Battery.Separator.Left.Symbol == "" {
		sampleConfig.Battery.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
	}

	if sampleConfig.Battery.Separator.Left.Font == "" {
		sampleConfig.Battery.Separator.Left.Font = sampleConfig.Separator.Left.Font
	}

	if sampleConfig.Battery.Separator.Left.FontSize == "" {
		sampleConfig.Battery.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.Battery.Separator.Left.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.Battery.Separator.Left.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.Battery.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.Battery.Separator.Left.FontSize, fallback to %s",
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.Battery.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}
	}

	if sampleConfig.Battery.Separator.Right.Color == "" {
		sampleConfig.Battery.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.Battery.Separator.Right.Background == "" {
		sampleConfig.Battery.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.Battery.Separator.Right.Symbol == "" {
		sampleConfig.Battery.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.Battery.Separator.Right.Font == "" {
		sampleConfig.Battery.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.Battery.Separator.Right.FontSize == "" {
		sampleConfig.Battery.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.Battery.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.Battery.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Battery.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.Battery.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Battery.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	// sampleConfig.CpuTemp.Enabled will be false if not set in config

	// No files configured - disable plugin
	if len(sampleConfig.CPUTemp.File) == 0 {
		sampleConfig.CPUTemp.Enabled = false
	}

	if sampleConfig.CPUTemp.Color == "" {
		sampleConfig.CPUTemp.Color = sampleConfig.Color
	}

	if sampleConfig.CPUTemp.Background == "" {
		sampleConfig.CPUTemp.Background = sampleConfig.Background
	}

	if sampleConfig.CPUTemp.Font == "" {
		sampleConfig.CPUTemp.Font = sampleConfig.Font
	}

	if sampleConfig.CPUTemp.FontSize == "" {
		sampleConfig.CPUTemp.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.CPUTemp.Separator.Left.Enabled {
		if sampleConfig.CPUTemp.Separator.Left.Color == "" {
			sampleConfig.CPUTemp.Separator.Left.Color = sampleConfig.Separator.Left.Color
		}

		if sampleConfig.CPUTemp.Separator.Left.Background == "" {
			sampleConfig.CPUTemp.Separator.Left.Background = sampleConfig.Separator.Left.Background
		}

		if sampleConfig.CPUTemp.Separator.Left.Symbol == "" {
			sampleConfig.CPUTemp.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
		}

		if sampleConfig.CPUTemp.Separator.Left.Font == "" {
			sampleConfig.CPUTemp.Separator.Left.Font = sampleConfig.Separator.Left.Font
		}

		if sampleConfig.CPUTemp.Separator.Left.FontSize == "" {
			sampleConfig.CPUTemp.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		} else {
			matched, err := regexp.MatchString(
				`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
				sampleConfig.CPUTemp.Separator.Left.FontSize,
			)

			if err != nil {
				log.Printf(
					"Unable to set sampleConfig.CPUTemp.Separator.Left.FontSize: %s, fallback to %s",
					err,
					sampleConfig.Separator.Left.FontSize,
				)

				sampleConfig.CPUTemp.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
			}

			if !matched {
				log.Printf(
					"Unable to set sampleConfig.CPUTemp.Separator.Left.FontSize, fallback to %s",
					sampleConfig.Separator.Left.FontSize,
				)

				sampleConfig.CPUTemp.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
			}
		}
	}

	if sampleConfig.CPUTemp.Separator.Right.Color == "" {
		sampleConfig.CPUTemp.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.CPUTemp.Separator.Right.Background == "" {
		sampleConfig.CPUTemp.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.CPUTemp.Separator.Right.Symbol == "" {
		sampleConfig.CPUTemp.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.CPUTemp.Separator.Right.Font == "" {
		sampleConfig.CPUTemp.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.CPUTemp.Separator.Right.FontSize == "" {
		sampleConfig.CPUTemp.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.CPUTemp.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf("Unable to set sampleConfig.CPUTemp.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.CPUTemp.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf("Unable to set sampleConfig.CPUTemp.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.CPUTemp.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	// sampleConfig.Vpn.Enabled will false if not set in config

	// No status file - disable plugin
	if sampleConfig.Vpn.StatusFile == "" {
		sampleConfig.Vpn.Enabled = false
	}

	if sampleConfig.Vpn.Color == "" {
		sampleConfig.Vpn.Color = sampleConfig.Color
	}

	if sampleConfig.Vpn.Background == "" {
		sampleConfig.Vpn.Background = sampleConfig.Background
	}

	if sampleConfig.Vpn.Font == "" {
		sampleConfig.Vpn.FontSize = sampleConfig.Font
	}

	if sampleConfig.Vpn.FontSize == "" {
		sampleConfig.Vpn.FontSize = sampleConfig.FontSize
	}

	// Check status file at least once per 3 seconds
	if sampleConfig.Vpn.MtimeThreshold < 3 {
		log.Printf("vpn.mtime_threshold not set, using 3")

		sampleConfig.Vpn.MtimeThreshold = 3
	}

	// sampleConfig.Vpn.DownColor will be empty string if no value set in config
	// sampleConfig.Vpn.UpColor will be empty string if no value set in config
	// sampleConfig.Vpn.TcpCheck.Enabled will false if not set in config

	if sampleConfig.Vpn.Separator.Left.Enabled {
		if sampleConfig.Vpn.Separator.Left.Color == "" {
			sampleConfig.Vpn.Separator.Left.Color = sampleConfig.Separator.Left.Color
		}

		if sampleConfig.Vpn.Separator.Left.Background == "" {
			sampleConfig.Vpn.Separator.Left.Background = sampleConfig.Separator.Left.Background
		}

		if sampleConfig.Vpn.Separator.Left.Symbol == "" {
			sampleConfig.Vpn.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
		}

		if sampleConfig.Vpn.Separator.Left.Font == "" {
			sampleConfig.Vpn.Separator.Left.Font = sampleConfig.Separator.Left.Font
		}

		if sampleConfig.Vpn.Separator.Left.FontSize == "" {
			sampleConfig.Vpn.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		} else {
			matched, err := regexp.MatchString(
				`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
				sampleConfig.Vpn.Separator.Left.FontSize,
			)

			if err != nil {
				log.Printf(
					"Unable to set sampleConfig.Vpn.Separator.Left.FontSize: %s, fallback to %s",
					err,
					sampleConfig.Separator.Left.FontSize,
				)

				sampleConfig.Vpn.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
			}

			if !matched {
				log.Printf(
					"Unable to set sampleConfig.Vpn.Separator.Left.FontSize, fallback to %s",
					sampleConfig.Separator.Left.FontSize,
				)

				sampleConfig.Vpn.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
			}
		}
	}

	if sampleConfig.Vpn.Separator.Right.Color == "" {
		sampleConfig.Vpn.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.Vpn.Separator.Right.Background == "" {
		sampleConfig.Vpn.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.Vpn.Separator.Right.Symbol == "" {
		sampleConfig.Vpn.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.Vpn.Separator.Right.Font == "" {
		sampleConfig.Vpn.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.Vpn.Separator.Right.FontSize == "" {
		sampleConfig.Vpn.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.Vpn.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.Vpn.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Vpn.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.Vpn.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.Vpn.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

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

	// sampleConfig.SimpleVolumePa.Enabled will false if not set in config
	if sampleConfig.SimpleVolumePa.Color == "" {
		sampleConfig.SimpleVolumePa.Color = sampleConfig.Color
	}

	if sampleConfig.SimpleVolumePa.Background == "" {
		sampleConfig.SimpleVolumePa.Background = sampleConfig.Background
	}

	if sampleConfig.SimpleVolumePa.Font == "" {
		sampleConfig.SimpleVolumePa.Font = sampleConfig.Font
	}

	if sampleConfig.SimpleVolumePa.FontSize == "" {
		sampleConfig.SimpleVolumePa.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.SimpleVolumePa.Symbol == "" {
		sampleConfig.SimpleVolumePa.Symbol = `🔊`
	}

	if sampleConfig.SimpleVolumePa.SymbolFont == "" {
		sampleConfig.SimpleVolumePa.SymbolFont = sampleConfig.SimpleVolumePa.Font
	}

	if sampleConfig.SimpleVolumePa.SymbolFontSize == "" {
		sampleConfig.SimpleVolumePa.SymbolFontSize = sampleConfig.SimpleVolumePa.FontSize
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

	if sampleConfig.SimpleVolumePa.WheelDown > 100 {
		sampleConfig.SimpleVolumePa.WheelDown = 5
	}

	if sampleConfig.SimpleVolumePa.MaxVolumeLimit <= 0 {
		sampleConfig.SimpleVolumePa.MaxVolumeLimit = 100
	}

	if sampleConfig.SimpleVolumePa.MaxVolumeLimit >= 120 {
		sampleConfig.SimpleVolumePa.MaxVolumeLimit = 100
	}

	if len(sampleConfig.SimpleVolumePa.RightClickCmd) > 0 {
		if sampleConfig.SimpleVolumePa.RightClickCmd[0] == "" {
			sampleConfig.SimpleVolumePa.RightClickCmd[0] = "true"
		}
	} else {
		sampleConfig.SimpleVolumePa.RightClickCmd = append(sampleConfig.SimpleVolumePa.RightClickCmd, "true")
	}

	if sampleConfig.SimpleVolumePa.Separator.Left.Color == "" {
		sampleConfig.SimpleVolumePa.Separator.Left.Color = sampleConfig.Separator.Left.Color
	}

	if sampleConfig.SimpleVolumePa.Separator.Left.Background == "" {
		sampleConfig.SimpleVolumePa.Separator.Left.Background = sampleConfig.Separator.Left.Background
	}

	if sampleConfig.SimpleVolumePa.Separator.Left.Symbol == "" {
		sampleConfig.SimpleVolumePa.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
	}

	if sampleConfig.SimpleVolumePa.Separator.Left.Font == "" {
		sampleConfig.SimpleVolumePa.Separator.Left.Font = sampleConfig.Separator.Left.Font
	}

	if sampleConfig.SimpleVolumePa.Separator.Left.FontSize == "" {
		sampleConfig.SimpleVolumePa.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.SimpleVolumePa.Separator.Left.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.SimpleVolumePa.Separator.Left.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.SimpleVolumePa.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.SimpleVolumePa.Separator.Left.FontSize, fallback to %s",
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.SimpleVolumePa.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}
	}

	if sampleConfig.SimpleVolumePa.Separator.Right.Color == "" {
		sampleConfig.SimpleVolumePa.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.SimpleVolumePa.Separator.Right.Background == "" {
		sampleConfig.SimpleVolumePa.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.SimpleVolumePa.Separator.Right.Symbol == "" {
		sampleConfig.SimpleVolumePa.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.SimpleVolumePa.Separator.Right.Font == "" {
		sampleConfig.SimpleVolumePa.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.SimpleVolumePa.Separator.Right.FontSize == "" {
		sampleConfig.SimpleVolumePa.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.SimpleVolumePa.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.SimpleVolumePa.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.SimpleVolumePa.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.SimpleVolumePa.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.SimpleVolumePa.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	// sampleConfig.Cron.Enabled will false if not set in config

	if len(sampleConfig.Cron.Tasks) == 0 {
		sampleConfig.Cron.Enabled = false
	}

	if sampleConfig.Cron.TimeZone == "" {
		sampleConfig.Cron.TimeZone = "GMT+0"
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

	if sampleConfig.NetIf.Color == "" {
		sampleConfig.NetIf.Color = sampleConfig.Color
	}

	if sampleConfig.NetIf.Background == "" {
		sampleConfig.NetIf.Background = sampleConfig.Background
	}

	if sampleConfig.NetIf.Font == "" {
		sampleConfig.NetIf.Font = sampleConfig.Font
	}

	if sampleConfig.NetIf.FontSize == "" {
		sampleConfig.NetIf.FontSize = sampleConfig.FontSize
	}

	if sampleConfig.NetIf.Separator.Left.Color == "" {
		sampleConfig.NetIf.Separator.Left.Color = sampleConfig.Separator.Left.Color
	}

	if sampleConfig.NetIf.Separator.Left.Background == "" {
		sampleConfig.NetIf.Separator.Left.Background = sampleConfig.Separator.Left.Background
	}

	if sampleConfig.NetIf.Separator.Left.Symbol == "" {
		sampleConfig.NetIf.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
	}

	if sampleConfig.NetIf.Separator.Left.Font == "" {
		sampleConfig.NetIf.Separator.Left.Font = sampleConfig.Separator.Left.Font
	}

	if sampleConfig.NetIf.Separator.Left.FontSize == "" {
		sampleConfig.NetIf.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.NetIf.Separator.Left.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.NetIf.Separator.Left.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.NetIf.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}

		if !matched {
			log.Printf("Unable to set sampleConfig.NetIf.Separator.Left.FontSize, fallback to %s",
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.NetIf.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}
	}

	if sampleConfig.NetIf.Separator.Right.Color == "" {
		sampleConfig.NetIf.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.NetIf.Separator.Right.Background == "" {
		sampleConfig.NetIf.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.NetIf.Separator.Right.Symbol == "" {
		sampleConfig.NetIf.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.NetIf.Separator.Right.Font == "" {
		sampleConfig.NetIf.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.NetIf.Separator.Right.FontSize == "" {
		sampleConfig.NetIf.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.NetIf.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.NetIf.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.NetIf.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.NetIf.Separator.Right.FontSize, fallback to %s",
				sampleConfig.NetIf.Separator.Right.FontSize,
			)

			sampleConfig.NetIf.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	// sampleConfig.AppButtons.Enabled will false if not set in config

	if sampleConfig.AppButtons.Separator.Left.Color == "" {
		sampleConfig.AppButtons.Separator.Left.Color = sampleConfig.Separator.Left.Color
	}

	if sampleConfig.AppButtons.Separator.Left.Background == "" {
		sampleConfig.AppButtons.Separator.Left.Background = sampleConfig.Separator.Left.Background
	}

	if sampleConfig.AppButtons.Separator.Left.Symbol == "" {
		sampleConfig.AppButtons.Separator.Left.Symbol = sampleConfig.Separator.Left.Symbol
	}

	if sampleConfig.AppButtons.Separator.Left.Font == "" {
		sampleConfig.AppButtons.Separator.Left.Font = sampleConfig.Separator.Left.Font
	}

	if sampleConfig.AppButtons.Separator.Left.FontSize == "" {
		sampleConfig.AppButtons.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.AppButtons.Separator.Left.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.AppButtons.Separator.Left.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.AppButtons.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}

		if !matched {
			log.Printf(
				"Unable to set sampleConfig.AppButtons.Separator.Left.FontSize, fallback to %s",
				sampleConfig.Separator.Left.FontSize,
			)

			sampleConfig.AppButtons.Separator.Left.FontSize = sampleConfig.Separator.Left.FontSize
		}
	}

	if sampleConfig.AppButtons.Separator.Right.Color == "" {
		sampleConfig.AppButtons.Separator.Right.Color = sampleConfig.Separator.Right.Color
	}

	if sampleConfig.AppButtons.Separator.Right.Background == "" {
		sampleConfig.AppButtons.Separator.Right.Background = sampleConfig.Separator.Right.Background
	}

	if sampleConfig.AppButtons.Separator.Right.Symbol == "" {
		sampleConfig.AppButtons.Separator.Right.Symbol = sampleConfig.Separator.Right.Symbol
	}

	if sampleConfig.AppButtons.Separator.Right.Font == "" {
		sampleConfig.AppButtons.Separator.Right.Font = sampleConfig.Separator.Right.Font
	}

	if sampleConfig.AppButtons.Separator.Right.FontSize == "" {
		sampleConfig.AppButtons.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
	} else {
		matched, err := regexp.MatchString(
			`^(xx-small|x-small|small|medium|large|x-large|xx-large|smaller|larger)$`,
			sampleConfig.AppButtons.Separator.Right.FontSize,
		)

		if err != nil {
			log.Printf(
				"Unable to set sampleConfig.AppButtons.Separator.Right.FontSize: %s, fallback to %s",
				err,
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.AppButtons.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}

		if !matched {
			log.Printf("Unable to set sampleConfig.AppButtons.Separator.Right.FontSize, fallback to %s",
				sampleConfig.Separator.Right.FontSize,
			)

			sampleConfig.AppButtons.Separator.Right.FontSize = sampleConfig.Separator.Right.FontSize
		}
	}

	if len(sampleConfig.Apps) == 0 {
		sampleConfig.AppButtons.Enabled = false
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
