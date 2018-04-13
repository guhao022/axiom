package webserver

import "github.com/num5/axiom"

type webserver struct {
}

func (w *webserver) Name() string {
	return "web server"
}

func (w *webserver) Run() error {
	return  nil
}

func (w *webserver) Close() error {
	return  nil
}

func (w *webserver) Receive(*axiom.Message) error {
	return  nil
}

func (w *webserver) Send(*axiom.Response, ...string) error {
	return  nil
}

func (w *webserver) Reply(*axiom.Response, ...string) error {
	return  nil
}