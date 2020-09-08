package im

import (
	"github.com/1024casts/snake/pkg/log"
	"github.com/gin-gonic/gin"
)

// @Summary 收消息
// @Description
// @Tags 私信
// @Accept  json
// @Produce  json
// @Param lastMId query string true "消息id，用于分页"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /im/chatlist [get]
func Msg(c *gin.Context) {
	log.Info("ChatList function called.")

}
