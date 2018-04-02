package shell

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"axiom"
)

func init() {
	axiom.RegisterAdapter("shell", New)
}

type adapter struct {
	axiom.BasicAdapter
	in   *bufio.Reader
	out  *bufio.Writer
	quit chan bool
}

// New returns an initialized adapter
func New(r *axiom.Robot) (axiom.Adapter, error) {
	adp := &adapter{
		out:  bufio.NewWriter(os.Stdout),
		in:   bufio.NewReader(os.Stdin),
		quit: make(chan bool),
	}
	adp.SetRobot(r)
	return adp, nil
}

func (a *adapter) Name() string {
	return `shell`
}

// Send sends a regular response
func (a *adapter) Send(res *axiom.Response, strings ...string) error {
	for _, str := range strings {
		err := a.writeString(str)
		if err != nil {
			log.Println("error: ", err)
			return err
		}
	}

	return nil
}

// Reply sends a direct response
func (a *adapter) Reply(res *axiom.Response, strings ...string) error {
	for _, str := range strings {
		s := res.UserName() + `: ` + str
		err := a.writeString(s)
		if err != nil {
			log.Println("error: ", err)
			return err
		}
	}

	return nil
}

// Emote performs an emote
func (a *adapter) Emote(res *axiom.Response, strings ...string) error {
	return nil
}

// Topic sets the topic
func (a *adapter) Topic(res *axiom.Response, strings ...string) error {
	return nil
}

// Play plays a sound
func (a *adapter) Play(res *axiom.Response, strings ...string) error {
	return nil
}

// Receive forwards a message to the robot
func (a *adapter) Receive(msg *axiom.Message) error {
	a.Robot.Receive(msg)
	return nil
}

// Run executes the adapter run loop
func (a *adapter) Run() error {
	prompt()

	go func() {
		for {
			line, _, err := a.in.ReadLine()
			message := a.newMessage(string(line))

			if err != nil {
				if err == io.EOF {
					break
					// a.Robot.signalChan <- syscall.SIGTERM
				}
				fmt.Println("error:", err)
			}
			a.Receive(message)
			prompt()
		}
	}()

	<-a.quit
	return nil
}

// Stop the adapter
func (a *adapter) Stop() error {
	a.quit <- true
	return nil
}

func prompt() {
	fmt.Print("> ")
}

// func newMessage(text string) *Message {
func (a *adapter) newMessage(text string) *axiom.Message {
	return &axiom.Message{
		ID:   "local-message",
		User: axiom.User{ID: "1", Name: "shell"},
		Room: "shell",
		Text: text,
	}
}

func (a *adapter) writeString(str string) error {
	msg := fmt.Sprintf("%s\n", strings.TrimSpace(str))

	if _, err := a.out.WriteString(msg); err != nil {
		return err
	}

	if err := a.out.Flush(); err != nil {
		return err
	}

	return nil
}
