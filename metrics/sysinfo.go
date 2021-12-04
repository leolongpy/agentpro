package metrics

import (
	"agentpro/models"
	"agentpro/settings"
)

func SysMetrics() (L []*models.MetricValue) {
	L = append(L, models.GaugeValue("sysinfo.innerip", settings.IP()))
	return
}
