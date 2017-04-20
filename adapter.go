// 工作平台
package axiom

// 用于实现机器人工作的适配器
type Adapter interface {
	Prepare() error              // 初始化
	GetName() string             // 适配器名称
	Process() error              // 处理
	Reply(Message, string) error // 回复
}
