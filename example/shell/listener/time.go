package listener

import (
	"axiom"
	"time"
)

type TimeListener struct{}

func (t *TimeListener) Handle() []*axiom.Listener {

	return []*axiom.Listener{
		{
			Regex: "time|时间|几点",
			HandlerFunc: func(c *axiom.Context) {
				layout := "2006-01-02 15:04:05"
				t := time.Now()
				c.Reply(c.Message.ToUser + " > 现在时间: " + t.Format(layout))
			},
		},
	}
}
