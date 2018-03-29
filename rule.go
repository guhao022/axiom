package axiom

import (
	"log"
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"text/template"
	"regexp"
)

// 解析器
type RuleParser interface {
	Name() string                           // 解析器名称
	Boot(*Robot)                            // 引导启动
	HelpMessage(Robot, string) string       // 帮助信息
	HandleMessage(Robot, Message) []Message // 解析消息
}

// 注册解析器
func RegisterRuleset(rule RuleParser) ListenerFunc {
	return func(bot *Robot) {
		log.Printf("bot: registering ruleset %T", rule)
		rule.Boot(bot)
		bot.rules = append(bot.rules, rule)
	}
}

type Rule struct {
	Regex        string
	HelpMessage  string
	ParseMessage func(Robot, string, []string) []string
}

type regexRuleset struct {
	regexes map[string]*template.Template
	rules   []Rule
}

func (r *regexRuleset) Name() string {
	return `Regex Ruleset`
}

func (r *regexRuleset) Boot(*Robot) {
	//
}

func (r *regexRuleset) HelpMessage(robot Robot, _ string) string {
	botName := robot.Name()
	var helpMsg string
	for _, rule := range r.rules {
		var finalRegex bytes.Buffer
		r.regexes[rule.Regex].Execute(&finalRegex, struct{ RobotName string }{botName})

		helpMsg = fmt.Sprintln(helpMsg, finalRegex.String(), "-", rule.HelpMessage)
	}
	return strings.TrimLeftFunc(helpMsg, unicode.IsSpace)
}

func (r regexRuleset) HandleMessage(robot Robot, in Message) []Message {
	for _, rule := range r.rules {
		botName := robot.Name()
		if in.Direct {
			botName = ""
		}

		var finalRegex bytes.Buffer
		if _, ok := r.regexes[rule.Regex]; !ok {
			r.regexes[rule.Regex] = template.Must(template.New(rule.Regex).Parse(rule.Regex))
		}
		r.regexes[rule.Regex].Execute(&finalRegex, struct{ RobotName string }{botName})
		sanitizedRegex := strings.TrimSpace(finalRegex.String())
		re := regexp.MustCompile(sanitizedRegex)
		matched := re.MatchString(in.Message)
		if !matched {
			continue
		}

		args := re.FindStringSubmatch(in.Message)
		if ret := rule.ParseMessage(robot, in.Message, args); len(ret) > 0 {
			var retMsgs []Message
			for _, m := range ret {
				retMsgs = append(
					retMsgs,
					Message{
						Room:       in.Room,
						ToUserID:   in.FromUserID,
						ToUserName: in.FromUserName,
						Message:    m,
					},
				)
			}
			return retMsgs
		}
	}

	return []Message{}
}

func NewRegex(rules []Rule) *regexRuleset {
	r := &regexRuleset{
		regexes: make(map[string]*template.Template),
		rules:   rules,
	}
	for _, rule := range rules {
		r.regexes[rule.Regex] = template.Must(template.New(rule.Regex).Parse(rule.Regex))
	}
	return r
}