package axiom

type Message struct {
	ID   string
	FromUser User
	ToUser	User
	Room string
	Text string
}
