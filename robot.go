package axiom

import (
	"log"
	"sync"
	"regexp"
)

const DefaultRobotName string = `Axiom`

// robot结构体，包含所有的内部相关数据
type robot struct {
	name        string         // robot 名称，如果没有指定，使用默认名称
	adapters    []Adapter      // 适配器，可注册多个
	ignoreUsers []string       // 用户黑名单
	preRegex    *regexp.Regexp // 匹配前缀，是否属于命令
	postRegex   *regexp.Regexp // 匹配内容
	mu          sync.RWMutex
	logger      *log.Logger
}

// 实例化一个robot实例
func newRobot() *robot {
	bot := new(robot)

	bot.name = DefaultRobotName

	return bot
}

func (bot *robot) SetName(name string) {
	bot.name = name
}
