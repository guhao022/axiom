package axiom

type Context struct {
	Matches []string
	Message Message
	Bot     *Robot
}

func (ctx *Context) Reply(message string) {
	ctx.Bot.Reply(ctx.Message, message)
}