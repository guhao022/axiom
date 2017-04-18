package axiom

// 监听者
type Listener interface {
	Handle() []*Listener
}



