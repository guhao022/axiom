package axiom

type Message struct {
	ID      int64      `json:"id"`                      // 消息ID
	Text    string     `json:"text"`                    // 消息内容
	Adapter []*Adapter `json:"adapter"`                 // 适配平台
	ReplyTo []*User    `json:"reply_to"`                // 输出适配器
	History *history   `json:"group_message,omitempty"` // 消息组
}

type history struct {
	ID      int64     `json:"id,omitempty"`          // 子消息ID
	Message []Message `json:"sub_message,omitempty"` // 消息内容
}


