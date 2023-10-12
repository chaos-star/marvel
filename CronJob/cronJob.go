package CronJob

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strings"
)

type CronJob struct {
	jobs map[string]SpecJob
	*cron.Cron
}

type SpecJob struct {
	Spec string
	Cmd  cron.Job
	Desc string
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

func (cj *CronJob) GetJobs() []*SpecJob {
	var pJobs []*SpecJob
	for index, _ := range cj.jobs {
		job := cj.jobs[index]
		pJobs = append(pJobs, &job)
	}
	return pJobs
}

func (cj *CronJob) Usage(tips ...string) {
	var (
		i   int
		msg strings.Builder
	)
	if len(tips) > 0 {
		for _, tip := range tips {
			msg.WriteString(fmt.Sprintf("%s\r\n", tip))
		}
		msg.WriteString("\r\n")
	}
	msg.WriteString(fmt.Sprintf("\033[1m%-5s %-20s %-30s\033[0m\r\n", "index", "command", "description"))
	for index, item := range cj.jobs {
		i++
		msg.WriteString(fmt.Sprintf("\033[36m%-5s\033[0m \033[32m%-20s\033[0m %-30s\r\n", fmt.Sprintf("%d.", i), index, item.Desc))
	}
	fmt.Println(msg.String())
}

func Initialize() *CronJob {
	return &CronJob{
		make(map[string]SpecJob),
		cron.New(),
	}
}
