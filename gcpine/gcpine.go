package gcpine

import (
	"context"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type (
	// GCPine - Toolkit of the LINE Bot to work for Google Cloud Platform.
	GCPine struct {
		ErrMessages []linebot.SendingMessage
		Function    map[EventType]PineFunction
		LiffFunc    map[string]PineLiffFunc
		*linebot.Client
	}

	// PineFunction - Event function for reply
	PineFunction func(ctx context.Context, pine *GCPine, event *linebot.Event) ([]linebot.SendingMessage, error)

	// PineLiffFunc - Functions for LIFF(TODO function)
	PineLiffFunc func(r *http.Request, w http.ResponseWriter)
)
