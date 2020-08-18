package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/zedisdog/armor/app"
	"github.com/zedisdog/armor/log"
	"sync"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

type Queue struct {
	host       string
	workNum    int
	jobTimeout int
	conns      chan *beanstalk.Conn
	app        *app.Armor
}

// Close 关闭
func (q *Queue) Close() {
	for i := 0; i < len(q.conns); i++ {
		conn := <-q.conns
		conn.Close()
	}
}

// Start 开始队列
func (q *Queue) Start(a *app.Armor) error {
	q.app = a
	for i := 0; i < q.workNum; i++ {
		q.startWorkers(a.CancelCxt, a.Wg)
	}

	return nil
}

// Dispatch 下发任务
func (q *Queue) Dispatch(job app.Job, ops ...func(*app.DispatchOptions)) error {
	conn := <-q.conns
	defer func() { q.conns <- conn }()
	Register(job)
	options := &app.DispatchOptions{
		Pri:   1024,
		Delay: 0 * time.Second,
		Ttr:   time.Duration(q.jobTimeout) * time.Second,
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

func (q *Queue) getConn(cxt context.Context) (*beanstalk.Conn, error) {
	for {
		select {
		case <-cxt.Done():
			return nil, errors.New("force quit")
		case conn := <-q.conns:
			return conn, nil
		default:
			if len(q.conns) < cap(q.conns) {
				conn, err := beanstalk.Dial("tcp", q.host)
				if err != nil {
					q.Close()
					return nil, err
				}
				q.conns <- conn
			}
		}
	}
}

func (q *Queue) putConn(conn *beanstalk.Conn) {
	if len(q.conns) < cap(q.conns) {
		q.conns <- conn
	} else {
		conn.Close()
	}
}

func (q *Queue) startWorkers(cxt context.Context, wg *sync.WaitGroup) {
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

func (q *Queue) process(cxt context.Context) {
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
		cxt := app.NewContext(id, q.app)
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

func New(v *viper.Viper) app.Queue {
	if v.GetBool("mq.beanstalkd.enable") {
		return &Queue{
			host:       v.GetString("mq.beanstalkd.host"),
			workNum:    v.GetInt("mq.beanstalkd.work_num"),
			jobTimeout: v.GetInt("mq.beanstalkd.job_timeout"),
			conns:      make(chan *beanstalk.Conn, v.GetInt("mq.beanstalkd.worker_num")*v.GetInt("mq.beanstalkd.conn_cap_times")),
		}
	}
	return nil
}

var ProviderSet = wire.NewSet(New)
