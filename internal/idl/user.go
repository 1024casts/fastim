package idl

import (
	"github.com/1024casts/banhui/model"
	"github.com/1024casts/banhui/util"
)

type TransUserInput struct {
	CurUser       *model.UserModel
	User          *model.UserModel
	Point         int
	ActivityCount int
	IsFollow      int
	IsFollowed    int
	InviteStatus  int
}

// 组装数据并输出
// 对外暴露的user结构，都应该经过此结构进行转换
func TransUser(input *TransUserInput) *model.UserInfo {
	user := input.User

	avatar := ""
	if user != nil && user.Id > 0 {
		avatar = user.Avatar
	}

	return &model.UserInfo{
		Id:             user.Id,
		Username:       user.Username,
		Avatar:         util.GetAvatarUrl(avatar),
		Sex:            user.Sex,
		PostCount:      user.PostCount,
		CommentCount:   user.CommentCount,
		ReplyCount:     user.ReplyCount,
		FeedCount:      user.PostCount + user.CommentCount + user.ReplyCount, // 动态数：贴子数+评论数+回复数
		FollowingCount: user.FollowingCount,                                  // 关注数，也是同学数
		FollowerCount:  user.FollowerCount,
		ClassCount:     user.ClassCount, // 班级数，包含: 加入的和创建的
		LikeCount:      0,
		PointCount:     input.Point,
		ActivityCount:  input.ActivityCount,
		IsFollow:       input.IsFollow,
		IsFollowed:     input.IsFollowed,
		InviteStatus:   input.InviteStatus,
	}
}
