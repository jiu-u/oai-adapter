package task

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var NeedCancelTaskErr = errors.New("bad request")

// Status 表示任务状态
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Result 存储任务结果和查询信息
type Result struct {
	Status      Status                 // 任务状态
	Data        interface{}            // 任务结果数据
	Error       error                  // 如果任务失败，这里存储错误
	CreatedAt   time.Time              // 任务创建时间
	CompletedAt time.Time              // 任务完成时间
	UpdatedAt   time.Time              // 最后更新时间
	ExternalID  string                 // 外部系统的任务ID
	StopPolling bool                   // 标记是否停止轮询
	Metadata    map[string]interface{} // 额外元数据
}

type PollFunc func(ctx context.Context, externalID string) (completed bool, result interface{}, err error)

// Poller 定义用于轮询任务状态的接口
type Poller interface {
	// PollTask 查询任务状态，返回是否完成、结果和错误
	PollTask(ctx context.Context, externalID string) (completed bool, result interface{}, err error)
}

// Manager 管理异步任务和轮询
type Manager struct {
	tasks       sync.Map      // 存储任务ID到任务结果的映射
	ttl         time.Duration // 任务结果保留时间
	pollInitial time.Duration // 初始轮询间隔
	pollMax     time.Duration // 最大轮询间隔
	ctx         context.Context
	cancel      context.CancelFunc
}

//var _ TaskManager = (*Manager)(nil)

type TaskManager interface {
	CreatePollingTask(ctx context.Context, externalID string, poller Poller, metadata map[string]interface{}) string
	GetTaskResult(taskID string) (*Result, bool)
	CancelTask(taskID string) bool
	CleanExpiredTasks()
	Close()
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(ttl, pollInitial, pollMax time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	tm := &Manager{
		ttl:         ttl,
		pollInitial: pollInitial,
		pollMax:     pollMax,
		ctx:         ctx,
		cancel:      cancel,
	}
	// 启动过期任务清理协程
	go tm.CleanExpiredTasks()
	return tm
}

// Close 关闭任务管理器，停止所有后台协程
func (tm *Manager) Close() {
	tm.cancel() // 取消上下文，通知所有协程退出
}

// CreatePollingTask 创建需要轮询的任务并返回本地任务ID
func (tm *Manager) CreatePollingTask(
	ctx context.Context,
	externalID string,
	poller Poller,
	metadata map[string]interface{},
) string {
	taskID := ulid.Make().String() // 生成本地唯一ID

	result := &Result{
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExternalID:  externalID,
		StopPolling: false,
		Metadata:    metadata,
	}

	tm.tasks.Store(taskID, result)

	// 启动轮询协程，传入特定的轮询器
	go tm.startPolling(ctx, taskID, poller)

	return taskID
}

// startPolling 开始轮询任务状态
func (tm *Manager) startPolling(ctx context.Context, taskID string, poller Poller) {
	// 创建带有取消功能的上下文
	pollCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 获取任务信息
	resultVal, exists := tm.tasks.Load(taskID)
	if !exists {
		return
	}
	result := resultVal.(*Result)
	externalID := result.ExternalID

	// 设置初始轮询间隔
	interval := tm.pollInitial

	// 开始轮询循环
	for {
		// 检查是否应该停止轮询
		resultVal, exists := tm.tasks.Load(taskID)
		if !exists || resultVal.(*Result).StopPolling {
			return
		}

		// 等待轮询间隔
		select {
		case <-time.After(interval):
			// 继续轮询
		case <-pollCtx.Done():
			// 上下文被取消，停止轮询
			return
		}

		// 更新任务状态为处理中
		result := resultVal.(*Result)
		if result.Status == StatusPending {
			result.Status = StatusProcessing
			result.UpdatedAt = time.Now()
			tm.tasks.Store(taskID, result)
		}

		// 调用轮询器查询任务状态
		completed, data, err := poller.PollTask(pollCtx, externalID)
		// 更新最后查询时间
		result.UpdatedAt = time.Now()

		if err != nil {
			if errors.Is(err, NeedCancelTaskErr) {
				result.Status = StatusFailed
				result.StopPolling = true
				result.Data = data
				tm.tasks.Store(taskID, result)
				return
			}
			// 轮询出错，但不一定表示任务失败，可能只是临时网络问题
			// 记录错误但继续轮询
			log.Printf("Error polling task %s: %v", taskID, err)
		}
		result.Data = data
		if completed {
			result.CompletedAt = time.Now()
			result.Status = StatusCompleted
			result.StopPolling = true
			tm.tasks.Store(taskID, result)
			return
		}
		tm.tasks.Store(taskID, result)

		// 使用指数退避增加轮询间隔
		interval = time.Duration(float64(interval) * 1.5)
		if interval > tm.pollMax {
			interval = tm.pollMax
		}
	}
}

// GetTaskResult 获取任务结果
func (tm *Manager) GetTaskResult(taskID string) (*Result, bool) {
	resultVal, exists := tm.tasks.Load(taskID)
	if !exists {
		return nil, false
	}

	return resultVal.(*Result), true
}

// CancelTask 取消任务轮询
func (tm *Manager) CancelTask(taskID string) bool {
	resultVal, exists := tm.tasks.Load(taskID)
	if !exists {
		return false
	}

	result := resultVal.(*Result)
	result.StopPolling = true
	tm.tasks.Store(taskID, result)
	return true
}

// CleanExpiredTasks 定期清理过期的任务
func (tm *Manager) CleanExpiredTasks() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			tm.tasks.Range(func(key, value interface{}) bool {
				result := value.(*Result)
				// 如果任务已完成或失败，且超过TTL时间，则删除
				if (result.Status == StatusCompleted || result.Status == StatusFailed) &&
					now.Sub(result.UpdatedAt) > tm.ttl {
					tm.tasks.Delete(key)
				}
				return true
			})
		case <-tm.ctx.Done():
			return
		}
	}
}
