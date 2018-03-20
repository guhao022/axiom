package v1_01

// 监听者
type ListenEvent interface {
	Handle() []*Listener
}

type Listener struct {
	Regex       string
	HandlerFunc ListenerFunc
}

type ListenerFunc func(c *Context)
