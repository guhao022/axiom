package main

import (
	"axiom"
	"./listener"
)

func main() {
	b := axiom.New("axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Start()
}
