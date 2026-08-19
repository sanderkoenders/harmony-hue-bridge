package ssdp

import (
	"fmt"
	"strings"
)

type Message struct {
	Method  string
	Headers map[string]string
}

func FromDataFrame(data []byte) (*Message, error) {
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	msg := &Message{
		Method:  strings.TrimSpace(lines[0]),
		Headers: make(map[string]string),
	}

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		name = strings.ToLower(strings.TrimSpace(name))
		value = strings.TrimSpace(value)

		msg.Headers[name] = value
	}

	return msg, nil
}

func NewMessage(method string) *Message {
	return &Message{
		Method:  method,
		Headers: make(map[string]string),
	}
}

func (msg *Message) AddHeader(name, value string) {
	name = strings.ToLower(strings.TrimSpace(name))
	value = strings.TrimSpace(value)

	msg.Headers[name] = value
}

func (m *Message) Header(name string) string {
	return m.Headers[strings.ToLower(name)]
}

func (m *Message) IsMSearch() bool {
	return strings.EqualFold(m.Method, "M-SEARCH * HTTP/1.1")
}

func (msg *Message) String() string {
	var sb strings.Builder
	sb.WriteString(msg.Method)
	sb.WriteString("\r\n")

	for name, value := range msg.Headers {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", name, value))
	}

	sb.WriteString("\r\n")

	return sb.String()
}
