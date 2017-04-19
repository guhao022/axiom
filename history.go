package axiom

import "time"

const (
	DEFAULT_HISTORY_ID       = 1
	DEFAULT_HISTORY_LATEST   = 0
	DEFAULT_HISTORY_COUNT    = 100
	DEFAULT_HISTORY_DEADLINE = 1800
)

type History struct {
	ID       int64     `json:"id,omitempty"`      // 子消息ID
	Latest   int64     `json:"latest"`            // 最后一条消息ID
	Message  []Message `json:"message,omitempty"` // 消息内容
	Count    int       `json:"count"`             // 消息最大数
	LastCall time.Time `json:"last_call"`         // 最后消息接收时间
	Deadline float64     `json:"deadline"`          // 过期时间(s)
}

func NewHistory() *History {
	return &History{
		ID:       DEFAULT_HISTORY_ID,
		Latest:   DEFAULT_HISTORY_LATEST,
		Count:    DEFAULT_HISTORY_COUNT,
		LastCall: time.Now(),
		Deadline: DEFAULT_HISTORY_DEADLINE,
	}
}

func (h *History) Insert(msg Message) {
	now := time.Now()
	subTime := now.Sub(h.LastCall)

	if subTime.Seconds() > h.Deadline {
		h = NewHistory()
	}

	if len(h.Message) >= h.Count {
		h.Message = append(h.Message[1:h.Count], msg)
	} else {
		h.Latest++
		h.Message[h.Latest] = msg
	}
}
