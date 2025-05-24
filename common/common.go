package common

import (
	"github.com/jiu-u/oai-adapter/pkg/task"
	"time"
)

// DefaultPollTaskOptions 提供合理的默认值
var DefaultPollTaskOptions = task.PollTaskOptions{
	MaxAttempts:     30,
	InitialInterval: time.Duration(float64(time.Second) * 1.5),
	MaxInterval:     time.Second * 30,
	BackoffFactor:   1.5,
	Timeout:         time.Minute * 10,
}
