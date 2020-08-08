package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/armor/app"
)

func InjectApp(a *app.Armor) func(*gin.Context) {
	return func(c *gin.Context) {
		c.Set("app", a)
		c.Next()
	}
}
