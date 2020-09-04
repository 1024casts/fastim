package model

import "sync"

type MsgModel struct {
	BaseModel
	ChatID      uint64 `gorm:"column:chat_id" json:"chat_id"`
	Content     string `gorm:"column:content" json:"content"`
	Extra       string `gorm:"column:extra" json:"extra"`
	LocalMid    uint64 `gorm:"column:local_mid" json:"local_mid"`
	MsgType     int    `gorm:"column:msg_type" json:"msg_type"`
	ReceiveType int    `gorm:"column:receive_type" json:"receive_type"`
	UserID      uint64 `gorm:"column:user_id" json:"user_id"`
}

// TableName sets the insert table name for this struct type
func (m *MsgModel) TableName() string {
	return "im_msg"
}

type MsgInfo struct {
	MsgId    uint64      `json:"msg_id"`
	IsSelf   int         `json:"is_self"`
	MsgType  int         `json:"msg_type"`
	Content  interface{} `json:"content"`
	LocalMid uint64      `json:"local_mid"`
	ShowTime string      `json:"show_time"`
}

type MsgList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*MsgInfo
}

type SendMsgInput struct {
	UserId      uint64 `json:"user_id"`
	YUserId     uint64 `json:"y_user_id"`
	ChatId      uint64 `json:"chat_id"`
	LocalMId    uint64 `json:"local_m_id"`
	MsgType     int    `json:"msg_type"`
	ReceiveType int    `json:"receive_type"`
	Content     string `json:"content"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Duration    int    `json:"duration"`
	ClassId     uint64 `json:"class_id"`
}
