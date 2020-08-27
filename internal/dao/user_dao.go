package dao

import "github.com/1024casts/fastim/internal/model"

type UserDao interface {
}

type userDao struct {
}

func NewUserDao() UserDao {
	return &userDao{}
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
	//result := repo.db.Self.Where("id = ?", id).First(user)

	result := model.DB.Where("id in (?)", ids).Find(&users)

	return users, result.Error
}
