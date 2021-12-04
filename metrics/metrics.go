package metrics

import (
	model "agentpro/models"
	"agentpro/settings"
)

type FuncsAndInterval struct {
	Fs       []func() []*model.MetricValue
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := settings.Config().Transfer.Interval
	Mappers = []FuncsAndInterval{
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				CpuMetrics,     // 监控指标-CPU
				NetMetrics,     // 监控指标-网络
				KernelMetrics,  // 监控指标-负载
				LoadAvgMetrics, // 监控指标-句柄数
				MemMetrics,     // 内存
				DiskIOMetrics,
				IOStatsMetrics, // 监控指标-IO
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				SocketStatSummaryMetrics, //监控指标-TCP连接
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				DeviceMetrics,
			},
			Interval: interval,
		},
		FuncsAndInterval{
			Fs: []func() []*model.MetricValue{
				SysMetrics,
			},
			Interval: interval,
		},
	}
}
