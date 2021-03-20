package gcpine

import (
	"context"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

// GCPine - Toolkit of the LINE Bot to work for Google Cloud Platform.
type GCPine struct {
	ErrMessages []linebot.SendingMessage
	Function    map[EventType]func(ctx context.Context, pine *GCPine, event *linebot.Event) ([]linebot.SendingMessage, error)
	LiffFunc    map[string]func(r *http.Request, w http.ResponseWriter)
	*linebot.Client
}
