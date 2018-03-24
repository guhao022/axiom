package axiom

const (
	// 私聊信息，用户和机器人单独聊天的信息
	DirectedMessage = "directed-message"

	// 频道信息，主要用于频道广播
	ChannelMessage  = "channel-message"
)

type Message struct {
	msgType string
	msg       string
	err       error
}
