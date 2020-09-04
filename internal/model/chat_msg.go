package model

type ChatMsgModel struct {
	ID     uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	ChatID uint64 `gorm:"column:chat_id" json:"chat_id"`
	MsgID  uint64 `gorm:"column:msg_id" json:"msg_id"`
}

// TableName sets the insert table name for this struct type
func (c *ChatMsgModel) TableName() string {
	return "im_chat_msg"
}
