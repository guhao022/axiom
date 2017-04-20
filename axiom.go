package axiom

import (
	"errors"
)

type Robot struct {
	name     string
	adapter  []Adapter
	listener []ListenEvent
	matcher  Matcher
}

func New(name ...string) *Robot {

	b := new(Robot)
	b.listener = []ListenEvent{}

	if len(name) > 0 {
		b.name = name[0]
	} else {
		b.name = "Axiom"
	}

	return b
}

func (b *Robot) SetMatcher(m Matcher) {
	b.matcher = m
}

// AddAdapter 向Robot中添加适配器
func (b *Robot) AddAdapter(a Adapter) {
	b.adapter = append(b.adapter, a)
}

// Start
func (b *Robot) Run() error {

	if len(b.adapter) <= 0 {
		return errors.New("You must add at least one adapter")
	}

	for _, adapter := range b.adapter {
		err := adapter.Prepare()

		if err != nil {
			return err
		}

		err = adapter.Process()

		if err != nil {
			return err
		}
	}

	return nil
}

// ListenFunc 添加自定义ListenerFunc
func (b *Robot) ListenFunc(regex string, handler ListenerFunc) {
	b.matcher.AddHandler(&Listener{regex, handler})
}

// Register 为Robot注册处理程序
func (b *Robot) Register(listener ...ListenEvent) {
	if len(listener) <= 0 {
		panic("监听器不能为空")
	}
	for _, l := range listener {
		handlers := l.Handle()
		for _, handler := range handlers {
			b.matcher.AddHandler(&Listener{handler.Regex, handler.HandlerFunc})
		}
	}
}

// ReceiveMessage 将适配器接收的消息传递给Handler
func (b *Robot) ReceiveMessage(message Message) {
	b.matcher.HandleMessage(message)
}

// Reply 通过适配器回复信息
func (b *Robot) Reply(m Message, message string) {
	if m.Adapter == nil {
		for _, adapter := range b.adapter {

			err := adapter.Reply(m, message)

			if err != nil {
				log.Errorf("适配器 [%s] 回复消息失败：%s...", adapter.GetName(), err.Error())
			}
		}
	}

}
