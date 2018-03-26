package axiom

type Event interface {
	Handle() []*Listener
}
