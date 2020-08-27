package dao

import (
	"github.com/jinzhu/gorm"

	"github.com/1024casts/fastim/internal/model"
)

type UserChatDao interface {
	CreateUserChat(db *gorm.DB, chatId uint64, userId, YUserId uint64) (*model.UserChatModel, error)
	GetUserChat(db *gorm.DB, userId uint64, YUserId uint64) (*model.UserChatModel, error)
	IncrUserChatMsgNum(db *gorm.DB, userId, YUserId uint64) error
	UpdateUserChatLastMsgId(db *gorm.DB, userId, YUserId, msgId uint64) error
	UpdateUserChatDelMsgId(db *gorm.DB, userId, YUserId, delMsgId uint64) error
	UpdateUserChatClearMsgId(db *gorm.DB, userId, YUserId, clearMsgId uint64) error
	GetUserChatList(db *gorm.DB, userId uint64, lastMId uint64, limit int) ([]*model.UserChatModel, error)
}

type userChatDao struct {
}

func NewUserChatDao() UserChatDao {
	return &userChatDao{}
}

func (repo *userChatDao) CreateUserChat(db *gorm.DB, chatId uint64, userId, YUserId uint64) (*model.UserChatModel, error) {
	userChat := &model.UserChatModel{
		ChatID:  chatId,
		UserID:  userId,
		YUserID: YUserId,
	}
	err := db.Create(&userChat).Error

	return userChat, err
}

func (repo *userChatDao) GetUserChat(db *gorm.DB, userId uint64, YUserId uint64) (*model.UserChatModel, error) {
	userChat := new(model.UserChatModel)
	err := db.Model(&userChat).Where("user_id=? and yuser_id=?", userId, YUserId).First(userChat).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return userChat, err
	}

	return userChat, nil
}

func (repo *userChatDao) IncrUserChatMsgNum(db *gorm.DB, userId, YUserId uint64) error {
	userChat := model.UserChatModel{}
	step := 1
	return db.Model(&userChat).Where("user_id=? and yuser_id=?", userId, YUserId).
		Updates(map[string]interface{}{"msg_num": gorm.Expr("msg_num + ?", step)}).Error
}

func (repo *userChatDao) UpdateUserChatLastMsgId(db *gorm.DB, userId, YUserId, msgId uint64) error {
	userChat := model.UserChatModel{}
	return db.Model(&userChat).Where("user_id=? and yuser_id=?", userId, YUserId).
		Updates(map[string]interface{}{"last_msg_id": msgId}).Error
}

func (repo *userChatDao) UpdateUserChatDelMsgId(db *gorm.DB, userId, YUserId, delMsgId uint64) error {
	userChat := model.UserChatModel{}
	return db.Model(&userChat).Where("user_id=? and yuser_id=?", userId, YUserId).
		Updates(map[string]interface{}{"del_msg_id": delMsgId}).Error
}

func (repo *userChatDao) UpdateUserChatClearMsgId(db *gorm.DB, userId, YUserId, clearMsgId uint64) error {
	userChat := model.UserChatModel{}
	return db.Model(&userChat).Where("user_id=? and yuser_id=?", userId, YUserId).
		Updates(map[string]interface{}{"clear_msg_id": clearMsgId}).Error
}

func (repo *userChatDao) GetUserChatList(db *gorm.DB, userId uint64, lastMId uint64, limit int) ([]*model.UserChatModel, error) {
	userChatList := make([]*model.UserChatModel, 0)
	result := db.Where("user_id=? and last_msg_id<=? and last_msg_id>del_msg_id", userId, lastMId).
		Order("last_msg_id desc").Limit(limit).
		Find(&userChatList)

	return userChatList, result.Error
}
