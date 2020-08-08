package app

import (
	"github.com/gin-gonic/gin"
)

type HttpServer interface {
	Start(armor *Armor)
	SetupRoutes(routes Routes) *gin.Engine
}

type RoutesMaker func(*gin.Engine)

type Routes []Route

type Route struct {
	Path        string
	Method      string
	Handler     interface{}
	Middlewares gin.HandlersChain
	Children    Routes
	DisplayName string
	Description string
}
