package cron

import (
	"agentpro/logger"
	"agentpro/settings"
	"context"
	"fmt"
	mvccpb2 "github.com/coreos/etcd/mvcc/mvccpb"
	jsoniter "github.com/json-iterator/go"
	Cron "github.com/robfig/cron/v3"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

const (
	KeyCreateChangeEvent = iota
	KeyUpdateChangeEvent
)

type KeyChangeEvent struct {
	Type  int
	Key   string
	Value []byte
}

type WatchKeyChangeResponse struct {
	Event      chan *KeyChangeEvent
	CancelFunc context.CancelFunc
	Watcher    clientv3.Watcher
}

func InitWatchTask() {
	if settings.Config().Registrar.Enable {
		WatchTaskList(GetTasksListUrl())
	}

}

func WatchTaskList(key string) (KeyChangeEventResponse *WatchKeyChangeResponse) {
	/*
		解决了Agent在运行的过程中，任务动态执行和调度
	*/

	watcher := clientv3.NewWatcher(client)
	watchChans := watcher.Watch(context.Background(), key, clientv3.WithPrefix())

	KeyChangeEventResponse = &WatchKeyChangeResponse{
		Event:   make(chan *KeyChangeEvent, 250),
		Watcher: watcher,
	}

	go func() {

		for ch := range watchChans {
			if ch.Canceled {
				goto End
			}
			for _, event := range ch.Events {
				handleKeyChangeEvent(event, KeyChangeEventResponse.Event)
			}
		}
	End:
		logger.StartupDebug("the watcher lose for key:", key)
	}()

	return
}

func handleKeyChangeEvent(event *clientv3.Event, events chan *KeyChangeEvent) {
	changeEvent := &KeyChangeEvent{
		Key: string(event.Kv.Key),
	}

	kv := clientv3.NewKV(client)

	switch event.Type {

	case mvccpb2.Event_EventType(mvccpb.PUT):
		mapTaskInfo := make(map[string]interface{})
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if event.IsCreate() {
			err := json.Unmarshal(event.Kv.Value, &mapTaskInfo)
			if err != nil {
				fmt.Println("Umarshal failed", err)
			}
			task_id := mapTaskInfo["task_id"].(string)
			cmd := mapTaskInfo["cmd"].(string)
			express := mapTaskInfo["express"].(string)
			exec := mapTaskInfo["executor"].(string)
			status := mapTaskInfo["status"].(string)

			if status == settings.TaskStatusRunning {
				logger.StartupDebug("启动任务:", task_id)
				cronid := CronStart(cmd, task_id, exec, express)
				UpdateCronId(int(cronid), task_id)
			}

			changeEvent.Type = KeyCreateChangeEvent
		} else {
			err := json.Unmarshal(event.Kv.Value, &mapTaskInfo)
			if err != nil {
				fmt.Println("Umarshal failed", err)
			}
			task_id := mapTaskInfo["task_id"].(string)
			status := mapTaskInfo["status"].(string)

			if status == settings.TaskStatusPause {
				mapTaskInfo := make(map[string]int)
				getResp, _ := kv.Get(context.TODO(), GetTasksCronIDUrl()+task_id)
				for _, kvpair := range getResp.Kvs {
					json.Unmarshal(kvpair.Value, &mapTaskInfo)
					cron_id := mapTaskInfo["cron_id"]
					id := Cron.EntryID(cron_id)
					XCron.Remove(id)
					DeleteTask(task_id)
				}
			}
			changeEvent.Type = KeyUpdateChangeEvent
		}

		changeEvent.Value = event.Kv.Value
	}
	events <- changeEvent
}
