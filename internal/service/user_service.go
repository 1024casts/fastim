package service

import (
	"github.com/1024casts/fastim/internal/dao"
	"github.com/1024casts/fastim/internal/model"
	"github.com/pkg/errors"
)

var UserSvc = NewUserService()

type UserService interface {
	CreateUser(user model.UserBaseModel) (id uint64, err error)
	GetUserById(id uint64) (*model.UserBaseModel, error)
	BatchGetUserListByIds(id []uint64) (map[uint64]*model.UserBaseModel, error)
}

// 校验码服务，生成校验码和获得校验码
type userService struct {
	userRepo dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		userRepo: dao.NewUserDao(),
	}
}

func (srv *userService) CreateUser(user model.UserBaseModel) (id uint64, err error) {
	id, err = srv.userRepo.CreateUser(user)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (srv *userService) GetUserById(id uint64) (*model.UserBaseModel, error) {
	userModel, err := srv.userRepo.GetUserById(id)
	if err != nil {
		return userModel, errors.Wrapf(err, "get user info err from db by id: %d", id)
	}

	return userModel, nil
}

// 批量获取
func (srv *userService) BatchGetUserListByIds(id []uint64) (map[uint64]*model.UserBaseModel, error) {
	userModels, err := srv.userRepo.GetUsersByIds(id)
	retMap := make(map[uint64]*model.UserBaseModel)

	if err != nil {
		return retMap, errors.Wrapf(err, "get user model err from db by id: %v", id)
	}

	for _, v := range userModels {
		retMap[v.Id] = v
	}

	return retMap, nil
}
