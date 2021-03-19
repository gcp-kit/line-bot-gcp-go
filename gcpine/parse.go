package gcpine

import (
	"encoding/json"
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

// ParseEvents - extract `[]*linebot.Event`
func ParseEvents(data []byte) ([]*linebot.Event, error) {
	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}

	if err := json.Unmarshal(data, request); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal: %w", err)
	}

	return request.Events, nil
}
