package wechat

import (
	"github.com/num5/axiom"
	"fmt"
	"github.com/KevinGong2013/wechat"
)

type axiomMsg struct {
	*axiom.Message
	to string
}

func NewAxiomMsg(wmsg wechat.EventMsgData) *axiomMsg {

	msg := &axiom.Message{
		Text: wmsg.Content,
	}

	return &axiomMsg{
		msg,
		wmsg.FromUserName,
	}
}

func (am *axiomMsg) Path() string {
	return `axiommsg`
}

func (am *axiomMsg) To() string {
	return am.to
}

func (am *axiomMsg) Content() map[string]interface{} {
	content := make(map[string]interface{}, 0)

	content["Type"] = 1
	content["Content"] = am.Text

	return content
}

func (am *axiomMsg) Description() string {
	return fmt.Sprintf(`[AxiomTextMessage] %s`, am.Text)
}

func (am *axiomMsg) String() string {
	return am.Text
}