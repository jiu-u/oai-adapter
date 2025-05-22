package base

import (
	"context"
	"github.com/jiu-u/oai-adapter/pkg/task"
)

type EmptyPoller struct {
	pollFunc task.PollFunc
}

func (p *EmptyPoller) PollTask(ctx context.Context, externalID string) (completed bool, result interface{}, err error) {
	return p.pollFunc(ctx, externalID)
}

func NewPoller(pollFunc task.PollFunc) *EmptyPoller {
	return &EmptyPoller{
		pollFunc: pollFunc,
	}
}
