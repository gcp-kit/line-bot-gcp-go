package gcpine

import "fmt"

var (
	// ErrEmptyMessages - []linebot.SendingMessage is empty
	ErrEmptyMessages = fmt.Errorf("no message to send")
	// ErrNoSetFunction - function to be executed individually is not set
	ErrNoSetFunction = fmt.Errorf("no set function")
	// ErrInvalidMessageType - invalid message type
	ErrInvalidMessageType = fmt.Errorf("invalid message type")
)
