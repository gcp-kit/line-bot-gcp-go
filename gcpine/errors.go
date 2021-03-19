package gcpine

import "fmt"

// ErrEmptyMessages - []linebot.SendingMessage is empty
var ErrEmptyMessages = fmt.Errorf("no message to send")
