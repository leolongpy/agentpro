package metrics

import (
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"fmt"
	"log"
)

func DeviceMetrics() (L []*models.MetricValue) {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Println(err)
		return
	}

	var diskTotal uint64 = 0
	var diskUsed uint64 = 0

	for idx := range mountPoints {
		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if err != nil {
			log.Println(err)
			continue
		}

		diskTotal += du.BlocksAll
		diskUsed += du.BlocksUsed

		tags := fmt.Sprintf("mount=%s,fstype=%s,__IP=%s", du.FsFile, du.FsVfstype, settings.IP())
		L = append(L, models.GaugeValue("df.bytes.used.percent", du.BlocksUsedPercent, tags))
		L = append(L, models.GaugeValue("df.inodes.used.percent", du.InodesUsedPercent, tags))
	}

	return
}
