package cron

import (
	"agentpro/logger"
	"agentpro/settings"
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func DeleteTask(TaskId string) error {
	kv := clientv3.NewKV(client)
	ctx, cancel := context.WithTimeout(context.Background(), 604800*time.Second)
	defer cancel()
	if _, err := kv.Delete(ctx, GetTasksIDUrl()+TaskId); err != nil {
		logger.StartupFatal("DeleteTask删除失败")
	}
	CronidKey := GetTasksCronIDUrl() + TaskId
	kv.Delete(ctx, CronidKey)
	return nil
}

func UpdateCronId(cronId int, taskId string) {
	kv := clientv3.NewKV(client)
	mapCronId := make(map[string]int)
	mapCronId["cron_id"] = cronId
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	updateValue, _ := json.Marshal(mapCronId)
	updateCronidkey := GetTasksIDUrl() + taskId
	kv.Put(context.TODO(), updateCronidkey, string(updateValue))
}

func GetAllTask() error {
	if err := initScheduler(); err != nil {
		return err
	}
	kv := clientv3.NewKV(client)
	ctx, cancel := context.WithTimeout(context.Background(), 604800*time.Second)
	defer cancel()
	if getResp, err := kv.Get(ctx, GetTasksIDUrl(), clientv3.WithPrefix()); err != nil {
		logger.StartupFatal("GetAllTasks与心跳服务器连接发生异常:", err)
	} else {
		logger.StartupInfo("获取到的计划任务列表keys:", GetTasksIDUrl())
		logger.StartupInfo("获取到的计划任务列表:", getResp.Kvs)
		maptaskInfo := make(map[string]interface{})
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		for _, kvpair := range getResp.Kvs {
			err := json.Unmarshal(kvpair.Value, &maptaskInfo)
			if err != nil {
				logger.StartupFatal("Umarshal failed:", err)
			}
			status := maptaskInfo["status"].(string)
			task_id := maptaskInfo["task_id"].(string)
			if status == settings.TaskStatusRunning {
				cmd := maptaskInfo["cmd"].(string)
				express := maptaskInfo["express"].(string)
				exec := maptaskInfo["executor"].(string)
				cronid := CronStart(cmd, task_id, exec, express)
				UpdateCronId(int(cronid), task_id)
			} else if status == settings.TaskStatusPause {
				DeleteTask(task_id)
			}
		}
	}
	return nil
}
