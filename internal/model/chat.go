package model

import (
	"sync"
)

type ChatModel struct {
	BaseModel
	SenderUID   uint64 `gorm:"column:sender_uid" json:"sender_uid"`
	ReceiverUID uint64 `gorm:"column:receiver_uid" json:"receiver_uid"`
	LastMsgID   int    `gorm:"column:last_msg_id" json:"last_msg_id"`
	MsgNum      int    `gorm:"column:msg_num" json:"msg_num"`
	IsDelete    int    `gorm:"column:is_delete" json:"is_delete"`
	Extra       string `gorm:"column:extra" json:"extra"`
}

// TableName sets the insert table name for this struct type
func (c *ChatModel) TableName() string {
	return "im_chat"
}

type ChatInfo struct {
	//User      *UserInfo `json:"user"`
	Msg       *MsgInfo `json:"msg"`
	ShowTime  string   `json:"show_time"`
	NewMsgNum int      `json:"new_msg_num"`
}

type UserChatList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*ChatInfo
}
