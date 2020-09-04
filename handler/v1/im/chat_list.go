package im

import (
	"strconv"
	"sync"

	"github.com/1024casts/fastim/internal/idl"

	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/pkg/log"
	"github.com/gin-gonic/gin"

	"github.com/1024casts/fastim/handler"
	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/fastim/internal/rstore"
	"github.com/1024casts/fastim/internal/service"
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
	log.Info("ChatList function called.")

	// 当前用户信息
	userId := handler.GetUserId(c)
	curUser, err := service.UserSvc.GetUserById(userId)
	if err != nil {
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	limit := 10
	pageValueStr := c.Query("lastMId")
	pageValue, _ := strconv.Atoi(pageValueStr)
	int64PageValue := uint64(pageValue)

	// 获取用户会话列表 user_chat
	userChatList, err := service.ImSvc.GetUserChatList(userId, int64PageValue, limit+1)
	if err != nil {
		log.Warnf("[chat_list] get user chat list err, %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	if len(userChatList) == 0 {
		// 返回空列表
		handler.SendResponse(c, errno.OK, handler.ListResponse{})
		return
	}

	// 处理翻页
	hasMore := 0
	if len(userChatList) > limit {
		hasMore = 1
		userChatList = userChatList[0 : len(userChatList)-1]

		lastUserChat := userChatList[len(userChatList)-1]
		int64PageValue = lastUserChat.LastMsgID
	}

	infos := make([]*model.ChatInfo, 0)

	var ids []uint64
	for _, obj := range userChatList {
		ids = append(ids, obj.Id)
	}

	wg := sync.WaitGroup{}
	list := model.UserChatList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.ChatInfo, len(userChatList)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// 获取会话的新消息数 redis
	// 一个人和多个人的所有消息数

	// 获取消息列表
	var msgIds []uint64
	for _, uc := range userChatList {
		if uc.LastMsgID > uc.ClearMsgID {
			msgIds = append(msgIds, uc.LastMsgID)
		}
	}
	msgMap := make(map[uint64]*model.MsgModel)
	if len(msgIds) > 0 {
		msgMap, err = service.ImSvc.GetMsgListByMsgIds(msgIds)
		if err != nil {
			log.Warnf("[chat_list] get msg list err, %v", err)
			handler.SendResponse(c, errno.InternalServerError, nil)
			return
		}
	}

	// 获取对方用户信息
	var userIds []uint64
	for _, uc := range userChatList {
		userIds = append(userIds, uc.YUserID)
	}
	userMap := make(map[uint64]*model.UserBaseModel)
	if len(userIds) > 0 {
		userMap, err = service.UserSvc.BatchGetUserListByIds(userIds)
		if err != nil {
			log.Warnf("[chat_list] get user list err, %v", err)
			handler.SendResponse(c, errno.InternalServerError, nil)
			return
		}
	}

	newMsgNumData, err := service.ImSvc.GetNewMsgNumData(userId)
	if err != nil {
		log.Warnf("[chat_list] get total new msg num err, %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	// 循环user_chat，trans chat
	for _, uc := range userChatList {
		// 过滤用户删除 和 空会话
		if uc.DelMsgID >= uc.LastMsgID {
			continue
		}

		wg.Add(1)
		go func(p *model.UserChatModel) {
			defer wg.Done()

			list.Lock.Lock()
			defer list.Lock.Unlock()

			user, ok := userMap[p.YUserID]
			if !ok {
				errChan <- err
				return
			}

			msg, ok := msgMap[p.LastMsgID]
			if !ok {
				errChan <- err
				return
			}

			// 获取会话的新消息数
			newMsgNum, ok := newMsgNumData.MsgNumList[p.YUserID]
			if !ok {
				newMsgNum = 0
			}

			transChatInput := &idl.TransChatInput{
				CurUser:   curUser,
				User:      user,
				Msg:       msg,
				NewMsgNum: newMsgNum,
			}
			chat, err := idl.TransChat(transChatInput)
			if err != nil {
				errChan <- err
				return
			}
			list.IdMap[p.Id] = chat
		}(uc)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		log.Warnf("[chatlist] err %v", err)
		handler.SendResponse(c, errno.InternalServerError, nil)
		return
	}

	for _, id := range ids {
		infos = append(infos, list.IdMap[id])
	}

	// 清空用户新消息状态
	rstore.NewNewMsgStatus().Del(userId)

	handler.SendResponse(c, nil, handler.ListResponse{
		HasMore:   hasMore,
		PageKey:   "lastMId",
		PageValue: pageValue,
		Items:     infos,
	})
}
