package axiom

type Matcher interface {
	AddHandler(l *Listener) error    // 添加处理程序
	HandleMessage(msg Message) error // 处理消息
}

func NewMatcher(bot *Robot) Matcher {
	return &matcher{
		Bot: bot,
	}
}

type matcher struct {
	Bot      *Robot
	handlers []*Listener
}

func (m *matcher) AddHandler(l *Listener) error {
	m.handlers = append(m.handlers, l)
	return nil
}

func (m *matcher) HandleMessage(message Message) error {
	for _, h := range m.handlers {

		matches := h.Regexp.FindStringSubmatch(message.Text)

		if len(matches) > 0 {
			c := &Context{
				Bot:     m.Bot,
				Matches: matches,
				Message: message,
			}
			h.HandlerFunc(c)

		}
		//continue
	}

	return nil
}
