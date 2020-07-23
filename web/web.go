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

func Start(cxt context.Context, wg *sync.WaitGroup, routes Routes) {
	srv := &http.Server{
		Handler: SetupRoutes(routes),
		Addr:    config.Instance().String("server.host") + ":" + strconv.Itoa(config.Instance().Int("server.port")),
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

func SetupRoutes(routes Routes) *gin.Engine {
	r := gin.Default()
	err := MakeRoutes(&r.RouterGroup, routes)
	if err != nil {
		panic(err)
	}
	return r
}
