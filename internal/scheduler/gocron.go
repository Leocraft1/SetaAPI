package scheduler

import (
	"fmt"
	"setaapi/internal/data"
	"setaapi/internal/repository"

	"github.com/go-co-op/gocron/v2"
)

func InitScheduler() (gocron.Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	_, err = s.NewJob(gocron.CronJob("*/20 * * * * *", true), gocron.NewTask(updateStopsTask), gocron.WithSingletonMode(gocron.LimitModeReschedule))
	_, err = s.NewJob(gocron.CronJob("*/20 * * * * *", true), gocron.NewTask(updateRoutesTask), gocron.WithSingletonMode(gocron.LimitModeReschedule))
	_, err = s.NewJob(gocron.CronJob("0 0 * * *", false), gocron.NewTask(updateRoutesStatusTask), gocron.WithSingletonMode(gocron.LimitModeReschedule))

	if err != nil {
		return nil, err
	}

	s.Start()
	return s, nil
}

func updateStopsTask() {
	fmt.Println("Task UpdateStops")
	data.UpdateStops()
}

func updateRoutesTask() {
	fmt.Println("Task UpdateRoutes")
	data.UpdateRoutes()
}

func updateRoutesStatusTask() {
	fmt.Println("Task UpdateRoutesStatus")
	repository.UpdateRoutesStatus()
}