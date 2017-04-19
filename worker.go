package axiom

import (
	"time"
	"sort"
)

const (
	DEFAULT_WORKER_ID       = 1
	DEFAULT_WORKER_LATEST   = 0
	DEFAULT_WORKER_COUNT    = 100
	DEFAULT_WORKER_DEADLINE = 1800
)

type Worker struct {
	ID       int64     `json:"id,omitempty"`      // 工作台ID
	Latest   int64     `json:"latest"`            // 下一条workerID
	Message  []Message `json:"message,omitempty"` // 消息内容
	Count    int       `json:"count"`             // 消息最大数
	LastCall time.Time `json:"last_call"`         // 最后消息接收时间
	Deadline float64   `json:"deadline"`          // 过期时间(s)
}

func (w Worker) Len() int {
	return len(w.Message)
}

func (w Worker) Swap(i, j int) {
	w.Message[i], w.Message[j] = w.Message[j], w.Message[i]
}

func (w Worker) Less(i, j int) bool {
	return w.Message[i].ID > w.Message[j].ID
}

type WS []Worker

func (ws WS) Len() int {
	return len(ws)
}

func (ws WS) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

func (ws WS) Less(i, j int) bool {
	return ws[i].ID > ws[j].ID
}

func NewWorker() *Worker {
	return &Worker{
		ID:       DEFAULT_WORKER_ID,
		Latest:   DEFAULT_WORKER_LATEST,
		Count:    DEFAULT_WORKER_COUNT,
		LastCall: time.Now(),
		Deadline: DEFAULT_WORKER_DEADLINE,
	}
}

func (w Worker) Insert(msg Message) {
	now := time.Now()
	subTime := now.Sub(w.LastCall)

	if subTime.Seconds() > w.Deadline {
		w = NewWorker()
	}

	if len(w.Message) >= w.Count {
		w.Message = append(w.Message[1:w.Count], msg)
	} else {
		w.Message[len(w.Message)] = msg
	}

	sort.Sort(w)
}
