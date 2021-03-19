package gcpine

import (
	"context"
	"fmt"
	"reflect"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Execute - separate execute for each event
func (g *GCPine) Execute(ctx context.Context, event *linebot.Event) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("panic: %+v", rec)
		}
	}()

	var eventType TracerName
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
		return fmt.Errorf("event type=%s: %w", eventType, ErrNoSetFunction)
	}

	stack, err := fn(ctx, g, event)
	if err != nil {
		return fmt.Errorf("error in event handler method: %w", err)
	}

	if len(stack) == 0 {
		return ErrEmptyMessages
	}

	if err = g.SendReplyMessage(event.ReplyToken, stack); err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}

	return nil
}