package gemini

import (
	"encoding/json"
	"testing"
)

func TestJSONData(t *testing.T) {
	input := `{"content": "Hello,
world!"}`
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}
