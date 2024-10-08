package lib

import (
	"log"
	"os"
	"syscall"
)

var PrintOutput = true

// SigHandler OS signal handler.
func (c *MyConfig) SigHandler() {
	for s := range c.Channels.SigChan {
		switch s {
		case syscall.SIGUSR1:
			log.Print("Got SIGUSR1, stopping output")

			c.Values.PrintOutput = false

		case syscall.SIGUSR2:
			log.Print("Got SIGUSR2, resuming output")

			c.Values.PrintOutput = true

			c.Channels.UpdateReady <- true

		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			os.Exit(0)

		// We have signal that we're not interested in, so make a new loop iteration.
		default:
			continue
		}
	}
}
