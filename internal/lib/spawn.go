package lib

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Spawner forks and execs given program, and also detaches form its control tty.
func (c *MyConfig) Spawner() {
	for prg := range c.Channels.RunChan {
		<-time.After(25 * time.Millisecond)

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
func (c *MyConfig) CleanZombies() {
	r := syscall.Rusage{}

	for {
		_, _ = syscall.Wait4(-1, nil, 0, &r)

		time.Sleep(1 * time.Minute)
	}
}

// RunProcess runs given command and returns contens of command stdout. On error prints error on stderr and returns
// empty syting.
func RunProcess(command []string) string {
	var (
		cmd            *exec.Cmd
		stdout, stderr bytes.Buffer
	)

	switch len(command) {
	case 1:
		cmd = exec.Command(command[0], "") //nolint: gosec
	default:
		cmd = exec.Command(command[0], command[1:]...) //nolint: gosec
	}

	cmd.Dir = "/"
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if stdout.String() != "" {
			log.Printf("Unable to run '%s': %s, %s", strings.Join(command, " "), err, stdout.String())
		} else {
			log.Printf("Unable to run '%s': %s", strings.Join(command, " "), err)
		}

		return ""
	}

	if stderr.String() != "" {
		log.Printf("%s: %s", strings.Join(command, " "), stderr.String())
	}

	return strings.TrimRight(stdout.String(), "\n")
}
