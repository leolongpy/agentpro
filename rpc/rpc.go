package rpc

import (
	"agentpro/logger"
	"agentpro/models"
	"agentpro/settings"
	"context"
	"math/rand"
	"time"
)

var requestCtx context.Context

func updateMetrics(addr string, metrics []*models.MetricValue, resp *models.TransferResponse) bool {
	logger.JobInfo(111, metrics, resp) //在大规模写入时可能会造成内存溢出
	return true
}

func SendMetrics(metrics []*models.MetricValue, resp *models.TransferResponse) {
	rand.Seed(time.Now().UnixNano())
	for _, i := range rand.Perm(len(settings.Config().Transfer.Addrs)) {
		addr := settings.Config().Transfer.Addrs[i]
		if updateMetrics(addr, metrics, resp) {
			break
		}
	}
}
