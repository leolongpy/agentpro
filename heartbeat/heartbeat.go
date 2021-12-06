package heartbeat

import (
	"agentpro/cron"
	"agentpro/settings"
	"time"
)

func AgentHealthChecks() {
	if settings.Config().Heartbeat.Enabled {
		go cron.RegisterPutHealth(time.Duration(settings.Config().Heartbeat.Interval)*time.Second, settings.IP())
	}
}
