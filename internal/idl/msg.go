package idl

import (
	"encoding/json"
	"time"

	"github.com/1024casts/fastim/internal/model"
	"github.com/1024casts/fastim/internal/service"
	"github.com/1024casts/fastim/pkg/utils"
	"github.com/1024casts/snake/pkg/log"
)

type TransMsgInput struct {
	CurUser *model.UserBaseModel
	Msg     *model.MsgModel
}

type textContent struct {
	Text string `json:"text"`
}

type picContent struct {
	Text   string `json:"text"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func TransMsg(input *TransMsgInput) (*model.MsgInfo, error) {
	msg := input.Msg

	contentMap := make(map[string]interface{})
	switch msg.MsgType {
	case service.MsgTypeText:
		content := textContent{}
		if err := json.Unmarshal([]byte(msg.Content), &content); err != nil {
			log.Warnf("[idl] unmarshal json for msg content err, %v", err)
			return nil, err
		}
		contentMap["text"] = content.Text
	case service.MsgTypePic:
		content := picContent{}
		if err := json.Unmarshal([]byte(msg.Content), &content); err != nil {
			log.Warnf("[idl] unmarshal json for msg content err, %v", err)
			return nil, err
		}
		// 图片需要组合字段
		contentMap["pic_url"] = utils.GetStaticImageUrl(content.Text)
		contentMap["width"] = content.Width
		contentMap["height"] = content.Height
	default:
		content := textContent{}
		if err := json.Unmarshal([]byte(msg.Content), &content); err != nil {
			log.Warnf("[idl] unmarshal json for msg content err, %v", err)
			return nil, err
		}
		contentMap["text"] = content.Text
	}

	// 是否是自己
	isSelf := 0
	if input.CurUser.ID == msg.UserID {
		isSelf = 1
	}

	return &model.MsgInfo{
		MsgId:    msg.ID,
		MsgType:  msg.MsgType,
		Content:  contentMap,
		IsSelf:   isSelf,
		LocalMid: msg.LocalMid,
		ShowTime: TransMsgTime(msg.CreatedAt),
	}, nil
}

func TransMsgTime(cTime time.Time) string {
	return time.Unix(cTime.Unix(), 00).Format("01月02日 15:04")
}
