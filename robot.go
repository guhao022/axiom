package axiom

import (
	"os"
	"log"
	"os/user"
)

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

	default_provider, err := NewProvider(robot)

	if err != nil {
		return nil, err
	}

	robot.provider = default_provider

	store, err := NewStore(robot)

	if err != nil {
		return nil, err
	}
	robot.store = store
	return robot, nil
}

func (robot *Robot) Receive(msg Message) error {

	user := msg.User
	if _, err := robot.users.Get(user.ID); err != nil {
		log.Printf("get user error: %v", err)
		robot.users.Set(user.ID, user)
		robot.users.Save()
	}

	for _, handler := range robot.handlers {
		response := NewResponseFromMessage(robot, msg)

		if err := handler.Handle(response); err != nil {
			return err
		}
	}
	return nil
}
