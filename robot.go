package axiom

import (
	"os"
	"os/signal"
	"syscall"
)

const (
	DefaultRobotName = `Axiom`
	DefaultRobotAlias = ``

	HEAR    = `HEAR`
	RESPOND = `RESPOND`
	TOPIC   = `TOPIC`
	ENTER   = `ENTER`
	LEAVE   = `LEAVE`
)

// Robot receives messages from an adapter and sends them to listeners
type Robot struct {
	Name       string
	Alias      string
	Adapter    Adapter
	Store      Store
	handlers   []handler
	Users      *UserMap
	//Auth       *Auth
	signalChan chan os.Signal
}

// Handlers returns the robot's handlers
func (robot *Robot) Handlers() []handler {
	return robot.handlers
}

// NewRobot returns a new Robot instance
func NewRobot() (*Robot, error) {
	robot := &Robot{
		Name:       DefaultRobotName,
		signalChan: make(chan os.Signal, 1),
	}

	adapter, err := NewAdapter(robot)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	robot.SetAdapter(adapter)

	store, err := NewStore(robot)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	robot.SetStore(store)

	robot.Users = NewUserMap(robot)
	//robot.Auth = NewAuth(robot)

	return robot, nil
}

// Handle registers a new handler with the robot
func (robot *Robot) Handle(handlers ...interface{}) {
	for _, h := range handlers {
		nh, err := NewHandler(h)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		robot.handlers = append(robot.handlers, nh)
	}
}

// Receive dispatches messages to our handlers
func (robot *Robot) Receive(msg *Message) error {
	log.Debugf("%s - robot received message", robot.Name)

	// check if we've seen this user yet, and add if we haven't.
	user := msg.User
	if _, err := robot.Users.Get(user.ID); err != nil {
		log.Debug(err)
		robot.Users.Set(user.ID, user)
		robot.Users.Save()
	}

	for _, handler := range robot.handlers {
		response := NewResponseFromMessage(robot, msg)

		if err := handler.Handle(response); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// Run initiates the startup process
func (robot *Robot) Run() error {
	log.Info("starting robot")

	// HACK
	log.Debugf("opening %s store connection", robot.Store.Name())
	go func() {
		robot.Store.Open()

		log.Debug("loading users from store")
		robot.Users.Load()
	}()

	log.Debugf("starting %s adapter", robot.Adapter.Name())
	go robot.Adapter.Run()

	// Start the HTTP server after the adapter, as adapter.Run() adds additional
	// handlers to the router.
	/*log.Debug("starting HTTP server")
	go func() {
		if err := http.ListenAndServe(`:9900`, Router); err != nil {
			log.Debug(err)
		}
	}()*/

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
	// Stop listening for new signals
	signal.Stop(robot.signalChan)

	// Initiate the shutdown process for our robot
	robot.Stop()

	return nil
}

// Stop initiates the shutdown process
func (robot *Robot) Stop() error {
	log.Info() // so we don't break up the log formatting when running interactively ;)

	log.Debugf("stopping %s adapter", robot.Adapter.Name())
	if err := robot.Adapter.Stop(); err != nil {
		return err
	}

	log.Debugf("closing %s store connection", robot.Store.Name())
	if err := robot.Store.Close(); err != nil {
		return err
	}

	log.Info("stopping robot")
	return nil
}

// SetName sets robot's name
func (robot *Robot) SetName(name string) {
	robot.Name = name
}

// SetAdapter sets robot's adapter
func (robot *Robot) SetAdapter(adapter Adapter) {
	robot.Adapter = adapter
}

// SetStore sets robot's adapter
func (robot *Robot) SetStore(store Store) {
	robot.Store = store
}
