package axiom

import "regexp"

// 监听者
type ListenEvent interface {
	Handle() []*ListenerHandler
}

type ListenerHandler struct {
	Regex       string
	HandlerFunc ListenerFunc
}

type Listener struct {
	Regexp       *regexp.Regexp
	HandlerFunc ListenerFunc
}

type ListenerFunc func(c *Context)
