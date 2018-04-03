package axiom

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"io"
	"log"
)

// Adapter interface
type Adapter interface {
	Name() string
	Run() error
	Stop() error

	Receive(*Message) error
	Send(*Response, ...string) error
	Emote(*Response, ...string) error
	Reply(*Response, ...string) error
	Topic(*Response, ...string) error
	Play(*Response, ...string) error
}

type adapter struct {
	name     string
	newFunc  func(*Robot) (Adapter, error)
	sendChan chan *Response
	recvChan chan *Message
}

// AvailableAdapters is a map of registered adapters
var AvailableAdapters = map[string]adapter{}

// RegisterAdapter registers an adapter
func RegisterAdapter(name string, newFunc func(*Robot) (Adapter, error)) {
	AvailableAdapters[name] = adapter{
		name:    name,
		newFunc: newFunc,
	}
}

// BasicAdapter declares common functions shared by all adapters
type BasicAdapter struct {
	*Robot
}

// SetRobot sets the adapter's Robot
func (a *BasicAdapter) SetRobot(r *Robot) {
	a.Robot = r
}

type cli struct {
	BasicAdapter
	in   *bufio.Reader
	out  *bufio.Writer
	quit chan bool
}

// New returns an initialized adapter
func NewCli(r *Robot) (Adapter, error) {
	adp := &cli{
		out:  bufio.NewWriter(os.Stdout),
		in:   bufio.NewReader(os.Stdin),
		quit: make(chan bool),
	}
	adp.SetRobot(r)
	return adp, nil
}

func (c *cli) Name() string {
	return `cli`
}

// Send sends a regular response
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

// Reply sends a direct response
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

// Emote performs an emote
func (c *cli) Emote(res *Response, strings ...string) error {
	return nil
}

// Topic sets the topic
func (c *cli) Topic(res *Response, strings ...string) error {
	return nil
}

// Play plays a sound
func (c *cli) Play(res *Response, strings ...string) error {
	return nil
}

// Receive forwards a message to the robot
func (c *cli) Receive(msg *Message) error {
	c.Robot.Receive(msg)
	return nil
}

// Run executes the adapter run loop
func (c *cli) Run() error {
	prompt()

	go func() {
		for {
			line, _, err := c.in.ReadLine()
			message := c.newMessage(string(line))

			if err != nil {
				if err == io.EOF {
					break
				}
				log.Printf("run %s error: %v", cli.Name, err)
			}
			c.Receive(message)
			prompt()
		}
	}()

	<-c.quit
	return nil
}

// Stop the adapter
func (c *cli) Stop() error {
	c.quit <- true
	return nil
}

func prompt() {
	fmt.Print("> ")
}

func (c *cli) newMessage(text string) *Message {
	return &Message{
		ID:   "local-message",
		User: User{ID: "1", Name: "cli"},
		Room: "cli",
		Text: text,
	}
}

func (c *cli) writeString(str string) error {
	msg := fmt.Sprintf("%s\n", strings.TrimSpace(str))

	if _, err := c.out.WriteString(msg); err != nil {
		return err
	}

	if err := c.out.Flush(); err != nil {
		return err
	}

	return nil
}
