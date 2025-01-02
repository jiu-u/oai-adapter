package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMessages(t *testing.T) {
	messages := []Message{
		&DeveloperMessage{
			Role:    "developer",
			Content: "123",
		},
		&SystemMessage{
			Role:    "system",
			Content: "123",
		},
	}

	res, err := json.Marshal(messages)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(res))
	var msgs []Message
	err = json.Unmarshal(res, &msgs)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", msgs)
	//for _, msg := range messages {
	//	_, err := msg.ToMessage()
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}
}
