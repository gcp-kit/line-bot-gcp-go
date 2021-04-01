package gcpine

import (
	"net/http"
)

// Props - props for common.
type Props interface {
	SetGCPine(pine *GCPine)
	SetSecret(secret string)
	ReceiveWebHook(r *http.Request, w http.ResponseWriter) error
}
