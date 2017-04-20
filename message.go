package axiom

type Message struct {
	ID      int64      `json:"id"`      // 消息ID
	Text    string     `json:"text"`    // 消息内容
	Adapter []*Adapter `json:"adapter"` // 适配平台
	History history    `json:"history"` // 消息历史
}

// 生成新的msg
func NewMessage(wid int64, msg string) Message {
	gener := NewID(wid)
	id := gener.Next()

	return Message{
		ID:      id,
		Text:    msg,
		Adapter: nil,
	}
}
