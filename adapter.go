// 工作平台
package axiom

import (
	"os"
	"strings"
	"bufio"
)

// 用于实现机器人工作的适配器
type Adapter interface {
	GetName() string             // 适配器名称
	Process() error              // 处理
	Send(Message, string) error              // 初始化
	Reply(Message, string) error // 回复
}

const DEFAULT_NAME = "Axiom"

// 默认实现shell适配器
type Shell struct {
	name string
	bot  *Robot
}

func NewShell(bot *Robot) *Shell {
	return &Shell{
		name: DEFAULT_NAME,
		bot:  bot,
	}
}

func (s *Shell) Prepare() error {
	return nil
}

func (s *Shell) GetName() string {
	return s.name
}

func (s *Shell) Process() error {

	for {
		scanner := bufio.NewScanner(os.Stdin)
		os.Stdout.WriteString(s.bot.name + " > ")
		scanner.Scan()

		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "quit" || line == "q" || line == "exit" {
			os.Stdout.WriteString("GoodBye!\n")
			return nil
		}

		v := NewMessage(1, line)
		s.bot.ReceiveMessage(v)
		//continue
	}

	return nil

}

func (s *Shell) Reply(msg Message, message string) error {

	os.Stdout.WriteString(message + "\n")
	return nil
}