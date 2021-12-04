package cron

import (
	"agentpro/logger"
	"agentpro/metrics"
	"agentpro/models"
	"agentpro/rpc"
	"agentpro/settings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func SendToTransfer(metrics []*models.MetricValue) {
	if len(metrics) == 0 {
		return
	}

	debug := settings.Config().Debug

	if debug {
		log.Printf("=> <Total=%d> %v\n", len(metrics), metrics[0])
	}

	var resp models.TransferResponse
	rpc.SendMetrics(metrics, &resp)

	if debug {
		log.Println("<=", &resp)
	}
}

func InitDataHistory() {
	for {
		metrics.UpdateCpuStat()
		metrics.UpdateDiskStats()
		time.Sleep(time.Second)
	}
}

func Collect() {

	if !settings.Config().Transfer.Enabled {
		return
	}

	if len(settings.Config().Transfer.Addrs) == 0 {
		return
	}

	for _, v := range metrics.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*models.MetricValue) {
	t := time.NewTicker(time.Second * time.Duration(sec)).C
	for {
		<-t

		hostname, err := settings.Hostname()
		if err != nil {
			continue
		}

		mvs := []*models.MetricValue{}
		//ignoreMetrics := settings.Config().IgnoreMetrics

		for _, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			log.Println("items:   ", items)
			logger.Debug(items)

			if len(items) == 0 {
				continue
			}

			for _, mv := range items {
				mvs = append(mvs, mv)
				// if b, ok := ignoreMetrics[mv.Metric]; ok && b {
				// 	continue
				// } else {
				// 	mvs = append(mvs, mv)
				// }
			}
		}

		now := time.Now().Unix()
		for j := 0; j < len(mvs); j++ {
			mvs[j].Step = sec
			mvs[j].Endpoint = hostname
			mvs[j].Timestamp = now
		}
		fmt.Println(mvs)
		//将指定内容写入到文件中
		out, err := json.Marshal(mvs)
		ss := string(out)
		err1 := ioutil.WriteFile("./output.txt", []byte(ss), 0666)
		log.Println(ss)
		if err1 != nil {
			fmt.Println("ioutil WriteFile error: ", err)
		}
		SendToTransfer(mvs)
	}
}
