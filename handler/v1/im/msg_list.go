package im

import (
	"strconv"
	"sync"

	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/1024casts/fastim/handler"
	"github.com/1024casts/fastim/internal/idl"
	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/fastim/internal/rstore"
	"github.com/1024casts/fastim/internal/service"
	"github.com/1024casts/fastim/pkg/utils"
)

// @Summary 会话消息列表
// @Description 拉取指定会话中的消息列表(默认返回最新的20条)，主要是拉取历史消息，下拉获取更多历史消息，增量消息通过/im/msg来获取，是否有新消息通过sync/status(定期轮询接口)中的has_new_msg来判断
// @Tags 私信
// @Accept  json
// @Produce  json
// @Param lastCMId query string false "分页参数, 默认为0, 为0时获取最新数据，否则获取历史消息"
// @Success 200 {object} handler.Response "{"code":0,"message":"OK","data":null}"
// @Router /im/msg/list [get]
func MsgList(c *gin.Context) {

	// 当前用户信息
	userId := handler.GetUserId(c)
	curUser, err := service.UserSvc.GetUserById(userId)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// 获取对方用户信息
	userIdStr := c.Query("user_id")
	YUserId, _ := strconv.Atoi(userIdStr)
	_, err = service.UserSvc.GetUserById(uint64(YUserId))
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// 获取分页参数
	lastCMIdStr := c.Query("lastCMId")
	lastCMId, _ := strconv.Atoi(lastCMIdStr)

	// 获取会话
	chatResp, err := service.ImSvc.FindChat(userId, uint64(YUserId), false)
	if err != nil {
		log.Warnf("[msg_list] find chat err, %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}
	chat := chatResp.Chat
	// userChat := chatResp.UserChat
	chatId := chat.ID

	limit := 3
	var pageValue uint64
	if lastCMId > 0 {
		pageValue = uint64(lastCMId)
	}
	// 获取会话消息
	chatMsgList, err := service.ImSvc.GetChatMsgListByChatId(chatId, pageValue, limit+1)
	if err != nil {
		log.Warnf("[msg_list] get chat msg err, %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// 处理分页
	hasMore := 0
	if len(chatMsgList) > limit {
		hasMore = 1
		// 最小的一个
		lastMsg := chatMsgList[len(chatMsgList)-1]
		pageValue = lastMsg.ID

		chatMsgList = chatMsgList[0 : len(chatMsgList)-1]
	}

	infos := make([]*model.MsgInfo, 0)
	var msgIds []uint64
	for _, obj := range chatMsgList {
		msgIds = append(msgIds, obj.MsgID)
	}

	wg := sync.WaitGroup{}
	list := model.MsgList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.MsgInfo, len(chatMsgList)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 获取消息列表
	// 将msgIds 由从大到小反转为从小到大
	msgIds = utils.Uint64SliceReverse(msgIds)
	msgMap, err := service.ImSvc.GetMsgListByMsgIds(msgIds)
	if err != nil {
		log.Warnf("[im] get msg list err, %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	for _, cm := range chatMsgList {
		wg.Add(1)
		go func(p *model.ChatMsgModel) {
			defer wg.Done()

			list.Lock.Lock()
			defer list.Lock.Unlock()

			msg, ok := msgMap[p.MsgID]
			if !ok {
				errChan <- errors.Wrapf(errors.New("msg_id not in map"), "msg_id: %d", p.MsgID)
				return
			}
			transMsgInput := &idl.TransMsgInput{
				CurUser: curUser,
				Msg:     msg,
			}
			msgInfo, err := idl.TransMsg(transMsgInput)
			if err != nil {
				log.Warnf("[msglist] trans msg err %v", err)
				errChan <- err
				return
			}

			list.IdMap[p.MsgID] = msgInfo
		}(cm)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		log.Warnf("[msg_list] chan err %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	for _, id := range msgIds {
		infos = append(infos, list.IdMap[id])
	}

	// 清空新消息状态
	rstore.NewNewMsgStatus().Del(userId)

	// 清空新消息数
	rstore.NewNewMsgNum().DelOne(userId, uint64(YUserId))

	handler.SendResponse(c, nil, handler.ListResponse{
		HasMore:   hasMore,
		PageKey:   "lastCMId",
		PageValue: int(pageValue),
		Items:     infos,
	})
}
