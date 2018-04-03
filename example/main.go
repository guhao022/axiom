package main

import (
	"os"
	"log"

	"axiom"
	_ "axiom/adapter/cli"
	_ "axiom/store/memory"
)

var pingHandler = axiom.Hear(`ping`, func(res *axiom.Response) error {
	return res.Send("PONG")
})

var echoHandler = axiom.Respond(`echo (.+)`, func(res *axiom.Response) error {
	return res.Reply(res.Match[1])
})

func run() int {
	robot, err := axiom.NewRobot()
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	robot.Handle(
		pingHandler,
		echoHandler,
	)

	if err := robot.Run(); err != nil {
		log.Printf("%s", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
