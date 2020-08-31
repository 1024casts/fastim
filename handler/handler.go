package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/pkg/errno"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ListResponse struct {
	HasMore   int         `json:"has_more"`
	PageKey   string      `json:"page_key"`
	PageValue int         `json:"page_value"`
	Items     interface{} `json:"items"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	// always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func GetUserId(c *gin.Context) uint64 {
	if c == nil {
		return 0
	}
	if v, exists := c.Get("uid"); exists {
		uid, ok := v.(uint64)
		if !ok {
			return 0
		}

		return uid
	}
	return 0
}
