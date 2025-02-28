package main

import (
	"io"
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

type MemTableMbox struct {
	messages []string
	uidls    []string
}

func (m *MemTableMbox) Stat() (numberOfMessages int, totalSize int, err error) {
	for i := range m.messages {
		totalSize += len(m.messages[i])
	}
	numberOfMessages = len(m.messages)
	return
}

func (m *MemTableMbox) List() (messageSizes []int, err error) {
	messageSizes = make([]int, len(m.messages))
	for i := range messageSizes {
		messageSizes[i] = len(m.messages[i])
	}
	return
}

func (m *MemTableMbox) ListOne(msgNumber int) (size int, err error) {
	return len(m.messages[msgNumber]), nil
}

func (m *MemTableMbox) Message(msgNumber int) (msgReader io.ReadCloser, err error) {
	return io.NopCloser(strings.NewReader(m.messages[msgNumber])), nil
}

func (m *MemTableMbox) Dele(msgNumber int) error {
	return nil
}

func (m *MemTableMbox) Uidl() (uidls []string, err error) {
	return m.uidls, nil
}

func (m *MemTableMbox) UidlOne(msgNumber int) (uidl string, err error) {
	return m.uidls[msgNumber], nil
}

func (m *MemTableMbox) Close() error {
	return nil
}

func (m *MemTableMbox) Add(msg string) error {
	parsedMsg, err := mail.ReadMessage(strings.NewReader(msg))
	if err == nil {
		m.AddWithUidl(msg, parsedMsg.Header.Get("message-id"))
		return nil
	}
	uidl, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	m.AddWithUidl(msg, uidl.String())
	return nil
}

func (m *MemTableMbox) AddWithUidl(msg, uidl string) {
	m.messages = append(m.messages, msg)
	m.uidls = append(m.uidls, uidl)
}
