package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/1024casts/fastim/internal/dao"
	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/fastim/internal/rstore"
	"github.com/1024casts/snake/pkg/log"
)

const (
	// 最大消息id
	MaxMsgID = 0xffffffffffff

	// 基础消息类型 1~50
	MsgTypeText  int = 1 //文本消息
	MsgTypePic   int = 2 //图片消息
	MsgTypeVoice int = 3 //语音消息
	MsgTypeVideo int = 4 //视频消息

	// 系统消息类型 51~100
	MsgTypeNoticeText int = 51 // 文字通知

	// 业务消息类型 101+
	MsgTypeInviteJoinClass int = 101 // 文章卡片消息

	// 接收类型
	ReceiveTypeBoth = 0 // 双方接收
	ReceiveTypePeer = 1 // 仅对方接收
	ReceiveTypeSelf = 2 // 仅自己接收
)

// 直接初始化，可以避免在使用时再实例化
//var ImService = NewIMService()

type IMService interface {
	// chat
	FindChat(userId uint64, YUserId uint64, isCreate bool) (*ChatResp, error)

	// msg
	SendMsg(input model.SendMsgInput) (*model.MsgModel, error)
	GetMsgListByMsgIds(msgIds []uint64) (map[uint64]*model.MsgModel, error)
	GetNewMsgNumData(userId uint64) (*newMsgNumStat, error)

	// user chat
	GetUserChatList(userId uint64, lastMId uint64, limit int) ([]*model.UserChatModel, error)
}

// 校验码服务，生成校验码和获得校验码
type imService struct {
	userRepo     dao.UserDao
	chatRepo     dao.ChatDao
	userChatRepo dao.UserChatDao
	msgRepo      dao.MsgDao
}

func NewIMService() IMService {
	return &imService{
		userRepo: dao.NewUserDao(),
	}
}

type ChatResp struct {
	Chat     *model.ChatModel
	UserChat *model.UserChatModel
}

// 查找2个用户之间的会话
// isCreate: 会话不存在时是否创建
// return: 存在返回user_chat与chat信息
func (i *imService) FindChat(userId uint64, YUserId uint64, isCreate bool) (*ChatResp, error) {
	if userId == YUserId {
		return nil, errors.New("do not chat with self")
	}

	chatResp := &ChatResp{}
	db := model.GetDB()

	userChat, err := i.userChatRepo.GetUserChat(db, userId, YUserId)
	if err != nil {
		log.Warnf("[imService] repo get user chat err, %v", err)
		return nil, err
	}

	if userChat.ChatID > 0 {
		chat, err := i.chatRepo.GetChat(db, userChat.ChatID)
		if err != nil {
			log.Warnf("[imService] repo get chat err, %v", err)
			return nil, err
		}

		chatResp.Chat = chat
		chatResp.UserChat = userChat

		return chatResp, nil
	} else {
		// 会话不存在，则创建会话
		if isCreate == false {
			return chatResp, nil
		}

		chatResp, err = i.createChat(userId, YUserId)
		if err != nil {
			log.Warnf("[imService] create chat  err, %v", err)
			return nil, err
		}

		return chatResp, nil
	}
}

func (i *imService) createChat(userId uint64, YUserId uint64) (*ChatResp, error) {
	db := model.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("[i] defer recover err, %v", r)
			tx.Rollback()
		}
	}()

	chat, err := i.chatRepo.CreateChat(tx, userId, YUserId)
	if err != nil {
		log.Warnf("[i] create chat err, %v", err)
		tx.Rollback()
		return nil, err
	}
	chatId := chat.ID

	// 写自己的user_chat
	userChat, err := i.userChatRepo.CreateUserChat(tx, chatId, userId, YUserId)
	if err != nil {
		log.Warnf("[i] create chat err, %v", err)
		tx.Rollback()
		return nil, err
	}

	// 写对方的user_chat
	_, err = i.userChatRepo.CreateUserChat(tx, chatId, YUserId, userId)
	if err != nil {
		log.Warnf("[i] create chat err, %v", err)
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		log.Warnf("[i] db commit err, %v", err)
		tx.Rollback()
		return nil, err
	}

	return &ChatResp{Chat: chat, UserChat: userChat}, nil
}

/**
 *	input:
 *  必选: user_id, yuser_id, chat_id, content
 *  可选:
 *      msg_type 默认文字消息
 *      rece_type 默认双方接收
 *      local_mid 默认0
 *      width,height 图片, 视频消息的宽高
 **/
func (i *imService) SendMsg(input model.SendMsgInput) (*model.MsgModel, error) {
	// 默认文字消息
	if input.MsgType == 0 {
		input.MsgType = MsgTypeText
	}

	contentMap := make(map[string]interface{})
	switch input.MsgType {
	case MsgTypeText:
		// 为了解析方便，所以这里也用map
		contentMap["text"] = input.Content
	case MsgTypePic:
		// 图片需要组合字段
		contentMap["pic_url"] = input.Content
		contentMap["width"] = input.Width
		contentMap["height"] = input.Height
	case MsgTypeInviteJoinClass:
		contentMap["class_id"] = input.ClassId
	default:
		contentMap["text"] = input.Content
	}

	contentByte, err := json.Marshal(contentMap)
	if err != nil {
		log.Warnf("[i] json marshal err, %v", err)
		return nil, err
	}
	content := string(contentByte)
	input.Content = content

	msg, err := i.addMsg(input)
	if err != nil {
		log.Warnf("[i] add msg err, %v", err)
		return nil, err
	}

	return msg, nil
}

// 写消息、更新关联表、写消息数到redis
func (i *imService) addMsg(input model.SendMsgInput) (*model.MsgModel, error) {
	db := model.GetDB()

	// 写入消息
	msgInput := model.MsgModel{
		ChatID:      input.ChatId,
		Content:     input.Content,
		LocalMid:    input.LocalMId,
		MsgType:     input.MsgType,
		ReceiveType: input.ReceiveType,
		UserID:      input.UserId,
		Extra:       "",
	}
	msgId, err := i.msgRepo.CreateMsg(db, msgInput)
	if err != nil {
		log.Warnf("[i] create msg err, %v", err)
		return nil, err
	}
	msgInput.ID = msgId
	fmt.Printf("msgInput: %v", msgInput)

	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//// 写chat_msg
	//_, err = i.imRepo.CreateChatMsg(tx, input.ChatId, msgId)
	//if err != nil {
	//	log.Warnf("[i] create chat msg err, %v", err)
	//	tx.Rollback()
	//	return nil, err
	//}
	//
	//// 写自己的 user_msg
	//if input.ReceiveType == ReceiveTypeBoth || input.ReceiveType == ReceiveTypeSelf {
	//	_, err := i.imRepo.CreateUserMsg(tx, input.UserId, msgId)
	//	if err != nil {
	//		log.Warnf("[i] create user msg err, %v", err)
	//		tx.Rollback()
	//		return nil, err
	//	}
	//}
	// 写对方的
	//if input.ReceiveType == ReceiveTypeBoth || input.ReceiveType == ReceiveTypePeer {
	//	_, err := i.imRepo.CreateUserMsg(tx, input.YUserId, msgId)
	//	if err != nil {
	//		log.Warnf("[i] create user msg err, %v", err)
	//		tx.Rollback()
	//		return nil, err
	//	}
	//}

	// 更新chat的 last_mid
	err = i.chatRepo.UpdateChatLastMsgId(tx, input.ChatId, msgId)
	if err != nil {
		log.Warnf("[i] update chat last mid err, %v", err)
		tx.Rollback()
		return nil, err
	}

	// 增加chat消息数
	err = i.chatRepo.IncrChatMsgNum(tx, input.ChatId)
	if err != nil {
		log.Warnf("[i] incr chat msg num err, %v", err)
		tx.Rollback()
		return nil, err
	}

	// 增加user_chat 消息数
	// 更新发送者消息计数, 通知消息或仅对方可见消息不计数
	if input.ReceiveType != ReceiveTypePeer {
		err = i.userChatRepo.IncrUserChatMsgNum(tx, input.UserId, input.YUserId)
		if err != nil {
			log.Warnf("[i] incr user chat msg num err, %v", err)
			tx.Rollback()
			return nil, err
		}
	}

	// 更新自己的last mid
	if input.ReceiveType == ReceiveTypeBoth || input.ReceiveType == ReceiveTypeSelf {
		err = i.userChatRepo.UpdateUserChatLastMsgId(tx, input.UserId, input.YUserId, msgId)
		if err != nil {
			log.Warnf("[i] update user chat lst msg id err, %v", err)
			tx.Rollback()
			return nil, err
		}
	}

	// 更新对方的last mid
	if input.ReceiveType == ReceiveTypeBoth || input.ReceiveType == ReceiveTypePeer {
		err = i.userChatRepo.UpdateUserChatLastMsgId(tx, input.YUserId, input.UserId, msgId)
		if err != nil {
			log.Warnf("[i] update user chat lst msg id err, %v", err)
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		log.Warnf("[i] tx commit err, %v", err)
		tx.Rollback()
		return nil, err
	}

	// 更新自己的新消息状态 redis
	if input.ReceiveType == ReceiveTypeBoth || input.ReceiveType == ReceiveTypeSelf {
		rstore.NewMsgStatus.Set(input.UserId)
	}

	// 更新对方的新消息状态
	if input.ReceiveType == ReceiveTypeBoth || input.ReceiveType == ReceiveTypePeer {
		rstore.NewMsgStatus.Set(input.YUserId)
	}

	// 更新对方会话的新消息数
	rstore.NewNewMsgNum().Incr(input.YUserId, input.UserId)

	return &msgInput, nil
}

func (i *imService) GetUserChatList(userId uint64, lastMId uint64, limit int) ([]*model.UserChatModel, error) {

	return nil, nil
}

func (i *imService) GetMsgListByMsgIds(msgIds []uint64) (map[uint64]*model.MsgModel, error) {

	return nil, nil
}

type newMsgNumStat struct {
	TotalNum   int            `json:"total_num"`
	MsgNumList map[uint64]int `json:"msg_num_list"`
}

func (i *imService) GetNewMsgNumData(userId uint64) (*newMsgNumStat, error) {

	return nil, nil
}
