package Job

import "github.com/robfig/cron/v3"

type Job struct {
	*cron.Cron
}

func Initialize() *Job {
	return &Job{
		cron.New(),
	}
}
