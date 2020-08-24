package dao

type MsgDao interface {
}

type msgDao struct {
}

func NewMsgDao() MsgDao {
	return &msgDao{}
}
