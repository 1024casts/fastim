package im

type SendRequest struct {
	UserId   uint64 `json:"user_id"`
	MsgType  int    `json:"msg_type"`
	Content  string `json:"content"`
	LocalMId uint64 `json:"local_m_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Duration int    `json:"duration"`
}
