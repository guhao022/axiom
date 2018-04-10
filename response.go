package axiom

// Response struct
type Response struct {
	Robot   *Robot
	Message *Message
	Match   []string
}

func NewResponse(robot *Robot, msg *Message) *Response {
	return &Response{
		Robot:   robot,
		Message: msg,
	}
}

func (res *Response) Text() string {

	return res.Message.Text
}

func (res *Response) Room() string {
	return res.Message.Room
}

func (res *Response) UserID() string {
	return res.Message.User.ID
}

func (res *Response) UserName() string {
	return res.Message.User.Name
}

func (res *Response) Send(strings ...string) error {
	if err := res.Robot.Provider().Send(res, strings...); err != nil {
		return err
	}
	return nil
}

func (res *Response) Reply(strings ...string) error {
	if err := res.Robot.Provider().Reply(res, strings...); err != nil {
		return err
	}
	return nil
}
