package main

import (
	"log"

	"github.com/go-co-op/gocron/v2"
)

// myCron contains cron handler.
var myCron gocron.Scheduler

// RunCron parses config part, related to cron tasks and kicks cron to run.
func RunCron() {
	var err error

	myCron, err = gocron.NewScheduler()

	if err != nil {
		log.Printf("Unable to create built-in cron: %s", err)

		return
	}

	for _, job := range Conf.Cron.Tasks {
		_, err := myCron.NewJob(
			gocron.CronJob(job.Time, false),
			gocron.NewTask(
				// Сука, просто заспавнить бинарь, ну почему это надо делать через жопу-то?
				any(
					func() bool {
						RunChan <- job.Cmd

						return true
					},
				),
			),
		)

		if err != nil {
			log.Printf("Unable to add cron job: %s", err)
		}
	}

	myCron.Start()
}
