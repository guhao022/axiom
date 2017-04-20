package main

import (
	"github.com/num5/axiom/example/shell/listener"
	"axiom"
)

func main() {
	b := axiom.New("Axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Run()
}
