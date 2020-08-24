package model

type UserChatModel struct {
	Id         uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	ChatID     uint64 `gorm:"column:chat_id" json:"chat_id"`
	ClearMsgID uint64 `gorm:"column:clear_msg_id" json:"clear_msg_id"`
	DelMsgID   uint64 `gorm:"column:del_msg_id" json:"del_msg_id"`
	Extra      string `gorm:"column:extra" json:"extra"`
	LastMsgID  uint64 `gorm:"column:last_msg_id" json:"last_msg_id"`
	MsgNum     int    `gorm:"column:msg_num" json:"msg_num"`
	UserID     uint64 `gorm:"column:user_id" json:"user_id"`
	YUserID    uint64 `gorm:"column:yuser_id" json:"yuser_id"`
}

// TableName sets the insert table name for this struct type
func (u *UserChatModel) TableName() string {
	return "im_user_chat"
}
