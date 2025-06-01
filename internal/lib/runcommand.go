package lib

import (
	"log"
	"time"
)

// RunCommands gets list of commands runs them sequentally and conacatenates their output.
func (c *MyConfig) RunCommand() {
	var (
		InitialDelay       = 100 * time.Millisecond
		LoopIterationDelay = time.Duration(c.CmdRun.Delay) * time.Second
		Delay              = InitialDelay
		ticker             = time.NewTicker(Delay)
	)

	for range ticker.C {
		if Delay == InitialDelay {
			Delay = LoopIterationDelay
			ticker.Reset(Delay)
		}

		var (
			outputString string
			command      = []string{c.CmdRun.Cmd}
		)

		if len(c.CmdRun.Args) > 0 {
			command = append(command, c.CmdRun.Args...)
		}

		outputString += RunProcess(command)

		log.Println(c.Values.RunCommandOutput)

		if c.Values.RunCommandOutput != outputString {
			c.Values.RunCommandOutput = outputString
			c.Channels.UpdateReady <- true
		}
	}
}
