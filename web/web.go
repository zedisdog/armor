package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"github.com/zedisdog/armor/app"
	"github.com/zedisdog/armor/log"
	"github.com/zedisdog/armor/web/middlewares"
	"net/http"
	"strconv"
	"time"
)

const POST = "POST"
const GET = "GET"
const PUT = "PUT"
const DELETE = "DELETE"
const OPTIONS = "OPTIONS"
const HEAD = "HEAD"
const STATIC = "STATIC"
const STATIC_FILE = "STATIC_FILE"
const STATIC_FS = "STATIC_FS"

type HttpServer struct {
	app  *app.Armor
	addr string
}

func (h *HttpServer) Start(a *app.Armor) {
	h.app = a
	srv := &http.Server{
		Handler: h.SetupRoutes(a.Routes),
		Addr:    h.addr,
	}
	a.Wg.Add(1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Log.Error(err)
		}
		a.Wg.Done()
	}()

	go func() {
		<-a.CancelCxt.Done()
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

func (h *HttpServer) SetupRoutes(routes app.Routes) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cros)
	r.Use(middlewares.InjectApp(h.app))
	err := MakeRoutes(&r.RouterGroup, routes)
	if err != nil {
		panic(err)
	}
	return r
}

func MakeRoutes(r *gin.RouterGroup, routes app.Routes) error {
	for _, value := range routes {
		if len(value.Children) > 0 {
			group := r.Group(value.Path)
			if len(value.Middlewares) > 0 {
				group.Use(value.Middlewares...)
			}
			err := MakeRoutes(group, value.Children)
			if err != nil {
				return err
			}
		} else {
			//if len(value.Middlewares) > 0 {
			//	missings := make(gin.HandlersChain, len(value.Middlewares))
			//	for _, middleware := range value.Middlewares {
			//		isExists := false
			//		for _, exists := range r.Handlers {
			//			if &middleware == &exists { // todo: 这个比较无效 考虑如何能比较两个函数
			//				isExists = true
			//				break
			//			}
			//		}
			//		if (!isExists) {
			//			missings = append(missings, middleware)
			//		}
			//	}
			//	r.Use(missings...)
			//}
			switch value.Method {
			case GET:
				r.GET(value.Path, value.Handler.(func(*gin.Context)))
			case POST:
				r.POST(value.Path, value.Handler.(func(*gin.Context)))
			case PUT:
				r.PUT(value.Path, value.Handler.(func(*gin.Context)))
			case DELETE:
				r.DELETE(value.Path, value.Handler.(func(*gin.Context)))
			case HEAD:
				r.HEAD(value.Path, value.Handler.(func(*gin.Context)))
			case OPTIONS:
				r.OPTIONS(value.Path, value.Handler.(func(*gin.Context)))
			case STATIC:
				r.Static(value.Path, value.Handler.(string))
			case STATIC_FILE:
				r.StaticFile(value.Path, value.Handler.(string))
			case STATIC_FS:
				r.StaticFS(value.Path, value.Handler.(http.FileSystem))
			default:
				return errors.New(fmt.Sprintf("unsupported http method <%s>", value.Method))
			}
		}
	}

	return nil
}

func New(v *viper.Viper) (app.HttpServer, error) {
	return &HttpServer{
		addr: v.GetString("server.host") + ":" + strconv.Itoa(v.GetInt("server.port")),
	}, nil
}

var ProviderSet = wire.NewSet(New)
