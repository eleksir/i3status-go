package lib

import (
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// RunChan to this channel we send command and its argv to execute as separate process.
var RunChan = make(chan []string, 128)

// Spawner forks and execs given program, and also detaches form its control tty.
func (c MyConfig) Spawner() {
	for prg := range RunChan {
		devnullR, _ := os.Open(os.DevNull)
		devnullW, _ := os.OpenFile(os.DevNull, os.O_WRONLY|os.O_APPEND, 0644)

		spa := syscall.SysProcAttr{
			Setsid:     true,
			Foreground: false,
		}

		pa := syscall.ProcAttr{
			Dir:   os.Getenv(`HOME`),
			Env:   os.Environ(),
			Files: []uintptr{devnullR.Fd(), devnullW.Fd(), devnullW.Fd()},
			Sys:   &spa,
		}

		bin, err := exec.LookPath(prg[0])

		if err != nil {
			log.Printf("Unable to spawn %s: %s", prg[0], err)
		}

		_, err = syscall.ForkExec(bin, prg, &pa)

		if err != nil {
			log.Printf("Unable to spawn %s: %s", prg[0], err)
		}

		// Close file errors we don't handle :) at all.
		_ = devnullR.Close()
		_ = devnullW.Close()
	}
}

// CleanZombies reaps processes spawned by Spawner() and already exited.
func (c MyConfig) CleanZombies() {
	r := syscall.Rusage{}

	for {
		_, _ = syscall.Wait4(-1, nil, 0, &r)

		time.Sleep(1 * time.Minute)
	}
}
