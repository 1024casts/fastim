package idl

import (
	"fmt"
	"time"

	"github.com/1024casts/fastim/internal/model"
)

// TransChatInput trans chat input data
type TransChatInput struct {
	CurUser   *model.UserBaseModel
	User      *model.UserBaseModel // 对方用户信息
	Msg       *model.MsgModel
	NewMsgNum int
}

// TransChat 组装数据并输出
// 对外暴露的user结构，都应该经过此结构进行转换
func TransChat(input *TransChatInput) (*model.ChatInfo, error) {
	transUserInput := &TransUserInput{
		User: input.User,
	}

	transMsgInput := &TransMsgInput{
		CurUser: input.CurUser,
		Msg:     input.Msg,
	}

	msg, err := TransMsg(transMsgInput)
	if err != nil {
		return nil, err
	}

	return &model.ChatInfo{
		User:      TransUser(transUserInput),
		Msg:       msg,
		ShowTime:  TransChatTime(input.Msg.CreatedAt),
		NewMsgNum: input.NewMsgNum,
	}, nil
}

func TransChatTime(cTime time.Time) string {
	duration := time.Now().Unix() - cTime.Unix()

	if duration < 60 {
		//return fmt.Sprintf("%d妙前", duration)
		return fmt.Sprintf("刚刚")
	} else if duration < 3600 {
		return fmt.Sprintf("%d分钟前", duration/60)
	} else if duration < 172800 {
		// 2天
		day := time.Unix(cTime.Unix(), 00).Format("20060102")
		today := time.Now().Format("20060102")
		if day == today {
			// 今天的小时+分
			return time.Unix(cTime.Unix(), 00).Format("15:04")
		}
		yesterday := time.Unix(cTime.Unix(), 00).Format("20060102")
		if day == yesterday {
			return "昨天"
		}
	} else if duration < 86400*365 {
		return time.Unix(cTime.Unix(), 00).Format("01月02日")
	}

	return time.Unix(cTime.Unix(), 00).Format("2006.01.02")
}
