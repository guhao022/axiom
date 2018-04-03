package axiom

// Response struct
type Response struct {
	Robot    *Robot
	Envelope *Envelope
	Message  *Message
	Match    []string
}

// Envelope contains metadata about the chat message.
type Envelope struct {
	Room    string
	User    *User
	Options map[string]interface{}
}

// Options type
type Options map[string]interface{}

// SetOptions sets the Envelope's Options
func (e *Envelope) SetOptions(opts map[string]interface{}) {
	e.Options = opts
}

// NewResponseFromMessage returns a new Response object with an associated Message
func NewResponseFromMessage(robot *Robot, msg *Message) *Response {
	return &Response{
		Robot: robot,
		Envelope: &Envelope{
			Room: msg.Room,
			User: &msg.User,
		},
		Message: msg,
	}
}

func (res *Response) Text() string {
	return res.Message.Text
}
