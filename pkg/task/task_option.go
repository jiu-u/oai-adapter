package task

import "time"

// PollTaskOptions 允许自定义轮询行为
type PollTaskOptions struct {
	MaxAttempts     int           // 最大尝试次数
	InitialInterval time.Duration // 初始间隔时间
	MaxInterval     time.Duration // 最大间隔时间
	BackoffFactor   float64       // 退避系数
	Timeout         time.Duration // 总超时时间
}
