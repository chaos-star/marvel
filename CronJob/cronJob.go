package CronJob

import "github.com/robfig/cron/v3"

type CronJob struct {
	jobs map[string]SpecJob
	*cron.Cron
}

type SpecJob struct {
	Spec string
	Cmd  cron.Job
}

func (cj *CronJob) JoinJobs(jobs map[string]SpecJob) {
	cj.jobs = jobs
	for _, item := range jobs {
		cj.AddJob(item.Spec, item.Cmd)
	}
}
func (cj *CronJob) GetJob(name string) *SpecJob {
	var pJob *SpecJob
	if job, ok := cj.jobs[name]; ok {
		pJob = &job
	}
	return pJob
}

func Initialize() *CronJob {
	return &CronJob{
		make(map[string]SpecJob),
		cron.New(),
	}
}
