package webhook

import (
	"github.com/num5/axiom"
	"net/http"
)

type webhook struct {
	axiom.BasicProvider
}

func NewWeb(r *axiom.Robot) (axiom.Provider, error) {
	wh := &webhook{}
	wh.SetRobot(r)
	wh.Robot.SetName("webhook")
	return wh, nil
}

func (wh *webhook) Name() string {
	return `webhook`
}

func (wh *webhook) Send(res *axiom.Response, strings ...string) error {

	return nil
}

func (wh *webhook) Reply(res *axiom.Response, strings ...string) error {

	return nil
}

func (wh *webhook) Run() error {

	return nil
}

func (wh *webhook) Close() error {

	return nil
}

func (wh *webhook) Handler(w http.ResponseWriter, r *http.Request) {

}

