package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type RoutesMaker func(*gin.Engine)

type Routes []Route

type Route struct {
	Path        string
	Method      string
	Handler     gin.HandlerFunc
	Middlewares gin.HandlersChain
	Children    Routes
	DisplayName string
	Description string
}

func MakeRoutes(r *gin.RouterGroup, routes Routes) error {
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
			case "GET":
				r.GET(value.Path, value.Handler)
			case "POST":
				r.POST(value.Path, value.Handler)
			case "PUT":
				r.PUT(value.Path, value.Handler)
			case "DELETE":
				r.DELETE(value.Path, value.Handler)
			case "HEAD":
				r.HEAD(value.Path, value.Handler)
			case "OPTIONS":
				r.OPTIONS(value.Path, value.Handler)
			default:
				return errors.New(fmt.Sprintf("unsupported http method <%s>", value.Method))
			}
		}
	}

	return nil
}
