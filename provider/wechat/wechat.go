package wechat

import (
	"github.com/num5/axiom"
	"github.com/KevinGong2013/wechat"
)

type weixin struct {
	axiom.BasicProvider
	bot *axiom.Robot
	wechat *wechat.WeChat
}

func newWeChat(r *axiom.Robot) *weixin {
	wx := new(weixin)

	wx.SetRobot(r)

	wc, err := wechat.NewBot(nil)
	if err != nil {
		panic(err)
	}

	wx.wechat = wc

	return wx
}

func (wx *weixin) Name() string {
	return "web wechat"
}

func (wx *weixin) Run() error {
	return nil
}

func (wx *weixin) Close() error {
	return nil
}

func (wx *weixin) Receive(*axiom.Message) error {
	return nil
}

func (wx *weixin) Send(*axiom.Response, ...string) error {
	return nil
}