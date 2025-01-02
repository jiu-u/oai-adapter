package openai

import (
	"encoding/json"
	"fmt"
	"github.com/jiu-u/oai-adapter/api"
	"testing"
)

func TestContentParse(t *testing.T) {
	mediaMsg1 := []api.MediaContent{
		{
			Type: "image_url",
			ImageUrl: &api.MessageImageUrl{
				Url:    "https://example.com/image.jpg",
				Detail: "This is an example image.",
			},
		},
		{
			Type: "Text",
			Text: "hello",
		},
	}
	byteMsg1, err := json.Marshal(mediaMsg1)
	if err != nil {
		t.Fatal(err)
	}
	mediaMsg2 := []api.MediaContent{
		{
			Type: "input_audio",
			InputAudio: &api.MessageInputAudio{
				Data:   "b64hello",
				Format: "mp3",
			},
		},
		{
			Type: "Text",
			Text: "hello",
		},
	}
	byteMsg2, err := json.Marshal(mediaMsg2)
	if err != nil {
		t.Fatal(err)
	}
	imput := api.ChatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []api.Message{
			{Role: "developer", Content: []byte(`"hello"`)},
			{Role: "user", Content: []byte(`"hello"`)},
			{Role: "assistant", Content: byteMsg1},
			{Role: "user", Content: []byte(`"hello"`)},
			{Role: "user", Content: byteMsg2},
		},
	}

	inputBytes, err := json.Marshal(imput)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(inputBytes))
	fmt.Println()
	var contentWarp api.ChatRequest
	err = json.Unmarshal(inputBytes, &contentWarp)
	if err != nil {
		t.Fatal(err)
	}

	for _, msg := range contentWarp.Messages {
		convertContent(msg.Content)
	}
}

func convertContent(input []byte) {
	var err error
	// 尝试将content解析为字符串
	var stringContent string
	err = json.Unmarshal(input, &stringContent)
	if err == nil {
		fmt.Println("Content is a string:", stringContent)
		return
	}
	// 尝试将content解析为[]MediaContent
	var mediaContent []api.MediaContent
	err = json.Unmarshal(input, &mediaContent)
	if err == nil {
		fmt.Println("Content is a slice of MediaContent:")
		for _, mc := range mediaContent {
			fmt.Printf("Type: %s\n", mc.Type)
			if mc.ImageUrl != nil {
				imageUrl, _ := json.Marshal(mc.ImageUrl)
				fmt.Printf("ImageUrl: %s\n", imageUrl)
			}
			if mc.InputAudio != nil {
				inputAudio, _ := json.Marshal(mc.InputAudio)
				fmt.Printf("InputAudio: %s\n", inputAudio)
			}
		}
		return
	}

	fmt.Println("Content is neither a string nor a slice of MediaContent")
}
