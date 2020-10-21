package dao

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/1024casts/fastim/internal/model"
	"github.com/pkg/errors"
)

type UserDao interface {
	CreateUser(user model.UserBaseModel) (id uint64, err error)
	GetUserById(id uint64) (*model.UserBaseModel, error)
	GetUsersByIds(ids []uint64) ([]*model.UserBaseModel, error)
	GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error)
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{
		db: db,
	}
}

func (repo *userDao) CreateUser(user model.UserBaseModel) (id uint64, err error) {
	err = model.DB.Create(&user).Error
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (repo *userDao) GetUserById(id uint64) (*model.UserBaseModel, error) {
	user := &model.UserBaseModel{}

	result := model.DB.Where("id = ?", id).First(user)

	return user, result.Error
}

func (repo *userDao) GetUsersByIds(ids []uint64) ([]*model.UserBaseModel, error) {
	users := make([]*model.UserBaseModel, 0)

	result := model.DB.Where("id in (?)", ids).Find(&users)

	return users, result.Error
}

// GetUserByPhone 根据手机号获取用户
func (repo *userDao) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := repo.db.Where("phone = ?", phone).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "[user_repo] get user err by phone")
	}

	return &user, nil
}
