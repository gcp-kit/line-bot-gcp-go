package gcpine

import (
	"context"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

type (
	// GCPine - Toolkit of the LINE Bot to work for Google Cloud Platform.
	GCPine struct {
		ErrMessages []linebot.SendingMessage
		Function    map[EventType]GCPineFunction
		LiffFunc    map[string]GCPineLiffFunc
		*linebot.Client
	}

	// GCPineFunction - Event function for reply
	GCPineFunction func(ctx context.Context, pine *GCPine, event *linebot.Event) ([]linebot.SendingMessage, error)

	// GCPineLiffFunc - Functions for LIFF(TODO function)
	GCPineLiffFunc func(r *http.Request, w http.ResponseWriter)
)
