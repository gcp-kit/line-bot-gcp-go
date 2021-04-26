package gcpine

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"golang.org/x/xerrors"
)

// SendReplyMessage - send reply message
func (g *GCPine) SendReplyMessage(token string, messages []linebot.SendingMessage) error {
	if _, err := g.ReplyMessage(token, messages...).Do(); err != nil {
		return xerrors.Errorf("faild to send reply message: %w", err)
	}

	return nil
}

// SendPushMessage - send push message
func (g *GCPine) SendPushMessage(uid string, messages []linebot.SendingMessage) error {
	if _, err := g.PushMessage(uid, messages...).Do(); err != nil {
		return xerrors.Errorf("faild to send push message: %w", err)
	}

	return nil
}
