package dao

import (
	"github.com/jinzhu/gorm"

	"github.com/1024casts/fastim/internal/model"
)

type MsgDao interface {
	CreateMsg(db *gorm.DB, msgModel model.MsgModel) (uint64, error)
	GetMsgByMsgId(db *gorm.DB, msgId uint64) (*model.MsgModel, error)
	GetMsgListByMsgIds(db *gorm.DB, msgIds []uint64) ([]*model.MsgModel, error)
	DelMsgByMsgId(db *gorm.DB, msgId uint64) error
}

type msgDao struct {
}

func NewMsgDao() MsgDao {
	return &msgDao{}
}

func (repo *msgDao) CreateMsg(db *gorm.DB, msgModel model.MsgModel) (uint64, error) {
	err := db.Create(&msgModel).Error
	return msgModel.ID, err
}

func (repo *msgDao) GetMsgByMsgId(db *gorm.DB, msgId uint64) (*model.MsgModel, error) {
	msg := &model.MsgModel{}

	result := db.Where("id=?", msgId).First(&msg)

	return msg, result.Error
}

func (repo *msgDao) GetMsgListByMsgIds(db *gorm.DB, msgIds []uint64) ([]*model.MsgModel, error) {
	msgs := make([]*model.MsgModel, 0)

	result := db.Where("id in (?)", msgIds).Find(&msgs)

	return msgs, result.Error
}

func (repo *msgDao) DelMsgByMsgId(db *gorm.DB, msgId uint64) error {
	msg := &model.MsgModel{}
	result := db.Where("id=?", msgId).Delete(&msg)

	return result.Error
}
