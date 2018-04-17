package axiom

// 用于实现机器人工作的适配器
type Adapter interface {
	Construct() error            // 初始化
	Process() error              // 处理
	Reply(Message, string) error // 回复
}
