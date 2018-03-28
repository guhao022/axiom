package regex

import "github.com/num5/axiom"

type Rule struct {
	Regex        string
	HelpMessage  string
	ParseMessage func(axiom.Robot, string, []string) []string
}

type regexRuleset struct {
	rules   []Rule
}
