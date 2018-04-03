package axiom

import "fmt"

// Adapter interface
type Adapter interface {
	Name() string
	Run() error
	Stop() error

	Receive(*Message) error
	Send(*Response, ...string) error
	Emote(*Response, ...string) error
	Reply(*Response, ...string) error
	Topic(*Response, ...string) error
	Play(*Response, ...string) error
}

type adapter struct {
	name     string
	newFunc  func(*Robot) (Adapter, error)
	sendChan chan *Response
	recvChan chan *Message
}

// AvailableAdapters is a map of registered adapters
var AvailableAdapters = map[string]adapter{}

// RegisterAdapter registers an adapter
func RegisterAdapter(name string, newFunc func(*Robot) (Adapter, error)) {
	AvailableAdapters[name] = adapter{
		name:    name,
		newFunc: newFunc,
	}
}

// BasicAdapter declares common functions shared by all adapters
type BasicAdapter struct {
	*Robot
}

// SetRobot sets the adapter's Robot
func (a *BasicAdapter) SetRobot(r *Robot) {
	a.Robot = r
}
