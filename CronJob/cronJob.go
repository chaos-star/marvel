package CronJob

import "github.com/robfig/cron/v3"

type CronJob struct {
	*cron.Cron
}

type SpecJob struct {
	Spec string
	Cmd  cron.Job
}

func (cj *CronJob) JoinJobs(jobs map[string]SpecJob) {
	for _, item := range jobs {
		cj.AddJob(item.Spec, item.Cmd)
	}
}

func Initialize() *CronJob {
	return &CronJob{
		cron.New(),
	}
}
