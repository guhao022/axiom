package axiom

type Provider interface {
	IncomingChannel() chan Message
	OutgoingChannel() chan Message
	Error() error
}
