package cron

import (
	"agentpro/logger"
	"agentpro/utils"
	"github.com/golang/protobuf/ptypes/timestamp"
	Cron "github.com/robfig/cron/v3"
	"time"
)

type ScheduledTask struct {
	taskId string
	cmd    string
	exec   string
}

func (t ScheduledTask) Run() {
	dt := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	out, exitCode, err := utils.CmdOutExitCodeBytes(t.exec, t.cmd)
	if err != nil {
		logger.StartupDebug("exec run error:", err)
	}
	cmdData := string(out)
	logger.StartupInfo("结果:", cmdData)
	enddt := &timestamp.Timestamp{Seconds: time.Now().Unix()}
	logger.StartupInfo("脚本执行信息:", cmdData, t.taskId, int32(exitCode), dt, enddt)
}

func CronStart(cmd, task_id, exec, express string) Cron.EntryID {
	// AddJob方法
	logger.StartupDebug("正式启动", cmd, task_id, exec, express)

	id, _ := XCron.AddJob(express, &ScheduledTask{
		task_id,
		cmd,
		exec,
	})
	return id
}
