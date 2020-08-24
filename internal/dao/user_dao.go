package dao

type UserDao interface {
}

type userDao struct {
}

func NewUserDao() UserDao {
	return &userDao{}
}
