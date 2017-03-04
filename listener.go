package axiom

import "regexp"

// 监听者
type Monitor interface {
	Handle() []*Listener
}

type Listener struct {
	Regex       string
	HandlerFunc ListenerFunc
}

type ListenerFunc func(c *Context)

type Matcher struct {
	Bot      *Robot
	handlers []*Listener
}

func (m *Matcher) AddHandler(l *Listener) {
	m.handlers = append(m.handlers, l)
}

func (m *Matcher) HandleMessage(message Message) {

	for _, h := range m.handlers {
		regexp, err := regexp.Compile(h.Regex)

		if err != nil {
			panic("regexp err: " + err.Error())
		}

		matches := regexp.FindStringSubmatch(message.Text)

		if len(matches) > 0 {
			c := &Context{
				Bot:     m.Bot,
				Matches: matches,
				Message: message,
			}
			h.HandlerFunc(c)
		}
	}
}
