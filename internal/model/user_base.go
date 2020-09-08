package model

import (
	"sync"
	"time"
)

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

type UserInfo struct {
	Id             uint64 `json:"id"`
	Phone          int    `json:"phone"`
	Username       string `json:"username"`
	Avatar         string `json:"avatar"`
	Sex            int    `json:"sex"`
	PostCount      int    `json:"post_count"`
	CommentCount   int    `json:"comment_count"`
	ReplyCount     int    `json:"reply_count"`
	FeedCount      int    `json:"feed_count"`
	FollowingCount int    `json:"following_count"`
	FollowerCount  int    `json:"follower_count"`
	ClassCount     int    `json:"class_count"`    // 总班级数
	LikeCount      int    `json:"like_count"`     // 总的喜欢数
	PointCount     int    `json:"point_count"`    // 总的斑点数
	ActivityCount  int    `json:"activity_count"` // 活动数
	IsFollow       int    `json:"is_follow"`      // 是否关注用户
	IsFollowed     int    `json:"is_followed"`    // 是否被粉丝关注
	InviteStatus   int    `json:"invite_status"`  // 邀请状态（在班级中邀请的)）
	CreatedAt      string `json:"created_at"`
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo
}
