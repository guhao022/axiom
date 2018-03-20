package main

import (
	"axiom"
	"axiom/v1.01/example/shell/listener"
)

func main() {
	b := axiom.New("Axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Run()
}
