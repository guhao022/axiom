package axiom

import (
	"net/http"
	"time"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"regexp"
	"os"
	"io"
	"os/exec"
	"runtime"
	"encoding/xml"
	"strconv"
)

const (
	DefaultUserAgent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
)

type WeiChat struct {
	bot *Robot
	wx *weixin
}

type weixin struct {
	qr_code_path string
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
	BaseRequest  BaseRequest
	syncHost     string

	http_client  *http.Client
}

type BaseRequest struct {
	//
	Uin int
	Sid string

}

func NewWeiChat(bot *Robot, qr_code_path string) *WeiChat {
	wx := &weixin{
		qr_code_path: qr_code_path,
	}
	return &WeiChat{
		bot: bot,
		wx: wx,
	}
}

// 第一步获取UUID
func (w *WeiChat) getUUID() bool {
	urlstr := "https://login.weixin.qq.com/jslogin"

	params := map[string]interface{} {
		"appid": "wx782c26e4c19acffb",
		"fun": "new",
		"lang": "zh_CN",
		"_": time.Now().Unix(),

	}

	data, err := w.httpPost(urlstr, params)
	if err != nil {
		return false
	}
	reg := regexp.MustCompile(`window.QRLogin.code = ([\d]+); window.QRLogin.uuid = "([\S]+)"`)
	req := reg.FindStringSubmatch(data)
	code := req[1]
	if code == "200" {
		w.wx.uuid = req[2]
		return true
	}
	return false
}

// 第二步获取二维码
func (w *WeiChat) getQRCode() bool {
	url := "https://login.weixin.qq.com/qrcode/" + w.wx.uuid
	params := map[string]interface{} {
		"t": "webwx",
		"_": time.Now().Unix(),
	}

	path := w.wx.qr_code_path + "/qrcode.jpg"
	out, err := os.Create(path)

	resp, err := w.httpPost(url, params)
	if err != nil {
		return false
	}
	_, err = io.Copy(out, bytes.NewReader(resp))
	if err != nil {
		return false
	} else {
		if runtime.GOOS == "darwin" {
			exec.Command("open", path).Run()
		} else {
			go func() {
				fmt.Println("please open on web broswer ip:99250/qrcode")
				http.HandleFunc("/qrcode", func(w http.ResponseWriter, req *http.Request) {
					http.ServeFile(w, req, "qrcode.jpg")
					return
				})
				http.ListenAndServe(":99250", nil)
			}()
		}
		return true
	}
}

// 第三步， 等待登录
func (w *WeiChat) waitForLogin(tip int) bool {
	time.Sleep(time.Duration(tip) * time.Second)

	url := "https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login"

	params := map[string]interface{} {
		"tip": tip,
		"uuid": w.wx.uuid,
		"_": time.Now().Unix(),
	}

	data, _ := w.httpGet(url, params)
	reg := regexp.MustCompile(`window.code=(\d+);`)
	req := reg.FindStringSubmatch(data)

	if len(req) > 1 {
		code := req[1]
		if code == "201" {
			return true

		} else if code == "200" {
			u_reg := regexp.MustCompile(`window.redirect_uri="(\S+?)";`)
			u_req := u_reg.FindStringSubmatch(data)

			if len(u_req) > 1 {
				r_uri := u_req[1] + "&fun=new"
				w.wx.redirect_uri = r_uri

				bu_reg := regexp.MustCompile(`/`)
				bu_req := bu_reg.FindAllStringIndex(r_uri, -1)

				w.wx.base_uri = r_uri[:bu_req[len(bu_req)-1][0]]

				return true
			}
			return false
		} else if code == "408" {
			fmt.Println("[登陆超时]")
		} else {
			fmt.Println("[登陆异常]")
		}
	}

	return false
}

// 第四步，登录获取Cookie
func (w *WeiChat) login(args ...interface{}) bool {
	data, _ := w.httpGet(w.wx.redirect_uri, "")

	type result struct {
		Skey        string `xml:"skey"`
		Wxsid       string `xml:"wxsid"`
		Wxuin       string `xml:"wxuin"`
		Pass_ticket string `xml:"pass_ticket"`
	}
	v := result{}
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false
	}
	w.wx.skey = v.Skey
	w.wx.sid = v.Wxsid
	w.wx.uin = v.Wxuin
	w.wx.pass_ticket = v.Pass_ticket
	w.wx.BaseRequest = make(map[string]interface{})
	w.wx.BaseRequest["Uin"], _ = strconv.Atoi(v.Wxuin)
	w.wx.BaseRequest["Sid"] = v.Wxsid
	w.wx.BaseRequest["Skey"] = v.Skey
	w.wx.BaseRequest["DeviceID"] = w.wx.device_id
	return true
}

// get 方法
func (w *WeiChat) httpGet(url string, params map[string]interface{}) ([]byte, error) {
	dataJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	data := bytes.NewBuffer(dataJson)

	request, _ := http.NewRequest("GET", url, data)

	request.Header.Add("Referer", "https://wx.qq.com/")

	request.Header.Add("User-agent", DefaultUserAgent)

	resp, err := w.wx.http_client.Do(request)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// post 方法
func (w *WeiChat) httpPost(url string, params map[string]interface{}) ([]byte, error) {
	postJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	postData := bytes.NewBuffer(postJson)

	request, err := http.NewRequest("POST", url, postData)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json;charset=utf-8")

	request.Header.Add("Referer", "https://wx.qq.com/")

	request.Header.Add("User-agent", DefaultUserAgent)

	resp, err := w.wx.http_client.Do(request)

	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}


// 初始化
func (w *WeiChat) Construct() error {
	return nil
}

// 解析
func (w *WeiChat) Process() error {

	return nil

}

// 回应
func (w *WeiChat) Reply(msg Message, message string) error {
	return nil
}
