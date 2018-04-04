package axiom

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Provider interface {
	Name() string
	Run() error
	Close() error

	IncomingChannel() chan Message
	OutgoingChannel() chan Message
}

var availableProviders map[string]func(*Robot) (Provider, error)

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
	in     chan Message
	out    chan Message
	quit   chan bool
	writer *bufio.Writer
}

func NewCli(r *Robot) (Provider, error) {
	c := &cli{
		out:    make(chan Message),
		in:     make(chan Message),
		quit:   make(chan bool),
		writer: bufio.NewWriter(os.Stdout),
	}
	c.SetRobot(r)
	return c, nil
}

func (c *cli) Name() string {
	return `cli`
}

func (c *cli) IncomingChannel() chan Message {
	return c.in
}

func (c *cli) OutgoingChannel() chan Message {
	return c.out
}

func (c *cli) Run() error {
	prompt()

	go func() {
		for {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				c.in <- Message{
					ID:   "local-message",
					User: User{ID: "1", Name: "cli"},
					Room: "cli",
					Text: scanner.Text(),
				}
			}

			//c.Receive(message)
			prompt()
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
