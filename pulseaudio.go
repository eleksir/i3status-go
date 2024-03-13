package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"

	p "github.com/mafik/pulseaudio"
)

var SoundVolume = "🔊:0%"
var pa *p.Client

var SVPAHandlerChan = make(chan ClickEvent, 256)

// TODO: https://twin.sh/articles/44/add-a-timeout-to-any-function-in-go timeout pulseaudio calls

// UpdateVolumeInfo updates info about current Sound Volume.
func UpdateVolumeInfo() {
	var err error

	pa, err = p.NewClient()

	// It can happen if no pulseaudio server running for current user.
	// If no server running we have to run one.
	if err != nil {
		if err := PaReinit(); err != nil {
			log.Print(err)

			return
		}
	}

	vol, err := pa.Volume()

	if err != nil {
		log.Printf("Unable get volume from pulseaudio server: %s", err)

		return
	}

	SoundVolume = fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
		Conf.SimpleVolumePa.Color,
		Conf.SimpleVolumePa.Background,
		Conf.SimpleVolumePa.SymbolFont,
		Conf.SimpleVolumePa.SymbolFontSize,
		Conf.SimpleVolumePa.Symbol,
	)

	SoundVolume += fmt.Sprintf(
		"<span color='%s' background='%s' font='%s' size='%s'>:%d%%</span>",
		Conf.SimpleVolumePa.Color,
		Conf.SimpleVolumePa.Background,
		Conf.SimpleVolumePa.Font,
		Conf.SimpleVolumePa.FontSize,
		int64(vol*100),
	)

	UpdateReady <- true

	for {
		// Subscribe to update notification channel, to get info that volume changed.
		pulseUpdate, err := pa.Updates()

		if err != nil {
			log.Printf("Unable to subscribe to pulseaudio updates: %s", err)

			return
		}

		// Rake update events.
		for range pulseUpdate {
			vol, err = pa.Volume()

			if err != nil {
				log.Printf("Unable get volume from pulseaudio server: %s", err)

				return
			}

			SoundVolume = fmt.Sprintf(
				"<span color='%s' background='%s' font='%s' size='%s'>%s</span>",
				Conf.SimpleVolumePa.Color,
				Conf.SimpleVolumePa.Background,
				Conf.SimpleVolumePa.SymbolFont,
				Conf.SimpleVolumePa.SymbolFontSize,
				Conf.SimpleVolumePa.Symbol,
			)

			SoundVolume += fmt.Sprintf(
				"<span color='%s' background='%s' font='%s' size='%s'>:%d%%</span>",
				Conf.SimpleVolumePa.Color,
				Conf.SimpleVolumePa.Background,
				Conf.SimpleVolumePa.Font,
				Conf.SimpleVolumePa.FontSize,
				int64(vol*100),
			)

			UpdateReady <- true
		}

		if err := PaReinit(); err != nil {
			log.Print(err)

			return
		}
	}

	// This code is unreachable :(
	pa.Close() //nolint:govet
}

// PaReinit re-inits pulseaudio and connection to it.
func PaReinit() error {
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
	if Conf.SimpleVolumePa.DontExitOnLogin {
		// Do not exit on any kind of login/logout events.
		cmd = exec.Command("pulseaudio", "--exit-idle-time=-1", "--start")
	} else {
		// Exit immediately on logout(?).
		cmd = exec.Command("pulseaudio", "--exit-idle-time=0", "--start")
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to initialize pulseaudio server instance: %w", err)
	}

	pa, err = p.NewClient()

	if err != nil {
		return fmt.Errorf("unable to make client connection to pulseaudio: %w", err)
	}

	return err
}

func SVPAHandler() {
	for e := range SVPAHandlerChan {
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
}
