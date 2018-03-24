package axiom

// 适配器
type Adapter interface {
	// 获取连接器名称
	GetName() string

	// 处理消息
	Process() error

	// 发送消息
	Send() error

	// 回答某个人的消息
	Reply() error
}

