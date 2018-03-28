package axiom

import "log"

// 解析器
type RuleParser interface {
	Name() string                           // 解析器名称
	Boot(*Robot)                            // 引导启动
	HelpMessage(Robot, string) string       // 帮助信息
	HandleMessage(Robot, Message) []Message // 解析消息
}

// 注册解析器
func RegisterRuleset(rule RuleParser) ListenerFunc {
	return func(bot *Robot) {
		log.Printf("bot: registering ruleset %T", rule)
		rule.Boot(bot)
		bot.rules = append(bot.rules, rule)
	}
}