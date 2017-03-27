package axiom

import (
	"bufio"
	"os"
	"strings"
	"os/user"
)

// 默认实现shell适配器
type Shell struct {
	bot     *Robot
}

func NewShell(bot *Robot) *Shell {

	return &Shell{
		bot: bot,
	}
}

func (s *Shell) Construct() error {
	return nil
}

func (s *Shell) Process() error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		os.Stdout.WriteString(s.bot.name + " > ")
		scanner.Scan()

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "quit" || line == "q" || line == "exit" {
			os.Stdout.WriteString("GoodBye!")
		}

		v := Message{
			ToUserName: u.Username,
			Text: line,
		}
		s.bot.ReceiveMessage(v)
	}

	return nil

}

func (s *Shell) Reply(msg Message, message string) error {

	os.Stdout.WriteString(message + "\n")
	return nil
}

