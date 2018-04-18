package wechat

import (
	"github.com/num5/axiom"
	"fmt"
)

type axiomMsg struct {
	*axiom.Message
	to string
}

func NewAxiomMsg(msg *axiom.Message, to string) *axiomMsg {
	return &axiomMsg{
		msg,
		to,
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