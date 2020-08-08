package app

import "time"

type Queue interface {
	Start(armor *Armor) error
	Dispatch(job Job, ops ...func(*DispatchOptions)) error
	Close()
}

// Job 队列任务接口
type Job interface {
	// Handle 执行任务时调用该方法
	Handle(c *Context)
}

// DispatchOptions 下发任务参数
type DispatchOptions struct {
	Pri   uint32
	Delay time.Duration
	Ttr   time.Duration
}

// WithPri beanstalk put参数pri
func WithPri(pri uint32) func(map[string]interface{}) {
	return func(options map[string]interface{}) {
		options["pri"] = pri
	}
}

// WithDelay beanstalk put参数delay
func WithDelay(delay time.Duration) func(map[string]interface{}) {
	return func(options map[string]interface{}) {
		options["delay"] = delay
	}
}

// WithTtr beanstalk put参数ttr
func WithTtr(ttr time.Duration) func(map[string]interface{}) {
	return func(options map[string]interface{}) {
		options["ttr"] = ttr
	}
}

// Context 自定义context
type Context struct {
	PutBack *Release
	JobID   uint64
	app     *Armor
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
func NewContext(jobID uint64, app *Armor) *Context {
	return &Context{
		JobID: jobID,
		app:   app,
	}
}
