package axiom

import (
	"sort"
	"sync"
	"time"
)

const (
	DEFAULT_HISTORY_COUNT    = 100
	DEFAULT_HISTORY_DEADLINE = 1800
)

type history struct {
	Message  []Message     `json:"message,omitempty"` // 消息内容
	Count    int           `json:"count"`             // 消息最大数
	LastCall time.Time     `json:"last_call"`         // 最后消息接收时间
	Deadline time.Duration `json:"deadline"`          // 过期时间(s)

	lock *sync.Mutex
}

func (h history) Len() int {
	return len(h.Message)
}

func (h history) Swap(i, j int) {
	h.Message[i], h.Message[j] = h.Message[j], h.Message[i]
}

func (h history) Less(i, j int) bool {
	return h.Message[i].ID > h.Message[j].ID
}

func NewHistory(count int) *history {
	var message []Message

	return &history{
		Message:  message,
		Count:    count,
		LastCall: time.Now(),
		Deadline: DEFAULT_HISTORY_DEADLINE,
		lock:     new(sync.Mutex),
	}
}

func (h *history) Insert(msg Message) error {
	h.lock.Lock()
	defer h.lock.Unlock()

	now := time.Now()
	subTime := now.Sub(h.LastCall)

	if subTime > h.Deadline {
		h.Flush()
	}

	if len(h.Message) >= h.Count {
		h.Message = append(h.Message[1:h.Count], msg)
	} else {
		h.Message = append(h.Message, msg)
	}

	sort.Sort(h)

	return nil
}

func (h *history) Flush() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.Message = nil
	h.LastCall = time.Now()
}

func (h *history) gc() {
	for {
		h.lock.Lock()
		now := time.Now()
		subTime := now.Sub(h.LastCall)
		if subTime > h.Deadline {
			h.Flush()
		}
		h.lock.Unlock()
		time.Sleep(h.Deadline)
	}
}
