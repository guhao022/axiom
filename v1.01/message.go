package v1_01

import "github.com/num5/ider"

type Message struct {
	ID      int64         `json:"id"`       // 消息ID
	Text    string        `json:"text"`     // 消息内容
	ReplyTo []interface{} `json:"reply_to"` // 接收方
}

// 生成新的msg
func NewMessage(wid int64, msg string) Message {
	id := ider.NewID(wid).Next()

	m := Message{
		ID:      id,
		Text:    msg,
		ReplyTo: nil,
	}

	return m
}
