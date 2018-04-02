package axiom

func New() (*Robot, error) {
	return NewRobot()
}

// Hear a message
func Hear(pattern string, fn func(res *Response) error) handler {
	return &Handler{Method: HEAR, Pattern: pattern, Run: fn}
}

// Respond creates a new listener for Respond messages
func Respond(pattern string, fn func(res *Response) error) handler {
	return &Handler{Method: RESPOND, Pattern: pattern, Run: fn}
}

// Topic returns a new listener for Topic messages
func Topic(pattern string, fn func(res *Response) error) handler {
	return &Handler{Method: TOPIC, Run: fn}
}

// Enter returns a new listener for Enter messages
func Enter(fn func(res *Response) error) handler {
	return &Handler{Method: ENTER, Run: fn}
}

// Leave creates a new listener for Leave messages
func Leave(fn func(res *Response) error) handler {
	return &Handler{Method: LEAVE, Run: fn}
}

// Close shuts down the robot. Unused?
func Close() error {
	return nil
}
