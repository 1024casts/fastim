package dao

import (
	"github.com/jinzhu/gorm"

	"github.com/1024casts/fastim/internal/model"
)

type ChatMsgDao interface {
	CreateChatMsg(db *gorm.DB, chatId, msgId uint64) (*model.ChatMsgModel, error)
	GetChatMsgListByChatId(db *gorm.DB, chatId uint64, chatMsgId uint64, limit int) ([]*model.ChatMsgModel, error)
}

type chatMsgDao struct{}

func NewCHatMsgDao() ChatMsgDao {
	return &chatMsgDao{}
}

func (repo *chatMsgDao) CreateChatMsg(db *gorm.DB, chatId, msgId uint64) (*model.ChatMsgModel, error) {
	chatMsg := &model.ChatMsgModel{
		ChatID: chatId,
		MsgID:  msgId,
	}
	err := db.Create(&chatMsg).Error
	return chatMsg, err
}

func (repo *chatMsgDao) GetChatMsgListByChatId(db *gorm.DB, chatId uint64, chatMsgId uint64, limit int) ([]*model.ChatMsgModel, error) {
	chatMsgList := make([]*model.ChatMsgModel, 0)
	result := db.Where("chat_id=? and id<=?", chatId, chatMsgId).Order("id desc").Limit(limit).Find(&chatMsgList)

	return chatMsgList, result.Error
}
