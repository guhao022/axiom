package axiom

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Provider interface {
	Name() string
	Run() error
	Close() error

	Receive(*Message) error
	Send(*Response, ...string) error
	Reply(*Response, ...string) error
}

var availableProviders = make(map[string]func(*Robot) (Provider, error))

func RegisterProvider(name string, f func(*Robot) (Provider, error)) {
	availableProviders[name] = f
}

type BasicProvider struct {
	*Robot
}

func (a *BasicProvider) SetRobot(r *Robot) {
	a.Robot = r
}

func NewProvider(robot *Robot) (Provider, error) {

	default_provider := `cli`
	if _, ok := availableProviders[default_provider]; !ok {
		return nil, fmt.Errorf("%s is not a registered provider", default_provider)
	}

	provider, err := availableProviders[default_provider](robot)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

type cli struct {
	BasicProvider
	quit   chan bool
	writer *bufio.Writer
}

func NewCli(r *Robot) (Provider, error) {
	c := &cli{
		quit:   make(chan bool),
		writer: bufio.NewWriter(os.Stdout),
	}
	c.SetRobot(r)
	return c, nil
}

func (c *cli) Name() string {
	return `cli`
}

func (c *cli) Receive(msg *Message) error {
	c.Robot.Receive(msg)
	return nil
}

func (c *cli) Send(res *Response, strings ...string) error {
	for _, str := range strings {
		err := c.writeString(str)
		if err != nil {
			log.Printf("send message error: %v", err)
			return err
		}
	}

	return nil
}

func (c *cli) Reply(res *Response, strings ...string) error {
	for _, str := range strings {
		s := res.UserName() + `: ` + str
		err := c.writeString(s)
		if err != nil {
			log.Printf("reply message error: %v", err)
			return err
		}
	}

	return nil
}

func (c *cli) Run() error {

	go func() {
		for {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()

			line := scanner.Text()
			line = strings.TrimSpace(line)

			msg := &Message{
				ID:   "local-message",
				FromUser: User{ID: "1", Name: "cli"},
				Room: "cli",
				Text: scanner.Text(),
			}

			c.Receive(msg)

			prompt()
			continue
		}
	}()

	<-c.quit
	return nil
}

func (c *cli) Close() error {
	c.quit <- true
	return nil
}

func prompt() {
	fmt.Print("> ")
}

func (c *cli) writeString(str string) error {
	msg := fmt.Sprintf("%s\n", strings.TrimSpace(str))

	if _, err := c.writer.WriteString(msg); err != nil {
		return err
	}

	if err := c.writer.Flush(); err != nil {
		return err
	}

	return nil
}

func init() {
	RegisterProvider(`cli`, NewCli)
}
