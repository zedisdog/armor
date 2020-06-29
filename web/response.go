package web

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func ParseValidateErrors(errors validator.ValidationErrors) gin.H {
	response := make(gin.H)
	response["message"] = "the request is validate failed"
	es := make(map[string]string)
	for _, e := range errors {
		es[strcase.ToSnake(e.Field())] = e.(error).Error()
	}

	response["errors"] = es

	return response
}
