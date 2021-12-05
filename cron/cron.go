package cron

import (
	"agentpro/logger"
	"agentpro/settings"
	"github.com/pkg/errors"
	Cron "github.com/robfig/cron/v3"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var (
	XCron *Cron.Cron = Cron.New(Cron.WithParser(Cron.NewParser(
		Cron.Minute | Cron.Hour | Cron.Dom | Cron.Month | Cron.Dow)))
	client   *clientv3.Client
	initEtcd sync.Once
)

func init() {
	XCron.Start()
}

func initScheduler() error {
	if !settings.Config().Registrar.Enable {
		err := errors.New("已经被初始化")
		logger.StartupDebug("cfg配置文件的Registrar选项被关闭")
		return err
	} else {
		addrs := settings.Config().Registrar.Addrs
		err := errors.New("已经被初始化")
		initEtcd.Do(func() {
			err = nil
			config := clientv3.Config{
				Endpoints:   addrs,
				DialTimeout: 5 * time.Second,
			}

			if client, err = clientv3.New(config); err != nil {
				logger.StartupFatal("与etcd服务器建立连接失败", err)
			} else {
				logger.StartupInfo("与etcd建立连接")
			}
		})
		return err
	}
}
