package axiom

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


