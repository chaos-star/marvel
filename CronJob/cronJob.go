package CronJob

import "github.com/robfig/cron/v3"

type CronJob struct {
	*cron.Cron
}

func Initialize() *CronJob {
	return &CronJob{
		cron.New(),
	}
}
