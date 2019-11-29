package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	ctx.JSON(http.StatusOK, resp)
}

func Fail(ctx *gin.Context, err error) {
	resp := Response{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, resp)
}
