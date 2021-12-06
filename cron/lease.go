package cron

import (
	"agentpro/logger"
	"agentpro/settings"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	keepResp     *clientv3.LeaseKeepAliveResponse
	keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
)

func KeepaliveError() error {
	return errors.New("Keepalived error")
}

func CreateLeaseToken() (clientv3.KV, clientv3.Lease) {
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	return kv, lease
}

type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)

}

var _ error = (*RetriableError)(nil)

func RegisterPutHealth(interval time.Duration, ip string) {
	kv, lease := CreateLeaseToken() // 创建租约
	err := retry.Do(
		func() error {
			keepRespChanError := make(chan error)
			var Leaseerr error
			var float64_Seconds float64 = interval.Seconds()
			seconds := int64(float64_Seconds)
			leaseResp, err := lease.Grant(context.TODO(), seconds*5)
			if err != nil {
				logger.StartupDebug("设置租约时间失败: ", err.Error())
			} else {
				leaseId := leaseResp.ID
				logger.StartupInfo("租约ID: ", leaseId)
				// 自动续租（他的底层会每次将租约信息扔到<- chan *clientv3.LeaseKeepAliveResponse管道中）
				if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
					logger.StartupDebug("keepalive failed", err)
				}

				go func() {
					for {
						select {
						case keepResp = <-keepRespChan:
							if keepResp == nil {
								logger.StartupInfo("租约失效")
								goto END // 失效必须跳出循环，才能让agent捕获异常
							}
						}

					}
				END:
					Leaseerr = KeepaliveError()
					keepRespChanError <- Leaseerr
				}()
				healthInfo := make(map[string]string)
				healthInfo["status"] = "Online"
				healthInfo["version"] = settings.VERSION
				updateValue, _ := json.Marshal(healthInfo)
				if _, err := kv.Put(context.TODO(), "/agent/health/"+ip, string(updateValue), clientv3.WithLease(leaseResp.ID)); err != nil {
					logger.StartupDebug("租约PUT失败：", err.Error())
				}
				err = <-keepRespChanError
			}
			return err
		},
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return 3 * time.Second
		}),
	)
	fmt.Println(err)
}
