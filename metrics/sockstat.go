package metrics

import (
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"fmt"
	"log"
)

func SocketStatSummaryMetrics() (L []*models.MetricValue) {
	ssMap, err := nux.SocketStatSummary()
	if err != nil {
		log.Println(err)
		return
	}

	tags := fmt.Sprintf("__IP=%s", settings.IP())
	for k, v := range ssMap {
		L = append(L, models.GaugeValue("ss."+k, v, tags))
	}

	return
}
