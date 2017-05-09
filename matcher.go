package axiom

import "regexp"

type Matcher interface {
	AddHandler(*Listener) error
	HandleMessage(Message) error
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
		regexp, err := regexp.Compile(h.Regex)

		if err != nil {
			return err
		}
		matches := regexp.FindStringSubmatch(message.Text)

		if len(matches) > 0 {
			c := &Context{
				Bot:     m.Bot,
				Matches: matches,
				Message: message,
			}
			h.HandlerFunc(c)
			return nil
		}

	}

	return nil
}


