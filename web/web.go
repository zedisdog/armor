package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/armor/config"
	"github.com/zedisdog/armor/log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func Start(cxt context.Context, wg *sync.WaitGroup, routesMaker *RoutesMaker) {
	srv := &http.Server{
		Handler: SetupRoutes(routesMaker),
		Addr:    config.Instance().String("server.host") + ":" + strconv.Itoa(config.Conf.Int("server.port")),
	}
	wg.Add(1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Log.Error(err)
		}
		wg.Done()
	}()

	go func() {
		<-cxt.Done()
		timeOutCxt, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := srv.Shutdown(timeOutCxt)
		if err != nil {
			log.Log.Error(err)
		} else {
			log.Log.Info("server will be safe shutdown in 30s")
		}
	}()
}

func SetupRoutes(routesMaker *RoutesMaker) *gin.Engine {
	r := gin.Default()
	(*routesMaker)(r)
	return r
}
