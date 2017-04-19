package axiom

import (
)

type Message struct {
	ID      int64      `json:"id"`       // 消息ID
	Text    string     `json:"text"`     // 消息内容
	Adapter []*Adapter `json:"adapter"`  // 适配平台
	ReplyTo []*User    `json:"reply_to"` // 输出适配器
}

/*
func NewMsg(historyID int64, msg string) *Message {
	worker := NewID(historyID)
	id := worker.Next()

	return &Message{
		ID:      id,
		Text:    msg,
		Adapter: nil,
		ReplyTo: nil,
	}
}
*/
