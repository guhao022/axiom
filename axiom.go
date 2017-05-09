package axiom

import (
	"errors"
	"log"
	"regexp"
)

type Robot struct {
	name     string
	adapter  Adapter
	listener []ListenEvent
	matcher  *Matcher

	done chan bool
}

func New(name ...string) *Robot {

	b := new(Robot)
	b.listener = []ListenEvent{}
	b.matcher = &Matcher{Bot: b}

	if len(name) > 0 {
		b.name = name[0]
	} else {
		b.name = "Axiom"
	}

	return b
}

// AddAdapter 向Robot中添加适配器
func (b *Robot) AddAdapter(a Adapter) {
	b.adapter = a
}

// Start，
func (b *Robot) Run() error {

	if b.adapter == nil {
		return errors.New("You must add at least one adapter")
	}

	err := b.adapter.Prepare()

	if err != nil {
		log.Printf("[%s] 适配器初始化错误：%v ", b.adapter.GetName(), err)
	}

	err = b.adapter.Process()

	if err != nil {
		log.Printf("[%s] 适配器处理错误：%v ", b.adapter.GetName(), err)
	}

	/*go func() {


	}()

	<- b.done*/

	return nil
}

// ListenFunc 添加自定义ListenerFunc
func (b *Robot) ListenFunc(regex string, handler ListenerFunc) error {
	regexp := regexp.MustCompile(regex)
	return b.matcher.AddHandler(&Listener{regexp, handler})
}

// Register 为Robot注册处理程序
func (b *Robot) Register(listener ...ListenEvent) error {

	if b.matcher == nil {
		return errors.New("没有设置消息处理器 [Matcher]")
	}

	if len(listener) <= 0 {
		return errors.New("监听器不能为空")
	}
	for _, l := range listener {
		handlers := l.Handle()
		for _, handler := range handlers {
			regexp := regexp.MustCompile(handler.Regex)
			return b.matcher.AddHandler(&Listener{regexp, handler.HandlerFunc})
		}
	}
	return nil
}

// ReceiveMessage 将适配器接收的消息传递给Handler
func (b *Robot) ReceiveMessage(message Message) error {
	return b.matcher.HandleMessage(message)
}

// Reply 通过适配器回复信息
func (b *Robot) Reply(m Message, message string) {

	err := b.adapter.Reply(m, message)

	if err != nil {
		log.Printf("适配器 [%s] 回复消息失败：%s...", b.adapter.GetName(), err.Error())
	}

}

func (b *Robot) Stop() {
	b.done <- true
}
