package main

import (
	"axiom"
	"axiom/example/shell/listener"
)

func main() {
	b := axiom.New("Axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Run()
}
