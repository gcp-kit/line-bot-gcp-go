package gcpine

import (
	"context"
	"log"
	"reflect"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/xerrors"
)

// Execute - separate execute for each event
func (g *GCPine) Execute(ctx context.Context, event *linebot.Event) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = xerrors.Errorf("panic: %+v", rec)
		}
	}()

	var eventType EventType
	switch event.Type {
	case linebot.EventTypeMessage:
		switch message := event.Message.(type) {
		case *linebot.TextMessage,
			*linebot.ImageMessage,
			*linebot.VideoMessage,
			*linebot.AudioMessage,
			*linebot.FileMessage,
			*linebot.LocationMessage,
			*linebot.StickerMessage:
			eventType = reflect.TypeOf(message).Elem().Name()
		default:
			return ErrInvalidMessageType
		}
	default:
		eventType = string(event.Type)
	}

	fn, ok := g.Function[eventType]
	if !ok {
		log.Printf("event type=%s: %+v", eventType, ErrNoSetFunction)
		return nil
	}

	stack, err := fn(ctx, g, event)
	if err != nil {
		return xerrors.Errorf("error in event handler method: %w", err)
	}

	// NOTE: will not reply
	if stack == nil {
		return nil
	}

	if len(stack) == 0 {
		return ErrEmptyMessages
	}

	if err = g.SendReplyMessage(event.ReplyToken, stack); err != nil {
		return xerrors.Errorf("could not send message: %w", err)
	}

	return nil
}
