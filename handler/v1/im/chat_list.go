package im

import (
	"github.com/gin-gonic/gin"
)

// @Summary 会话列表
// @Description
// @Tags 私信
// @Accept  json
// @Produce  json
// @Param lastMId query string true "消息id，用于分页"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /im/chatlist [get]
func ChatList(c *gin.Context) {

}
