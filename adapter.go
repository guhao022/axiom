package axiom

const (
	DefaultAdapter = `shell`
)

type Adapter interface {
	IncomingChannel() chan Message
	OutgoingChannel() chan Message
	Name() string
}

var availableAdapters []Adapter




