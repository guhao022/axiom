package axiom

import "net/http"

type WeiChat struct {
	bot *Robot
	wx *weixin
}

type weixin struct {
	// 本次微信登录需要的UUID
	uuid string
	base_uri string
	redirect_uri string
	uin string
	sid string
	skey string
	pass_ticket string
	device_id string
	synckey string
	SyncKey map[string]interface{}
	User map[string]interface{}
	BaseRequest  map[string]interface{}
	syncHost     string
	http_client  *http.Client
}

func NewWeiChat(bot *Robot) *WeiChat {
	wx := &weixin{

	}
	return &WeiChat{
		bot: bot,
	}
}

// 第一步获取UUID
func getUUID() {
	//
}


func (w *WeiChat) Construct() error {
	return nil
}

func (w *WeiChat) Process() error {

	return nil

}

func (w *WeiChat) Reply(msg Message, message string) error {
	return nil
}
