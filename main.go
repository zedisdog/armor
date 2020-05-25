package armor

import (
	"context"
	"github.com/zedisdog/armor/queue"
	"github.com/zedisdog/armor/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Start(makeRoutes *web.MakeRoutes) error {
	cxt, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	err := queue.Start(cxt, &wg)
	defer queue.Close()
	if err != nil {
		return err
	}
	web.Start(cxt, &wg, makeRoutes)

	<-sigs
	cancel()
	wg.Wait()

	return nil
}
