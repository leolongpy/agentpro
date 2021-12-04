package metrics

import (
	"agentpro/logger"
	"agentpro/models"
	"agentpro/nux"
	"agentpro/settings"
	"agentpro/utils"
	"fmt"
	"os/exec"
	"strconv"
	"sync"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

func CpuIdleBusy() map[string]string {
	psLock.RLock()
	defer psLock.RUnlock()
	cmd := "/usr/bin/tsar --check 2>/dev/null|awk -F= '{print $2}'|awk '{print $1}'"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		logger.ToMOCDebug("Failed to execute command: %s", cmd)
	}
	Busy := string(out)
	if Busy == "" {
		Busy = "0.0"
	}
	Busy = utils.Strip(Busy, "\n")
	IdleFloat, _ := strconv.ParseFloat(Busy, 32)
	Idle := 100.0 - IdleFloat
	IdleFormat := strconv.FormatFloat(Idle, 'f', -1, 32)

	IdleBusyMap := make(map[string]string)
	IdleBusyMap["idle"] = IdleFormat
	IdleBusyMap["busy"] = Busy
	return IdleBusyMap
}

func CpuUser() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
}

func CpuSystem() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
}

func CpuIowait() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
}

func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

func CpuMetrics() []*models.MetricValue {
	if !CpuPrepared() {
		return []*models.MetricValue{}
	}

	tags := fmt.Sprintf("__IP=%s", settings.IP())
	CpuIdleBusyMap := CpuIdleBusy()
	idle := models.GaugeValue("cpu.idle", CpuIdleBusyMap["idle"], tags)
	busy := models.GaugeValue("cpu.busy", CpuIdleBusyMap["busy"], tags)
	user := models.GaugeValue("cpu.user", CpuUser(), tags)
	system := models.GaugeValue("cpu.system", CpuSystem(), tags)
	iowait := models.GaugeValue("cpu.iowait", CpuIowait(), tags)
	return []*models.MetricValue{idle, busy, user, system, iowait}
}
