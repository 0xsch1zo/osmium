package http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

const (
	errSerializationFmt = "Failed to serialize register response with: %w"
	errUnauthorized     = "Unauthorized"
	errFailedFlushSse   = "Failed to flush sse headers"
)

func sendSSE(w http.ResponseWriter, messageType string, message string) error {
	_, err := w.Write([]byte(fmt.Sprintf("event: %s\ndata: %s\n", messageType, message)))
	f, ok := w.(http.Flusher)
	if !ok {
		ApiErrorHandler(errors.New(errFailedFlushSse), w)
	}

	f.Flush()
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

func SocketErrorHandler(err error, conn *websocket.Conn) {
	target := &teamserver.ClientError{}

	var message string
	var code int
	if errors.As(err, &target) {
		code = websocket.ClosePolicyViolation
		message = err.Error()
	} else { // Default to internal
		code = websocket.CloseInternalServerErr
		message = "An internal server error has occured."
	}

	log.Print(err)
	conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, message),
		time.Now().Add(1*time.Second),
	)
}
