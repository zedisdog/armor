package queue

import (
	"time"
)

// Context 自定义context
type Context struct {
	PutBack *Release
	JobID   uint64
	Payload interface{}
}

// Deadline 没有使用
func (q Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done 没有使用
func (q Context) Done() <-chan struct{} {
	return nil
}

// Err 没有使用
func (q Context) Err() error {
	return nil
}

// Value 没有使用
func (q Context) Value(key interface{}) interface{} {
	return nil
}

// Release 释放任务时的可选参数
type Release struct {
	Pri   uint32
	Delay time.Duration
}

// Release 释放任务方法
func (q *Context) Release(ops ...func(release *Release)) {
	release := &Release{
		Pri:   1024,
		Delay: 0 * time.Second,
	}
	for _, fn := range ops {
		fn(release)
	}

	q.PutBack = release
}

// NewContext 实例化自定义context
func NewContext(jobID uint64) *Context {
	return &Context{
		JobID: jobID,
	}
}
