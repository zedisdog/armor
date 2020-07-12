package armor

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/zedisdog/armor/model"
	"github.com/zedisdog/armor/queue"
	"github.com/zedisdog/armor/web"
)

type armor struct {
	startQueue bool
	addr       string
}

type ConfigFunc func(*armor)

func WithQueue(enabled bool) ConfigFunc {
	return func(a *armor) {
		a.startQueue = enabled
	}
}

func SetAddr(addr string) ConfigFunc {
	return func(a *armor) {
		a.addr = addr
	}
}

func NewArmor(migrate model.AutoMigrate, configs ...ConfigFunc) *armor {
	if migrate != nil {
		migrate(model.DB)
	}
	a := &armor{}
	for _, config := range configs {
		config(a)
	}
	return a
}

func (a *armor) Start(makeRoutes web.MakeRoutes) error {
	cxt, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	if a.startQueue {
		err := queue.Start(cxt, &wg)
		defer queue.Close()
		if err != nil {
			return err
		}
	}
	if a.addr == "" {
		a.addr = ":80"
	}
	web.Start(cxt, &wg, a.addr, makeRoutes)

	<-sigs
	cancel()
	wg.Wait()

	return nil
}
