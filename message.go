package axiom

import (
	"sync"
)

type Message struct {
	ID      int64      `json:"id"`       // 消息ID
	User    *User      `json:"user"`     // 发送用户
	Text    string     `json:"text"`     // 消息内容
	Adapter []*Adapter `json:"adapter"`  // 适配平台
	ReplyTo []*User    `json:"reply_to"` // 输出适配器
}

func NewMsg()

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
