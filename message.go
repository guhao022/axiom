package axiom

import "sync"

type Message struct {
	ID      int64      `json:"id"`                      // 消息ID
	Text    string     `json:"text"`                    // 消息内容
	Adapter []*Adapter `json:"adapter"`                 // 适配平台
	ReplyTo []*User    `json:"reply_to"`                // 输出适配器
	History *history   `json:"group_message,omitempty"` // 消息组
}

type history struct {
	ID      int64     `json:"id,omitempty"`          // 子消息ID
	Message []Message `json:"sub_message,omitempty"` // 消息内容
}

func (h *history) Init() *history {
	//id := NewSafeID()
	// 获取储存的ID
	// 生成新的ID
}

func (h *history) Insert(id int64, msg string) error {
	return nil
}

func (h *history) List(id int64) error {

	return nil
}

type IDGenerator interface {
	Next() int64
}

func NewSafeID(startID int64) IDGenerator {
	return &safeID{
		nextID: startID,
		mutex:  &sync.Mutex{},
	}
}

type safeID struct {
	nextID int64
	mutex  *sync.Mutex
}

func (s *safeID) Next() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	id := s.nextID
	s.nextID++
	return id
}