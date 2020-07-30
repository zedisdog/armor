package test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zedisdog/armor/auth"
	"github.com/zedisdog/armor/model"
	"net/http/httptest"
)

var (
	token string
)

func Act(account model.HasId) {
	var err error
	token, err = auth.GenerateToken(account)
	if err != nil {
		panic(err)
	}
}

func Post(handler gin.HandlerFunc, data map[string]interface{}) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.POST("/", handler)

	jsonData, _ := json.Marshal(data)
	c.Request = httptest.NewRequest("POST", "/",
		bytes.NewReader(jsonData),
	)

	if token != "" {
		c.Request.Header.Set("Authorization", token)
	}

	r.ServeHTTP(w, c.Request)

	return w, c
}

func Get(handler gin.HandlerFunc) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/", nil)

	if token != "" {
		r.Use(auth.Middleware)
		c.Request.Header.Set("Authorization", token)
	}

	r.GET("/", handler)

	r.ServeHTTP(w, c.Request)

	return w, c
}

type Data map[string]interface{}
