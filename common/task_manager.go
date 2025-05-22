package common

import (
	"github.com/jiu-u/oai-adapter/pkg/task"
	"sync"
	"time"
)

var (
	defaultTaskManager task.TaskManager
	initTaskOnce       sync.Once
)

func SetDefaultTaskManager(tm task.TaskManager) {
	if defaultTaskManager != nil {
		defaultTaskManager.Close()
	}
	defaultTaskManager = tm
}

func GetDefaultTaskManager() task.TaskManager {
	initTaskOnce.Do(func() {
		defaultTaskManager = task.NewTaskManager(time.Hour*24, time.Second, time.Minute)
	})
	return defaultTaskManager
}
