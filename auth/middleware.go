package auth

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func Middleware(c *gin.Context) {
	token := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	claims, err := ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(401, map[string]string{"message": "未授权的访问"})
		return
	}
	id, err := strconv.ParseUint(claims.Id, 10, 64)
	c.Set("id", id)
}
