package axiom

import (
	"strings"
	"fmt"
	"sync"
	"log"
)

const DefaultRobotName string = `Axiom`

// Robot结构体，包含所有的内部相关数据
type Robot struct {
	name        string
	providerIn  chan Message
	providerOut chan Message
	rules       []RuleParser

	once sync.Once

	//brain brain.Memorizer
}

var processOnce sync.Once

//
type ListenerFunc func(*Robot)

func NewBot(name string) *Robot {
	bot := &Robot{
		name: name,
		providerIn: make(chan Message),
		providerOut: make(chan Message),
	}

	bot.name = name

	return bot
}

func (bot *Robot) Process() {
	processOnce.Do(func() {

		for in := range bot.providerIn {
			if strings.HasPrefix(in.Message, bot.Name()+" help") {
				go func(robot Robot, msg Message) {
					helpMsg := fmt.Sprintln("available commands:")
					for _, rule := range bot.rules {
						helpMsg = fmt.Sprintln(helpMsg, rule.HelpMessage(robot, in.Room))
					}
					bot.providerOut <- Message{
						Room:       msg.Room,
						ToUserID:   msg.FromUserID,
						ToUserName: msg.FromUserName,
						Message:    helpMsg,
					}
				}(*bot, in)
				continue
			}
			go func(robot Robot, msg Message) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("panic recovered when parsing message: %#v. Panic: %v", msg, r)
					}
				}()
				for _, rule := range bot.rules {
					responses := rule.HandleMessage(robot, msg)
					for _, resp := range responses {
						bot.providerOut <- resp
					}
				}
			}(*bot, in)
		}
	})
}

func (bot *Robot) Name() string {
	return bot.name
}

func (bot *Robot) MessageProviderOut() chan Message {
	return bot.providerOut
}