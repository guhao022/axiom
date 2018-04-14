package irc

import (
	"axiom"
	irc "github.com/thoj/go-ircevent"
)

type ircprovider struct {
	axiom.BasicProvider

	conn     *irc.Connection
}

func NewIRC(r *axiom.Robot) (axiom.Provider, error) {
	irc := &ircprovider{}
	irc.SetRobot(r)
	irc.Robot.SetName("irc")
	return irc, nil
}

func (i *ircprovider) Name() string {
	return `irc`
}

func (i *ircprovider) Send(res *axiom.Response, strings ...string) error {
	for _, str := range strings {
		s := &ircPayload{
			Channel: res.Message.Room,
			Text:    str,
		}
		i.conn.Privmsg(s.Channel, s.Text)
	}
	return nil
}

func (i *ircprovider) Reply(res *axiom.Response, strings ...string) error {

	return nil
}

func (i *ircprovider) Run() error {

	return nil
}

func (i *ircprovider) Close() error {

	return nil
}
