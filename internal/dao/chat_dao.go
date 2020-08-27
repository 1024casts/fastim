package dao

import (
	"github.com/jinzhu/gorm"

	"github.com/1024casts/fastim/internal/model"
)

type ChatDao interface {
	CreateChat(db *gorm.DB, userId uint64, YUserId uint64) (*model.ChatModel, error)
	GetChat(db *gorm.DB, chatId uint64) (*model.ChatModel, error)
	UpdateChatLastMsgId(db *gorm.DB, chatId, msgId uint64) error
	IncrChatMsgNum(db *gorm.DB, chatId uint64) error
	GetChatListByChatIds(db *gorm.DB, chatIds []uint64) ([]*model.ChatModel, error)
}

type chatDao struct {
}

func NewChatDao() ChatDao {
	return &chatDao{}
}

func (repo *chatDao) CreateChat(db *gorm.DB, userId uint64, YUserId uint64) (*model.ChatModel, error) {
	chat := &model.ChatModel{
		SenderUID:   userId,
		ReceiverUID: YUserId,
	}
	err := db.Create(&chat).Error
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (repo *chatDao) GetChat(db *gorm.DB, chatId uint64) (*model.ChatModel, error) {
	chat := &model.ChatModel{}
	err := db.Model(&chat).Where("id=?", chatId).First(&chat).Error

	return chat, err
}

func (repo *chatDao) UpdateChatLastMsgId(db *gorm.DB, chatId, msgId uint64) error {
	chat := model.ChatModel{}
	return db.Model(&chat).Where("id=?", chatId).
		Updates(map[string]interface{}{"last_msg_id": msgId}).Error
}

func (repo *chatDao) IncrChatMsgNum(db *gorm.DB, chatId uint64) error {
	chat := model.ChatModel{}
	step := 1
	return db.Model(&chat).Where("id=?", chatId).
		Updates(map[string]interface{}{"msg_num": gorm.Expr("msg_num + ?", step)}).Error
}

func (repo *chatDao) GetChatListByChatIds(db *gorm.DB, chatIds []uint64) ([]*model.ChatModel, error) {
	chats := make([]*model.ChatModel, 0)

	result := db.Where("id in (?)", chatIds).Find(&chats)

	return chats, result.Error
}
