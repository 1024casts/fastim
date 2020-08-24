package service

import (
	"github.com/1024casts/banhui/model"
	"github.com/1024casts/fastim/internal/dao"
)

type UserService interface {
	GetUserById(id uint64) (*model.UserModel, error)
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

func (srv *userService) GetUserById(id uint64) (*model.UserModel, error) {

	return nil, nil
}
