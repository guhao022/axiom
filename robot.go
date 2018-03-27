package axiom

const DefaultRobotName string = `Axiom`

// Robot结构体，包含所有的内部相关数据
type Robot struct {
	name        string
	providerIn  chan Message
	providerOut chan Message
	matchers    []Matcher

	//brain brain.Memorizer
}

