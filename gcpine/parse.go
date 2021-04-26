package gcpine

import (
	"encoding/json"

	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/xerrors"
)

// ParseEvents - extract `[]*linebot.Event`
func ParseEvents(data []byte) ([]*linebot.Event, error) {
	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}

	if err := json.Unmarshal(data, request); err != nil {
		return nil, xerrors.Errorf("failed to json unmarshal: %w", err)
	}

	return request.Events, nil
}
