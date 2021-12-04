package metrics

import (
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"fmt"
	"log"
)

func MemMetrics() []*models.MetricValue {
	m, err := nux.MemInfo()
	if err != nil {
		log.Println(err)
		return nil
	}

	memFree := m.MemFree + m.Buffers + m.Cached + m.Sreclaimable
	memUsed := m.MemTotal - memFree

	pmemFree := 0.0
	pmemUsed := 0.0
	if m.MemTotal != 0 {
		pmemFree = float64(memFree) * 100.0 / float64(m.MemTotal)
		pmemUsed = float64(memUsed) * 100.0 / float64(m.MemTotal)
	}

	pswapUsed := 0.0
	if m.SwapTotal != 0 {
		pswapUsed = float64(m.SwapUsed) * 100.0 / float64(m.SwapTotal)
	}
	tags := fmt.Sprintf("__IP=%s", settings.IP())
	return []*models.MetricValue{
		models.GaugeValue("mem.memused", memUsed, tags),
		models.GaugeValue("mem.memfree", memFree, tags),
		models.GaugeValue("mem.memfree.percent", pmemFree, tags),
		models.GaugeValue("mem.memused.percent", pmemUsed, tags),
		models.GaugeValue("mem.swapused.percent", pswapUsed, tags),
	}

}
