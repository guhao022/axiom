package axiom

const (
	DefaultAdapter = `shell`
)

type Provider interface {
	IncomingChannel() chan Message
	OutgoingChannel() chan Message
	Name() string
}

var availableProviders []func(func(string) string) (Provider, bool)

func Detect(getenv func(string) string) Provider {
	for _, ap := range availableProviders {
		if ret, ok := ap(getenv); ok {
			if ret.Error() != nil {
				log.Printf("providers: %T %v", ret, ret.Error())
				continue
			}
			return ret
		}
	}
	log.Println("providers: no message provider found.")
	log.Println("providers: falling back to CLI.")
	return CLI()
}


