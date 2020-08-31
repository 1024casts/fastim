package im

import (
	"github.com/1024casts/fastim/internal/idl"
	"github.com/gin-gonic/gin"

	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"

	"github.com/1024casts/fastim/handler"
	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/fastim/internal/service"
)

// @Summary 发送消息
// @Description 内容使用json格式，eg: {"user_id": 200002, "msg_type":1, "content": "聊天的内容或者图片key"}
// @Tags 私信
// @Accept  json
// @Produce  json
// @Param user_id body int true "聊天对象的uid"
// @Param msg_type body string true "消息类型 1:文本消息 2:图片 3:音频消息 4:视频消息 51:系统通知文本消息 101:文章卡片消息"
// @Param content body string true "消息内容，如果是文字则为内子内容，如果是图片(音频或视频)则为图片key(音频或视频)"
// @Param width body int false "上传图片的宽度"
// @Param height body int false "上传图片的高度"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /im/send [post]
func Send(c *gin.Context) {
	log.Info("Send function called.")
	// Binding the data.
	var req SendRequest
	if err := c.Bind(&req); err != nil {
		log.Warnf("[send] bind req param err: %+v", err)
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}
	log.Warnf("[send] get req param: %v", req)
	// 检查参数
	if checkSendParam(&req) == false {
		log.Warnf("[send] check req param err")
		handler.SendResponse(c, errno.ErrParam, nil)
		return
	}

	userSrv := service.NewUserService()

	// 当前用户信息
	userId := handler.GetUserId(c)
	curUser, err := userSrv.GetUserById(userId)
	if err != nil {
		log.Warnf("[send] get user info err: %+v, user_id:%d", err, userId)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// 获取对方用户的信息
	YUserId := req.UserId
	YUser, err := userSrv.GetUserById(YUserId)
	if err != nil {
		log.Warnf("[send] get user info err: %+v, yuser_id:%d", err, YUserId)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// 是否拉黑对方
	// 是否被对方拉黑

	imSrv := service.NewIMService()

	// 获取用户会话
	chatResp, err := imSrv.FindChat(userId, YUserId, true)
	if err != nil {
		log.Warnf("[send] find chat err, %v, user_id: %d, yuser_id: %d", err, userId, YUserId)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}
	chat := chatResp.Chat
	chatId := chat.ID
	//userChat := chatResp.UserChat
	receiveType := service.ReceiveTypeBoth

	// 发私信
	msgInput := model.SendMsgInput{
		UserId:      userId,
		YUserId:     YUserId,
		ChatId:      chatId,
		LocalMId:    req.LocalMId,
		MsgType:     req.MsgType,
		ReceiveType: receiveType,
		Content:     req.Content,
		Width:       req.Width,
		Height:      req.Height,
		Duration:    req.Duration,
	}
	msg, err := imSrv.SendMsg(msgInput)
	if err != nil {
		log.Warnf("[send] send msg err, %v, input: %+v", err, msgInput)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// push 对方，长连接推送 + app 推送

	// push 自己

	// return chat
	transChatInput := &idl.TransChatInput{
		CurUser: curUser,
		User:    YUser,
		Msg:     msg,
	}
	chatInfo, err := idl.TransChat(transChatInput)
	if err != nil {
		log.Warnf("[send] trans chat err, %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	handler.SendResponse(c, nil, gin.H{
		"curUser": curUser,
		"YUser":   YUser,
		"msg":     msg,
		"chat":    chatInfo,
	})
}

// 检查参数
func checkSendParam(req *SendRequest) bool {
	if req.MsgType == 0 {
		return false
	}
	switch req.MsgType {
	case service.MsgTypeText:
		if req.Content == "" {
			return false
		}
	case service.MsgTypePic:
		if req.Content == "" {
			return false
		}
		if req.Width <= 0 {
			return false
		}
		if req.Height <= 0 {
			return false
		}
	}
	return true
}
