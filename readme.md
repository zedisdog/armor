# armor

## web server

使用gin

路由函数
```go
package main

import (
	"github.com/zedisdog/armor"
	"github.com/zedisdog/armor/log"
)

func MakeRoutes(r *gin.Engine) {
    r.POST("/login", v1.Login)
}

func main() {
    err := armor.Start(MakeRoutes)
    if err != nil {
        log.Log.WithError(err).Info("server start failed")
    }
}
```

## auth

使用jwt

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/zedisdog/armor/auth"
    "github.com/zedisdog/armor/model"
    "strconv"
)

type Account struct {
	model.Model
}

func MakeRoutes(r *gin.Engine) {
    r.POST("/login", generateToken)
    r.Use(auth.Auth)
    r.GET("/account", parseToken)
}

func main() {
    err := armor.Start(MakeRoutes)
    if err != nil {
        log.Log.WithError(err).Info("server start failed")
    }
}

func generateToken(c *gin.Context) {
    a := Account{Id: 123}
    token, _ := auth.GenerateToken(&a)
    c.JSON(200, gin.H{
        "token": token,
    })
}

func parseToken(c *gin.Context) {
    i, _ := c.Get("id")
    c.JSON(200, gin.H{
        "id": strconv.FormatUint(i.(uint64), 10),
    })
}
```

## orm

使用gorm

```go
package models

import "github.com/zedisdog/armor/model"

type Account struct {
	model.Model
	Username string `gorm:"type:varchar(255) COMMENT '用户名';" json:"username"`
	Password string `gorm:"type:varchar(255) COMMENT '密码';" json:"-"`
	Roles    []Role `gorm:"many2many:role_user;" json:"roles"`
}
```

## queue

使用beanstalkd

## log

使用logrus

```go
package main

import (
	"github.com/zedisdog/armor/log"
)

func main() {
    log.Log.WithError(err).Error("server start failed")
}
```

## test

使用goconvey