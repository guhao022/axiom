package axiom

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
)

type Provider interface {
	IncomingChannel() chan Message
	OutgoingChannel() chan Message
	Error() error
}

func RegisterProvider(provider Provider) ListenerFunc {
	return func(bot *Robot) {
		log.Printf("bot: changing message provider %T\n", provider)
		bot.providerIn = provider.IncomingChannel()
		bot.providerOut = provider.OutgoingChannel()
	}
}

// 默认实现CLI
type providerCLI struct {
	in  chan Message
	out chan Message
}

func CLI() *providerCLI {
	cli := &providerCLI{
		in:  make(chan Message),
		out: make(chan Message),
	}
	go cli.loop()
	return cli
}

func (c *providerCLI) IncomingChannel() chan Message {
	return c.in
}

func (c *providerCLI) OutgoingChannel() chan Message {
	return c.out
}

func (c *providerCLI) Error() error {
	return nil
}

func (c *providerCLI) loop() {
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c.in <- Message{
				Room:         "CLI",
				FromUserID:   "CLI",
				FromUserName: "CLI",
				Message:      scanner.Text(),
			}
		loop:
			for {
				select {
				case msg := <-c.out:
					os.Stdout.WriteString(processOutMessage(msg))
				default:
					break loop
				}
			}
		}
	}()
	go func() {
		for msg := range c.out {
			os.Stdout.WriteString(processOutMessage(msg))
		}
	}()
}

func processOutMessage(msg Message) string {
	var finalMsg bytes.Buffer
	template.Must(template.New("tmpl").Parse(msg.Message)).Execute(&finalMsg, struct{ User string }{msg.ToUserID})

	return fmt.Sprintln("\nout:>", msg.Room, msg.ToUserID, msg.ToUserName, ":", finalMsg.String())
}
