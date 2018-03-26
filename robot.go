package axiom

import (
	"log"
	"sync"
	"fmt"
)

const DefaultRobotName string = `Axiom`

// 适配器
var adapters map[string]Adapter

// 注册适配器
func RegisterAdpater(name string,adp Adapter) {
	adapters[name] = adp
}

// robot结构体，包含所有的内部相关数据
type robot struct {
	name string // robot 名称，如果没有指定，使用默认名称
	ignoreUsers []User // 用户黑名单
	mu          sync.RWMutex
	logger      *log.Logger
}

// 实例化一个robot实例
func NewRobot() *robot {
	bot := new(robot)

	bot.name = DefaultRobotName

	bot.ignoreUsers = make([]User, 0)

	return bot
}

// 获取名称
func (bot *robot) GetName() string {
	return bot.name
}

// 设置名称
func (bot *robot) SetName(name string) *robot {
	bot.name = name

	return bot
}

// 屏蔽用户
func (bot *robot) IgnoreUsers(user ...User) *robot {
	if len(user) > 0 {
		bot.ignoreUsers = append(bot.ignoreUsers, user...)
	}

	return bot
}

// 接收消息并传送给时间处理器
func (bot *robot) Receive(msg *Message) {
	//
}

// 运行
func (bot *robot) Start() error {
	if len(adapters) <= 0 {
		return fmt.Errorf("You must add at least one adapter")
	}

	for name, adp := range adapters {
		err := adp.Process()

		if err != nil {
			return fmt.Errorf("The %s adapter process error, error is: %v", name, err)
		}
	}

	return nil
}

// 发送消息
func (bot *robot) Send() {

}
