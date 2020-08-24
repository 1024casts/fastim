package dao

type ChatDao interface {
}

type chatDao struct {
}

func NewChatDao() ChatDao {
	return &chatDao{}
}
