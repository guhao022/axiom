package axiom

import (
	"fmt"
)

type Context struct {
	Matches []string
	Message Message
	Bot     *Robot
}

func (ctx *Context) Reply(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	ctx.Bot.Reply(ctx.Message, message)
}
