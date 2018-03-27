package axiom

import (
	"time"
)

type User struct {
	ID                string   `json:"id"`                  // 用户ID
	Color             string   `json:"color"`               // 输出颜色
	NickName          string   `json:"nick_name"`           // 用户昵称
	RealName          string   `json:"real_name"`           // 用户真实姓名
	Email             string   `json:"email"`               // 用户邮箱
	Phone             string   `json:"phone"`               // 用户手机号
	IsBot             bool     `json:"is_bot"`              // 是否是机器人
	IsAdmin           bool     `json:"is_admin"`            // 是否是管理员
	IsRestricted      bool     `json:"is_restricted"`       // 是否被限制
	IsUltraRestricted bool     `json:"is_ultra_restricted"` // 是否为超级限制，任何功能都不能使用
	Online            string   `json:"online"`              // 是否在线
	Presence          Presence `json:"presence"`            // 线上信息
}

// 在线信息
type Presence struct {
	Online          bool     `json:"online,omitempty"`           // 是否在线
	AutoAway        bool     `json:"auto_away,omitempty"`        // 是否自动离线
	ManualAway      bool     `json:"manual_away,omitempty"`      // 是否手动离线
	ConnectionCount int      `json:"connection_count,omitempty"` // 连接次数
	LastActivity    JSONTime `json:"last_activity,omitempty"`    // 上次登录时间
}

type JSONTime int64

func (t JSONTime) String() string {
	tm := t.Time()
	return tm.Format("2006-01-02 15:04:05")
}

func (t JSONTime) Time() time.Time {
	return time.Unix(int64(t), 0)
}
