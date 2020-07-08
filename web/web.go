package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/armor/log"
	"net/http"
	"sync"
	"time"
)

func Start(cxt context.Context, wg *sync.WaitGroup, makeRoutes MakeRoutes) {
	srv := &http.Server{
		Handler: SetupRoutes(makeRoutes),
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

func SetupRoutes(makeRoutes MakeRoutes) *gin.Engine {
	r := gin.Default()
	makeRoutes(r)
	return r
}
