// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

// Main connection against ESL - Gotta add more description here
type SocketConnection struct {
	net.Conn
	err chan error
	m   chan *Message
}

// Dial - Will establish timedout dial against specified address. In this case, it will be freeswitch server
func (c *SocketConnection) Dial(network string, addr string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// Send - Will send raw message to open net connection
func (c *SocketConnection) Send(cmd string) error {

	if strings.Contains(cmd, "\r\n") {
		fmt.Errorf(EInvalidCommandProvided, cmd)
	}

	fmt.Fprintf(c, "%s\r\n\r\n", cmd)

	return nil
}

// SendMany - Will loop against passed commands and return 1st error if error happens
func (c *SocketConnection) SendMany(cmds []string) error {

	for _, cmd := range cmds {
		if err := c.Send(cmd); err != nil {
			return err
		}
	}

	return nil
}

// Execute - Helper fuck to execute commands with its args and sync/async mode
func (c *SocketConnection) Execute(command, args string, sync bool) (m *Message, err error) {
	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, "", "")
}

// ExecuteUUID - Helper fuck to execute uuid specific commands with its args and sync/async mode
func (c *SocketConnection) ExecuteUUID(uuid string, command string, args string, sync bool) (m *Message, err error) {
	return c.SendMsg(map[string]string{
		"call-command":     "execute",
		"execute-app-name": command,
		"execute-app-arg":  args,
		"event-lock":       strconv.FormatBool(sync),
	}, uuid, "")
}

// SendMsg - Basically this func will send message to the opened connection
func (c *SocketConnection) SendMsg(msg map[string]string, uuid, data string) (m *Message, err error) {

	b := bytes.NewBufferString("sendmsg")

	if uuid != "" {
		if strings.Contains(uuid, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		b.WriteString(" " + uuid)
	}

	b.WriteString("\n")

	for k, v := range msg {
		if strings.Contains(k, "\r\n") {
			return nil, fmt.Errorf(EInvalidCommandProvided, msg)
		}

		if v != "" {
			if strings.Contains(v, "\r\n") {
				return nil, fmt.Errorf(EInvalidCommandProvided, msg)
			}

			b.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		}
	}

	b.WriteString("\n")

	if msg["content-length"] != "" && data != "" {
		b.WriteString(data)
	}

	if _, err := b.WriteTo(c); err != nil {
		return nil, err
	}

	select {
	case err := <-c.err:
		return nil, err
	case m := <-c.m:
		return m, nil
	}
}

// OriginatorAdd - Will return originator address known as net.RemoteAddr()
// This will actually be a freeswitch address
func (c *SocketConnection) OriginatorAddr() net.Addr {
	return c.RemoteAddr()
}

// ReadMessage - Will read message from channels and return them back accordingy.
// If error is received, error will be returned. If not, message will be returned back!
func (c *SocketConnection) ReadMessage() (*Message, error) {
	Debug("Waiting for connection message to be received ...")

	select {
	case err := <-c.err:
		return nil, err
	case msg := <-c.m:
		return msg, nil
	}
}

// Handle - Will handle new messages and close connection when there are no messages left to process
func (c *SocketConnection) Handle() {

	done := make(chan bool)

	go func() {
		for {
			msg, err := newMessage(bufio.NewReaderSize(c, ReadBufferSize), true)

			if err != nil {
				c.err <- err
				done <- true
				break
			}

			c.m <- msg
		}
	}()

	<-done

	// Closing the connection now as there's nothing left to do ...
	c.Close()
}

// Close - Will close down net connection and return error if error happen
func (c *SocketConnection) Close() error {
	if err := c.Conn.Close(); err != nil {
		return err
	}

	return nil
}
