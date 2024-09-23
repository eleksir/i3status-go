package lib

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	p "github.com/mafik/pulseaudio"
)

// TODO: https://twin.sh/articles/44/add-a-timeout-to-any-function-in-go timeout pulseaudio calls

// UpdateVolumeInfo updates info about current Sound Volume.
func (c *MyConfig) UpdateVolumeInfo() {
	var err error

	c.Values.PA, err = p.NewClient()

	// It can happen if no pulseaudio server running for current user.
	// If no server running we have to run one.
	if err != nil {
		if err := c.PaReinit(); err != nil {
			log.Print(err)

			return
		}
	}

	vol, err := c.Values.PA.Volume()

	if err != nil {
		log.Printf("Unable get volume from pulseaudio server: %s", err)

		return
	}

	c.Values.SoundVolume = fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
		c.SimpleVolumePa.Color,
		c.SimpleVolumePa.Background,
		c.SimpleVolumePa.SymbolFont,
		c.SimpleVolumePa.SymbolFontSize,
		c.SimpleVolumePa.Symbol,
	)

	c.Values.SoundVolume += fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>:%d%%</span>",
		c.SimpleVolumePa.Color,
		c.SimpleVolumePa.Background,
		c.SimpleVolumePa.Font,
		c.SimpleVolumePa.FontSize,
		int64(vol*100),
	)

	c.Channels.UpdateReady <- true

	for {
		// Subscribe to update notification channel, to get info that volume changed.
		pulseUpdate, err := c.Values.PA.Updates()

		if err != nil {
			log.Printf("Unable to subscribe to pulseaudio updates: %s", err)

			return
		}

		// Rake update events.
		for range pulseUpdate {
			vol, err = c.Values.PA.Volume()

			if err != nil {
				log.Printf("Unable get volume from pulseaudio server: %s", err)

				return
			}

			c.Values.SoundVolume = fmt.Sprintf(
				"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
				c.SimpleVolumePa.Color,
				c.SimpleVolumePa.Background,
				c.SimpleVolumePa.SymbolFont,
				c.SimpleVolumePa.SymbolFontSize,
				c.SimpleVolumePa.Symbol,
			)

			c.Values.SoundVolume += fmt.Sprintf(
				"<span color='%s' background='%s' font='%s' size='%s'>:%d%%</span>",
				c.SimpleVolumePa.Color,
				c.SimpleVolumePa.Background,
				c.SimpleVolumePa.Font,
				c.SimpleVolumePa.FontSize,
				int64(vol*100),
			)

			c.Channels.UpdateReady <- true
		}

		if err := c.PaReinit(); err != nil {
			log.Print(err)

			return
		}
	}

	// This code is unreachable :(
	c.Values.PA.Close() //nolint:govet
}

// PaReinit re-inits pulseaudio and connection to it.
func (c *MyConfig) PaReinit() error {
	var err error

	// If we out of updates, seems someone killed pulseaudio server. Restart it.
	cmd := exec.Command("pulseaudio", "--check")

	if err := cmd.Run(); err == nil {
		log.Printf("pulseaudio already running, but seems not responding, try gracefully kill it")

		cmd = exec.Command("pulseaudio", "--kill")

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("unable to kill pulseaudio: %w", err)
		}

		cnt := 0

		for {
			cmd := exec.Command("pulseaudio", "--check")

			// PA not running. Breakout of loop.
			if err := cmd.Run(); err != nil {
				break
			}

			if cnt >= 5 {
				return errors.New("timeout waiting pulseaudio to exit") //nolint: goerr113
			}

			log.Print("Waiting for pulseaudio to exit")

			cnt++

			time.Sleep(300 * time.Millisecond)
		}
	}

	// PA has annoying feature: socket activation. Sound system must be persistent and available to all users regardless
	// of init or other things. They just fail to understand that. Or they fucked up pulseaudio just because they scary
	// of security something responsibility. Anyway this is grand design flaw, but we still have to bear with it.
	// This setting defines if pulseaudio will be run in manner that allows it to exit is other login detected (how
	// they guess it - I do not know). In general case we do not want to allow pa to exit and leave our session without
	// audio.
	if c.SimpleVolumePa.DontExitOnLogin {
		// Do not exit on any kind of login/logout events.
		cmd = exec.Command("pulseaudio", "--exit-idle-time=-1", "--start")
	} else {
		// Exit immediately on logout(?).
		cmd = exec.Command("pulseaudio", "--exit-idle-time=0", "--start")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to initialize pulseaudio server instance: %w", err)
	}

	c.Values.PA, err = p.NewClient()

	if err != nil {
		return fmt.Errorf("unable to make client connection to pulseaudio: %w", err)
	}

	return err
}

func (c *MyConfig) SVPAHandler() {
	for e := range c.Channels.SVPAHandlerChan {
		if e.Button == 3 {
			c.Channels.RunChan <- c.SimpleVolumePa.RightClickCmd

			continue
		}

		vol, err := c.Values.PA.Volume()

		if err != nil {
			if err := c.PaReinit(); err != nil {
				log.Printf("Unable to get pulseaudio volume: %s", err)
			} else {
				vol, err = c.Values.PA.Volume()

				if err != nil {
					log.Printf("Unable to get volume pulseaudio server behaves weirdly: %s", err)
				}
			}
		}

		switch e.Button {
		case c.SimpleVolumePa.WheelUp:
			vol += float32(c.SimpleVolumePa.Step) / 100

			if vol > (float32(c.SimpleVolumePa.MaxVolumeLimit) / 100) {
				vol = float32(c.SimpleVolumePa.MaxVolumeLimit) / 100
			}

			if err := c.Values.PA.SetVolume(vol); err != nil {
				log.Printf("Unable to set pulseaudio volume: %s", err)
			}

		case c.SimpleVolumePa.WheelDown:
			vol -= float32(c.SimpleVolumePa.Step) / 100

			if vol < 0 {
				vol = 0
			}

			if err := c.Values.PA.SetVolume(vol); err != nil {
				log.Printf("Unable to set pulseaudio volume: %s", err)
			}
		}
	}
}
