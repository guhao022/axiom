package main

import (
	"axiom"
	"github.com/num5/axiom/example/shell/listener"
)

func main() {
	b := axiom.New("Axiom")

	b.AddAdapter(axiom.NewShell(b))

	b.Register(&listener.TimeListener{})

	b.Run()
}
