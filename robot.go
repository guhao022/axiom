package axiom

import "os"

const (
	DefaultRobotName = `Axiom`
)

type Robot struct {
	name       string
	provider   Provider
	store      Store
	handlers   []handler
	users      *UserMap
	signalChan chan os.Signal
}

func New() (*Robot, error) {
	robot := &Robot{
		name:       DefaultRobotName,
		signalChan: make(chan os.Signal, 1),
	}

	default_adp, err := robot.newAdapter()

	if err != nil {
		return nil, err
	}

	robot.Adapter = default_adp

	return robot, nil
}

func (robot *Robot) Receive(msg *Message) error {

	// check if we've seen this user yet, and add if we haven't.
	user := msg.User
	if _, err := robot.Users.Get(user.ID); err != nil {
		log.Printf("get user error: %v", err)
		robot.Users.Set(user.ID, user)
		robot.Users.Save()
	}

	for _, handler := range robot.handlers {
		response := NewResponseFromMessage(robot, msg)

		if err := handler.Handle(response); err != nil {
			return err
		}
	}
	return nil
}
