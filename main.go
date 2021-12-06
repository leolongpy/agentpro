package main

import (
	"agentpro/cron"
	"agentpro/heartbeat"
	"agentpro/http"
	"agentpro/logger"
	"agentpro/metrics"
	"agentpro/settings"
	"fmt"
	goHttp "net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	logger.StartupDebug("begin")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGALRM)
	go func() {
		signalType := <-ch
		signal.Stop(ch)
		fmt.Println("退出...")
		fmt.Println("收到OS信号类型 : ", signalType)
		logger.StartupDebug("收到OS信号类型 : ", signalType)
		// 删除pid
		os.Remove(settings.Config().Pid)
		os.Exit(0)
	}()
	if len(os.Args) != 2 {
		fmt.Printf("使用说明 : %s [start|stop|version] \n ", os.Args[0])
		os.Exit(0) // 安全退出
	}
	settings.LoadConfiguration()
	settings.InitLocalIp()
	metrics.BuildMappers()

	cron.GetAllTask()
	cron.InitWatchTask()
	heartbeat.AgentHealthChecks()
	go cron.InitDataHistory()
	cron.Collect()
	if strings.ToLower(os.Args[1]) == "main" {
		go func() {
			goHttp.ListenAndServe("0.0.0.0:16060", nil)
		}()
		http.Start()
	}
	settings.HandleControl(os.Args[1])
}
