package axiom

// 解析器
type Matcher interface {
	Name() string                           // 解析器名称
	Boot(*Robot)                            // 引导启动
	HelpMessage(Robot, string) string       // 帮助信息
	HandleMessage(Robot, Message) []Message // 解析消息
}
