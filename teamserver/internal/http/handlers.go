package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

const (
	errSerializationFmt = "Failed to serialize register response with: %w"
	errUnauthorized     = "Unauthorized"
)

func sendEventMessage(w http.ResponseWriter, message string) error {
	_, err := w.Write([]byte(
		`event: message
data: ` + message + "\n"))
	return err
}

func ApiErrorHandler(err error, w http.ResponseWriter) {
	target := &teamserver.ClientError{}

	if errors.As(err, &target) {
		api.RequestErrorHandler(w, err)
	} else { // Default to internal
		api.InternalErrorHandler(w)
	}

	log.Print(err)
}
