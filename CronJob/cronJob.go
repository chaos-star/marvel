package CronJob

import "github.com/robfig/cron/v3"

type CronJob struct {
	*cron.Cron
}
type Job cron.Job

func Initialize() *CronJob {
	return &CronJob{
		cron.New(),
	}
}
