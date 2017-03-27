package main

import (
	"./listener"
	"axiom"
)

func main() {
	b := axiom.New("axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Start()
}
