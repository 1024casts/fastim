package service

import (
	"github.com/1024casts/banhui/model"
	"github.com/1024casts/fastim/internal/dao"
)

// 直接初始化，可以避免在使用时再实例化
//var UserService = NewUserService()

type IMService interface {
	// user chat
	GetUserChatList(userId uint64, lastMId uint64, limit int) ([]*model.UserChatModel, error)

	// msg
	GetMsgListByMsgIds(msgIds []uint64) (map[uint64]*model.MsgModel, error)
	GetNewMsgNumData(userId uint64) (*newMsgNumStat, error)
}

// 校验码服务，生成校验码和获得校验码
type imService struct {
	userRepo dao.UserDao
}

func NewIMService() IMService {
	return &imService{
		userRepo: dao.NewUserDao(),
	}
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
