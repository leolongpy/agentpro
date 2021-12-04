package metrics

import (
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"fmt"
	"log"
)

func NetMetrics() []*models.MetricValue {
	return CoreNetMetrics(settings.Config().Collector.IfacePrefix)
}

func CoreNetMetrics(ifacePrefix []string) []*models.MetricValue {

	netIfs, err := nux.NetIfs(ifacePrefix)
	if err != nil {
		log.Println(err)
		return []*models.MetricValue{}
	}

	cnt := len(netIfs)
	ret := make([]*models.MetricValue, cnt*23)

	for idx, netIf := range netIfs {
		iface := fmt.Sprintf("iface=%s,__IP=%s", netIf.Iface, settings.IP())
		ret[idx*23+0] = models.CounterValue("net.if.in.bytes", netIf.InBytes, iface)
		ret[idx*23+1] = models.CounterValue("net.if.in.packets", netIf.InPackages, iface)
		ret[idx*23+2] = models.CounterValue("net.if.in.errors", netIf.InErrors, iface)
		ret[idx*23+3] = models.CounterValue("net.if.in.dropped", netIf.InDropped, iface)
		ret[idx*23+4] = models.CounterValue("net.if.in.fifo.errs", netIf.InFifoErrs, iface)
		ret[idx*23+5] = models.CounterValue("net.if.in.frame.errs", netIf.InFrameErrs, iface)
		ret[idx*23+6] = models.CounterValue("net.if.in.compressed", netIf.InCompressed, iface)
		ret[idx*23+7] = models.CounterValue("net.if.in.multicast", netIf.InMulticast, iface)
		ret[idx*23+8] = models.CounterValue("net.if.out.bytes", netIf.OutBytes, iface)
		ret[idx*23+9] = models.CounterValue("net.if.out.packets", netIf.OutPackages, iface)
		ret[idx*23+10] = models.CounterValue("net.if.out.errors", netIf.OutErrors, iface)
		ret[idx*23+11] = models.CounterValue("net.if.out.dropped", netIf.OutDropped, iface)
		ret[idx*23+12] = models.CounterValue("net.if.out.fifo.errs", netIf.OutFifoErrs, iface)
		ret[idx*23+13] = models.CounterValue("net.if.out.collisions", netIf.OutCollisions, iface)
		ret[idx*23+14] = models.CounterValue("net.if.out.carrier.errs", netIf.OutCarrierErrs, iface)
		ret[idx*23+15] = models.CounterValue("net.if.out.compressed", netIf.OutCompressed, iface)
		ret[idx*23+16] = models.CounterValue("net.if.total.bytes", netIf.TotalBytes, iface)
		ret[idx*23+17] = models.CounterValue("net.if.total.packets", netIf.TotalPackages, iface)
		ret[idx*23+18] = models.CounterValue("net.if.total.errors", netIf.TotalErrors, iface)
		ret[idx*23+19] = models.CounterValue("net.if.total.dropped", netIf.TotalDropped, iface)
		ret[idx*23+20] = models.GaugeValue("net.if.speed.bits", netIf.SpeedBits, iface)
		ret[idx*23+21] = models.CounterValue("net.if.in.percent", netIf.InPercent, iface)
		ret[idx*23+22] = models.CounterValue("net.if.out.percent", netIf.OutPercent, iface)
	}
	return ret
}
