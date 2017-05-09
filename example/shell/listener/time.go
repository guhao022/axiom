package listener

import (
	"axiom"
	"time"
	"github.com/FrankWong1213/golang-lunar"
)

type TimeListener struct{}

func (t *TimeListener) Handle() []*axiom.Listener {

	return []*axiom.Listener{
		{
			Regex: "time|时间|几点",
			HandlerFunc: func(c *axiom.Context) {
				layout := "2006-01-02 15:04:05"
				t := time.Now()
				c.Reply(" 现在时间: " + t.Format(layout))
			},
		}, {
			Regex: "明天农历|明天是农历|农历明天",
			HandlerFunc: func(c *axiom.Context) {
				t := time.Now().AddDate(0, 0, 1)
				c.Reply(" 明天是农历 " + lunar.Lunar(t.Format("20060102")))
			},
		}, {
			Regex: "今天农历|今天是农历|农历今天|农历",
			HandlerFunc: func(c *axiom.Context) {
				t := time.Now()
				c.Reply(" 今天是农历 " + lunar.Lunar(t.Format("20060102")))
			},
		},
	}
}
