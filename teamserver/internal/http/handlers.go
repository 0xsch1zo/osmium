package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

const (
	errSerializationFmt = "Failed to serialize register response with: %w"
	errUnauthorized     = "Unauthorized"
)

func sendSSE(w http.ResponseWriter, messageType string, message string) error {
	_, err := w.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n", messageType, message)))
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
