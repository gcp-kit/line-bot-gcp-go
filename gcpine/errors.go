package gcpine

import "golang.org/x/xerrors"

var (
	// ErrEmptyMessages - []linebot.SendingMessage is empty
	ErrEmptyMessages = xerrors.Errorf("no message to send")
	// ErrNoSetFunction - function to be executed individually is not set
	ErrNoSetFunction = xerrors.Errorf("no set function")
	// ErrInvalidMessageType - invalid message type
	ErrInvalidMessageType = xerrors.Errorf("invalid message type")
)
