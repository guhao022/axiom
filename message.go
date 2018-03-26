package axiom

import "regexp"

const (
	// 私聊信息，用户和机器人单独聊天的信息
	DirectedMessage = "directed-message"

	// 频道信息，主要用于频道广播
	ChannelMessage = "channel-message"
)

type Message struct {
	user      *User          // 消息发送者
	msg       string         // 消息内容
	msgType   string         // 消息类型
	preRegex  *regexp.Regexp // 匹配前缀，是否属于命令
	postRegex *regexp.Regexp // 匹配内容
	err       error
}
