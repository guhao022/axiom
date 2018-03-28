package axiom

import "log"

// 保存机器人发送、接收消息所有的元数据
type Message struct {
	Room         string
	FromUserID   string
	FromUserName string
	ToUserID     string
	ToUserName   string
	Message      string
	Direct       bool
}

func MessageProvider(provider Provider) ListenerFunc {
	return func(bot *Robot) {
		log.Printf("bot: changing message provider %T\n", provider)
		bot.providerIn = provider.IncomingChannel()
		bot.providerOut = provider.OutgoingChannel()
	}
}
