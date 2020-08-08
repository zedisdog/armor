package middlewares

import (
	"github.com/gin-gonic/gin"
	app2 "github.com/zedisdog/armor/app"
	"github.com/zedisdog/armor/auth"
	"net/http"
	"strconv"
	"strings"
)

func Auth(c *gin.Context) {
	app := c.MustGet("app").(*app2.Armor)
	arr := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(arr) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问"})
		return
	}
	token := arr[1]
	claims, err := auth.ParseToken(token, []byte(app.Config.GetString("jwt.key")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"message": "未授权的访问"})
		return
	}
	id, err := strconv.ParseUint(claims.Id, 10, 64)
	c.Set("id", id)
}
