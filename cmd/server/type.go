package main

import "encoding/json"

type Message interface {
	ToMessage() (json.RawMessage, error)
}

type DeveloperMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

func (m DeveloperMessage) ToMessage() (json.RawMessage, error) {
	return json.Marshal(m)
}

type SystemMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

func (m SystemMessage) ToMessage() (json.RawMessage, error) {
	return json.Marshal(m)
}

type UserMessage struct {
	Role    string      `json:"role"`
	Content UserContent `json:"content"`
	Name    string      `json:"name,omitempty"`
}

type UserContent interface {
	ToContent() (json.RawMessage, error)
}

type StringContent string

func (m StringContent) ToContent() (json.RawMessage, error) {
	return json.Marshal(m)
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
