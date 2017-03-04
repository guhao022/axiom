package axiom

import "errors"

type Robot struct {
	name     string
	adapter  Adapter
	monitors []Monitor
	matcher  *Matcher
}

func New(arguments ...string) *Robot {

	b := new(Robot)
	b.monitors = []Monitor{}
	b.matcher = &Matcher{Bot: b}

	if len(arguments) > 0 {
		b.name = arguments[0]
	} else {
		b.name = "Axiom"
	}

	return b
}

// AddAdapter 向Robot中添加适配器
func (b *Robot) AddAdapter(a Adapter) {
	b.adapter = a
}

// Start
func (b *Robot) Start() error {

	if b.adapter == nil {
		return errors.New("You must add at least one adapter")
	}

	err := b.adapter.Construct()

	if err != nil {
		return err
	}

	return b.adapter.Process()
}

// ListenFunc 添加自定义ListenerFunc
func (b *Robot) ListenFunc(regex string, handler ListenerFunc) {
	b.matcher.AddHandler(&Listener{regex, handler})
}

// Register 为Robot注册处理程序
func (b *Robot) Register(h Monitor) {
	handlers := h.Handle()
	for _, handler := range handlers {
		b.matcher.AddHandler(&Listener{handler.Regex, handler.HandlerFunc})
	}
}

// ReceiveMessage 将适配器接收的消息传递给Handler
func (b *Robot) ReceiveMessage(message Message) {
	b.matcher.HandleMessage(message)
}

// Reply 通过适配器回复信息
func (b *Robot) Reply(m Message, message string) {
	b.adapter.Reply(m, message)
}