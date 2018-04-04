package main

import (
	"os"
	"axiom"
	"log"
)

func run() int {
	robot, err := axiom.New()
	if err != nil {
		log.Print(err)
		return 1
	}

	tableFlipHandler := &axiom.Handler{
		Pattern: `tableflip`,
		Run: func(res *axiom.Response) error {
			return res.Send(`(╯°□°）╯︵ ┻━┻`)
		},
	}

	robot.Handle(
		tableFlipHandler,

		// Or use a axiom.Handler structure complete with usage...
		&axiom.Handler{
			Pattern: `SYN`,
			Usage:   `axiom syn - replies with "ACK"`,
			Run: func(res *axiom.Response) error {
				return res.Reply("ACK")
			},
		},
	)

	robot.Run()

	return 0
}

func main() {
	os.Exit(run())
}
