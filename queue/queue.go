package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/zedisdog/armor/config"
	"github.com/zedisdog/armor/log"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

var Queue *queue

type queue struct {
	conns chan *beanstalk.Conn
}

func Instance() *queue {
	if Queue == nil {
		Init()
	}

	return Queue
}

func Init() {
	Queue = &queue{
		conns: make(chan *beanstalk.Conn, config.Conf.Int("mq.beanstalkd.worker_num")*config.Conf.Int("mq.beanstalkd.conn_cap_times")),
	}
}

// Close 关闭
func (q *queue) Close() {
	for i := 0; i < len(q.conns); i++ {
		conn := <-q.conns
		conn.Close()
	}
}

func (q *queue) getConn(cxt context.Context) (*beanstalk.Conn, error) {
	for {
		select {
		case <-cxt.Done():
			return nil, errors.New("force quit")
		case conn := <-q.conns:
			return conn, nil
		default:
			if len(q.conns) < cap(q.conns) {
				conn, err := beanstalk.Dial("tcp", config.Conf.String("mq.beanstalkd.host"))
				if err != nil {
					q.Close()
					return nil, err
				}
				q.conns <- conn
			}
		}
	}
}

func (q *queue) putConn(conn *beanstalk.Conn) {
	if len(q.conns) < cap(q.conns) {
		q.conns <- conn
	} else {
		conn.Close()
	}
}

// Start 开始队列
func (q *queue) Start(cxt context.Context, wg *sync.WaitGroup) error {
	for i := 0; i < config.Conf.Int("mq.beanstalkd.worker_num"); i++ {
		q.startWorkers(cxt, wg)
	}

	return nil
}

func (q *queue) startWorkers(cxt context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for {
			select {
			case <-cxt.Done():
				log.Log.Info("graceful down")
				wg.Done()
				return
			default:
				q.process(cxt)
			}
			time.Sleep(3 * time.Second)
		}
	}()
}

func (q *queue) process(cxt context.Context) {
	conn, err := q.getConn(cxt)
	if err != nil {
		log.Log.WithError(err).Error("can not make conn with beanstalkd")
		return
	}
	defer q.putConn(conn)
	id, body, err := conn.Reserve(3 * time.Second)
	if err != nil {
		if !errors.Is(err, beanstalk.ErrTimeout) {
			log.Log.WithError(err).Warn("beanstalk error")
		}
	} else {
		log.Log.WithField("job", string(body)).Info("job got")
		job := jsonToJob(body)
		cxt := NewContext(id)
		job.Handle(cxt)
		log.Log.WithField("job", string(body)).Info("job processed")
		if cxt.PutBack != nil {
			err := conn.Release(cxt.JobID, cxt.PutBack.Pri, cxt.PutBack.Delay)
			if err != nil {
				log.Log.WithField("job", string(body)).WithError(err).Info("release job failed")
				return
			}
			log.Log.WithField("job", string(body)).Info("job processed")
		} else {
			err := conn.Delete(cxt.JobID)
			if err != nil {
				log.Log.WithField("job", string(body)).WithError(err).Info("delete job failed")
				return
			}
			log.Log.WithField("job", string(body)).Info("job deleted")
		}
	}
}

// Dispatch 下发任务
func (q *queue) Dispatch(job Job, ops ...func(*DespatchOptions)) error {
	conn := <-q.conns
	defer func() { q.conns <- conn }()
	Register(job)
	options := &DespatchOptions{
		Pri:   1024,
		Delay: 0 * time.Second,
		Ttr:   time.Duration(config.Conf.Int("mq.beanstalkd.job_timeout")) * time.Second,
	}
	for _, fn := range ops {
		fn(options)
	}

	jobJSON := jobToJSON(job)

	start := time.Now()
	_, err := conn.Put(jobJSON, options.Pri, options.Delay, options.Ttr)
	cost := time.Since(start)
	fmt.Printf("%+v", cost)
	return err
}

// DespatchOptions 下发任务参数
type DespatchOptions struct {
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
