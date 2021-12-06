package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	config  clientv3.Config
	client  *clientv3.Client
	err     error
	kv      clientv3.KV
	putResp *clientv3.PutResponse
)

var taskId = "000001"

func RunTask() string {
	mMap := make(map[string]interface{})
	mMap["cmd"] = "/tmp/demo.sh"
	mMap["task_id"] = taskId
	mMap["status"] = "pause" // pause 暂停
	mMap["executor"] = "bash"
	mMap["express"] = "* * * * *"
	j, err := json.Marshal(mMap)
	fmt.Println(string(j), err)
	return string(j)
}

func main() {
	// etcd客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"0.0.0.0:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	// 操作k-v 用于读写etcd的键值对
	kv = clientv3.NewKV(client)
	jsonStr := RunTask()
	fmt.Println(jsonStr)
	// TODO() 占位就可以
	// 如果在这里不指定WithPrevKV，那就无法得到putResp
	putResp, err := kv.Put(
		context.TODO(),
		"/cron/jobs/172.29.85.56/"+taskId,
		jsonStr,
		clientv3.WithPrevKV())
	fmt.Println(putResp, err)
}
