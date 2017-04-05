package main

import (
	"./listener"
	"axiom"
)

func main() {
	b := axiom.New("Axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Start()
}
