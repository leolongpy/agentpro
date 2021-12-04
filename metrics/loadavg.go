package metrics

import (
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"fmt"
	"log"
)

func LoadAvgMetrics() []*models.MetricValue {
	load, err := nux.LoadAvg()
	if err != nil {
		log.Println(err)
		return nil
	}
	tags := fmt.Sprintf("__IP=%s", settings.IP())
	return []*models.MetricValue{
		models.GaugeValue("load.1min", load.Avg1min, tags),
		models.GaugeValue("load.5min", load.Avg5min, tags),
		models.GaugeValue("load.15min", load.Avg15min, tags),
	}

}
