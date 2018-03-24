package axiom

import (
	"github.com/num5/logger"
	"net/http"
	"sync"
	"time"
	"crypto/tls"
	"net"
)

const DefaultRobotName string = `Axiom`

// robot结构体，包含所有的内部相关数据
type Robot struct {
	name	string
	adapters     []Adapter
	//admins	[]Users
	tokens     []string
	httpClient *http.Client

	// 设置是否验证 TLS，true表示需要验证证书，false不验证
	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to man-in-the-middle attacks.
	// This should be used only for testing.
	insecureSkipVerify bool
	//store *Store
	//scripts      *Scripts
	mutex  *sync.Mutex
	logger *logger.Log
}

// 实例化一个robot实例
func NewRobot(name ...string) *Robot {
	bot := new(Robot)

	if len(name) > 0 {
		bot.name = name[0]
	} else {
		bot.name = DefaultRobotName
	}

	bot.httpClient = &http.Client{}

	return bot
}

func (bot *Robot) createHttpClient() {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: bot.insecureSkipVerify},
	}
	bot.httpClient.Transport = tr
}