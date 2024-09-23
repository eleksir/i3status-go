package lib

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"
)

// UpdateCPUTemperature gets and updates average CPU cores temperature.
func (c *MyConfig) UpdateCPUTemperature() {
	n := len(c.CPUTemp.File)

	for {
		var (
			temperature = make([]int64, n)
			tSum        int64
			tAvg        int64
		)

		for i, filename := range c.CPUTemp.File {
			file, err := os.Open(filename)

			if err != nil {
				log.Printf("Unable to open %s: %s", filename, err)
			} else {
				reader := bufio.NewReader(file)
				line, _, err := reader.ReadLine()

				if err != nil {
					log.Printf("Unable to read from %s: %s", filename, err)
					err = file.Close()

					if err != nil {
						log.Printf("Unable to close %s: %s", filename, err)
					}
				} else {
					err = file.Close()

					if err != nil {
						log.Printf("Unable to close %s: %s", filename, err)
					} else {
						temp, err := strconv.ParseInt(string(line), 10, 32)

						if err != nil {
							log.Printf("Unable to convert string to number from file %s: %s", filename, err)
						} else {
							if temp > 1000 {
								temp /= 1000
							}

							temperature[i] = temp
						}
					}
				}
			}
		}

		if len(temperature) == 1 {
			tAvg = temperature[0]
		} else {
			for _, t := range temperature {
				tSum += t
			}

			tAvg = tSum / int64(len(temperature))
		}

		if c.Values.CPUTemperature != tAvg {
			c.Values.CPUTemperature = tAvg
			c.Channels.UpdateReady <- true
		}

		time.Sleep(3 * time.Second)
	}
}
