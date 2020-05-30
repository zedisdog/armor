package auth

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func Auth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	claims, err := ParseToken(token)
	if err != nil {
		_ = c.AbortWithError(401, err)
	}
	id, err := strconv.ParseUint(claims.Id, 10, 64)
	c.Set("id", id)
}
