package axiom

// 适配器
type Adapter interface {
	GetName() string // 获取连接器名称
	Process(message Message) error  // 处理消息
	Send() error     // 发送消息
	Reply() error    // 回复消息
}
