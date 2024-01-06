package main

import (
	"log"
	"os"
	"syscall"
)

var PrintOutput = true
var sigChan = make(chan os.Signal, 1)

// sigHandler OS signal handler.
func sigHandler() {
	for s := range sigChan {
		switch s {
		case syscall.SIGUSR1:
			log.Print("Got SIGUSR1, stopping output")

			PrintOutput = false

		case syscall.SIGUSR2:
			log.Print("Got SIGUSR2, resuming output")

			PrintOutput = true

			UpdateReady <- true

		// We have signal that we're not interested in, so make a new loop iteration.
		default:
			continue
		}
	}
}
