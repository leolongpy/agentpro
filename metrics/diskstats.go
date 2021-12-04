package metrics

import (
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"fmt"
	"log"
	"strings"
	"sync"
)

var (
	diskStatsMap = make(map[string][2]*nux.DiskStats)
	dsLock       = new(sync.RWMutex)
)

func UpdateDiskStats() error {
	dsList, err := nux.ListDiskStats()
	if err != nil {
		return err
	}

	dsLock.Lock()
	defer dsLock.Unlock()
	for i := 0; i < len(dsList); i++ {
		device := dsList[i].Device
		diskStatsMap[device] = [2]*nux.DiskStats{dsList[i], diskStatsMap[device][0]}
	}
	return nil
}

func IOReadSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadSectors - arr[1].ReadSectors
}

func IOWriteSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteSectors - arr[1].WriteSectors
}

func IOMsecTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecTotal - arr[1].MsecTotal
}

func TS(arr [2]*nux.DiskStats) uint64 {
	return uint64(arr[0].TS.Sub(arr[1].TS).Nanoseconds() / 1000000)
}

func IODelta(device string, f func([2]*nux.DiskStats) uint64) uint64 {
	val, ok := diskStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

func DiskIOMetrics() (L []*models.MetricValue) {

	dsList, err := nux.ListDiskStats()
	if err != nil {
		log.Println(err)
		return
	}

	for _, ds := range dsList {
		if !ShouldHandleDevice(ds.Device) {
			continue
		}

		device := fmt.Sprintf("device=%s,__IP=%s", ds.Device, settings.IP())

		L = append(L, models.CounterValue("disk.io.read_sectors", ds.ReadSectors, device))
		L = append(L, models.CounterValue("disk.io.write_sectors", ds.WriteSectors, device))
	}
	return
}

func IOStatsMetrics() (L []*models.MetricValue) {
	dsLock.RLock()
	defer dsLock.RUnlock()

	for device, _ := range diskStatsMap {
		if !ShouldHandleDevice(device) {
			continue
		}

		tags := fmt.Sprintf("device=%s,__IP=%s", device, settings.IP())
		delta_rsec := IODelta(device, IOReadSectors)
		delta_wsec := IODelta(device, IOWriteSectors)
		use := IODelta(device, IOMsecTotal)

		duration := IODelta(device, TS)

		L = append(L, models.GaugeValue("disk.io.read_bytes", float64(delta_rsec)*512.0, tags))
		L = append(L, models.GaugeValue("disk.io.write_bytes", float64(delta_wsec)*512.0, tags))
		tmp := float64(use) * 100.0 / float64(duration)
		if tmp > 100.0 {
			tmp = 100.0
		}
	}

	return
}

func ShouldHandleDevice(device string) bool {
	normal := len(device) == 3 && (strings.HasPrefix(device, "sd") || strings.HasPrefix(device, "vd"))
	aws := len(device) >= 4 && strings.HasPrefix(device, "xvd")
	return normal || aws
}
