package axiom

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultRobotName = `Axiom`
)

type Robot struct {
	name     string
	alias    string
	provider Provider
	store    Store

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

	default_store, err := NewStore(robot)

	if err != nil {
		return nil, err
	}
	robot.store = default_store

	robot.users = NewUserMap(robot)

	return robot, nil
}

func (robot *Robot) SetAlias(alias string) {
	robot.alias = alias
}

func (robot *Robot) SetName(name string) {
	robot.name = name
}

func (robot *Robot) GetName() string {
	return robot.name
}

func (robot *Robot) SetProvider(p Provider) {
	robot.provider = p
}

func (robot *Robot) Provider() Provider {
	return robot.provider
}

func (robot *Robot) SetStore(store Store) {
	robot.store = store
}

func (robot *Robot) Handlers() []handler {
	return robot.handlers
}

func (robot *Robot) Receive(msg *Message) error {

	user := msg.FromUser
	if _, err := robot.users.Get(user.ID); err != nil {
		log.Printf("get user error: %v", err)
		robot.users.Set(user.ID, user)
		robot.users.Save()
	}

	for _, handler := range robot.handlers {
		response := NewResponse(robot, msg)

		if err := handler.Handle(response); err != nil {
			return err
		}
	}
	return nil
}

func (robot *Robot) Handle(handlers ...handler) {
	for _, h := range handlers {

		robot.handlers = append(robot.handlers, h)
	}
}

func (robot *Robot) Run() {
	log.Printf("starting robot")

	log.Printf("opening %s store connection", robot.store.Name())
	go func() {
		robot.store.Open()

		log.Printf("loading users from store")
		robot.users.Load()
	}()

	log.Printf("starting %s provider", robot.provider.Name())
	go robot.provider.Run()

	signal.Notify(robot.signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	stop := false
	for !stop {
		select {
		case sig := <-robot.signalChan:
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				stop = true
			}
		}
	}

	signal.Stop(robot.signalChan)

	robot.Stop()
}

func (robot *Robot) Stop() error {

	log.Printf("stopping %s provider", robot.provider.Name())
	if err := robot.provider.Close(); err != nil {
		return err
	}

	log.Printf("closing %s store connection", robot.store.Name())
	if err := robot.store.Close(); err != nil {
		return err
	}

	log.Printf("stopping robot")
	return nil
}
