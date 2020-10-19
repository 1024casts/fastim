package service

import (
	"context"
	"strconv"

	"github.com/1024casts/fastim/internal/dao"
	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/snake/pkg/token"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var UserSvc = NewUserService()

type UserService interface {
	CreateUser(user model.UserBaseModel) (id uint64, err error)
	GetUserById(id uint64) (*model.UserBaseModel, error)
	BatchGetUserListByIds(id []uint64) (map[uint64]*model.UserBaseModel, error)
	PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error)
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

func (srv *userService) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	userModel, err := srv.userRepo.GetUserByPhone(ctx, phone)
	if err != nil || gorm.IsRecordNotFoundError(err) {
		return userModel, errors.Wrapf(err, "get user info err from db by phone: %d", phone)
	}

	return userModel, nil
}

// PhoneLogin 邮箱登录
func (srv *userService) PhoneLogin(ctx context.Context, phone int64, verifyCode int) (tokenStr string, err error) {
	// 如果是已经注册用户，则通过手机号获取用户信息
	u, err := srv.GetUserByPhone(ctx, phone)
	if err != nil {
		return "", errors.Wrapf(err, "[login] get u info err")
	}

	// 否则新建用户信息, 并取得用户信息
	if u.ID == 0 {
		u := model.UserBaseModel{
			Phone:    phone,
			Username: strconv.Itoa(int(phone)),
		}
		u.ID, err = srv.userRepo.CreateUser(u)
		if err != nil {
			return "", errors.Wrapf(err, "[login] create user err")
		}
	}

	// 签发签名 Sign the json web token.
	tokenStr, err = token.Sign(ctx, token.Context{UserID: u.ID, Username: u.Username}, "")
	if err != nil {
		return "", errors.Wrapf(err, "[login] gen token sign err")
	}
	return tokenStr, nil
}

// 批量获取
func (srv *userService) BatchGetUserListByIds(id []uint64) (map[uint64]*model.UserBaseModel, error) {
	userModels, err := srv.userRepo.GetUsersByIds(id)
	retMap := make(map[uint64]*model.UserBaseModel)

	if err != nil {
		return retMap, errors.Wrapf(err, "get user model err from db by id: %v", id)
	}

	for _, v := range userModels {
		retMap[v.ID] = v
	}

	return retMap, nil
}
