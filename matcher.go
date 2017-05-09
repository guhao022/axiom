package axiom

/*type Matcher interface {
	AddHandler(l *Listener) error    // 添加处理程序
	HandleMessage(msg Message) error // 处理消息
}*/

func NewMatcher(bot *Robot) *Matcher {
	return &Matcher{
		Bot: bot,
	}
}

type Matcher struct {
	Bot      *Robot
	handlers []*Listener
}

func (m *Matcher) AddHandler(l *Listener) error {
	m.handlers = append(m.handlers, l)
	return nil
}

func (m *Matcher) HandleMessage(message Message) error {
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
