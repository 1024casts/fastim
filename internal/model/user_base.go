package model

import "time"

type UserBaseModel struct {
	BaseModel
	Phone          int       `gorm:"column:phone" json:"phone"`
	Username       string    `gorm:"column:username" json:"username"`
	Avatar         string    `gorm:"column:avatar" json:"avatar"`
	Password       string    `gorm:"column:password" json:"password"`
	OpenID         string    `gorm:"column:open_id" json:"open_id, omitempty"`
	UnionID        string    `gorm:"column:union_id" json:"union_id, omitempty"`
	Sex            int       `gorm:"column:sex" json:"sex"`
	PostCount      int       `gorm:"column:post_count" json:"post_count"`
	CommentCount   int       `gorm:"column:comment_count" json:"comment_count"`
	ReplyCount     int       `gorm:"column:reply_count" json:"reply_count"`
	ClassCount     int       `gorm:"column:class_count" json:"class_count"`
	FollowingCount int       `gorm:"column:following_count" json:"following_count"`
	FollowerCount  int       `gorm:"column:follower_count" json:"follower_count"`
	Status         int       `gorm:"column:status" json:"status"`
	LastLoginIP    string    `gorm:"column:last_login_ip" json:"last_login_ip"`
	LastLoginTime  time.Time `gorm:"column:last_login_time" json:"last_login_time"`
}

// TableName sets the insert table name for this struct type
func (u *UserBaseModel) TableName() string {
	return "user_base"
}
