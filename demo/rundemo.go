package main

import (
	"fmt"
	LFCron "github.com/robfig/cron/v3"
	"os/exec"
	"time"
)

var XXCron *LFCron.Cron = LFCron.New(LFCron.WithParser(
	LFCron.NewParser(LFCron.Minute | LFCron.Hour | LFCron.Dom | LFCron.Month | LFCron.Dow),
))

func init() {
	XXCron.Start()
}

type TestJob struct {
	task_id string
	cmd     string
}

func (t TestJob) Run() {
	dt := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("运行中...", dt)
	c1, err := exec.Command("bash", "-c", t.cmd).Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(c1))
}

func CronDemo() LFCron.EntryID {
	spec := "* * * * *"
	id, err := XXCron.AddJob(spec, &TestJob{
		"1",
		"df -h",
	})
	fmt.Println(err)
	return id
}

func main() {
	CronDemo()
	select {}
}
