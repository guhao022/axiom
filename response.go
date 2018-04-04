package axiom

// Response struct
type Response struct {
	Robot    *Robot
	Message  *Message
	Match    []string
}

// NewResponseFromMessage returns a new Response object with an associated Message
func NewResponseFromMessage(robot *Robot, msg *Message) *Response {
	return &Response{
		Robot: robot,
		Message: msg,
	}
}

func (res *Response) Text() string {

	return res.Message.Text
}
