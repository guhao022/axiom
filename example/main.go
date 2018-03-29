package main

import (
	"axiom"
)

func main() {
	axiom.NewBot().Process()
}


var regexRules = []axiom.Rule{
	{
		`jump`, `tells the robot to jump`,
		func(bot axiom.Robot, msg string, matches []string) []string {
			var ret []string
			ret = append(ret, "{{ .User }}, How high?")

			return ret
		},
	},
}

func init() {
	axiom.RegisterProvider(axiom.CLI())
	axiom.RegisterRuleset(axiom.NewRegex(regexRules))
}
