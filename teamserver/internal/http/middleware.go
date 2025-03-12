package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(handlers ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(handlers) - 1; i >= 0; i-- {
			handler := handlers[i]
			next = handler(next)
		}

		return next
	}
}

func JsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func ServerSentEvents(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func (server *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var username api.AuthorizedRequest
		body := bytes.Buffer{}
		rCopy := io.TeeReader(r.Body, &body)
		err := json.NewDecoder(rCopy).Decode(&username)
		if err != nil {
			starget := &json.SyntaxError{}
			utarget := &json.UnmarshalTypeError{}
			if errors.As(err, &starget) || errors.As(err, &utarget) {
				api.RequestErrorHandler(w, err)
				return
			}

			api.InternalErrorHandler(w)
			log.Print(err)
			return
		}

		r.Body = io.NopCloser(&body)

		sessionToken, err := r.Cookie("session_token")
		if errors.Is(err, http.ErrNoCookie) {
			api.RequestErrorHandler(w, errors.New(errUnauthorized))
			return
		} else if err != nil {
			api.InternalErrorHandler(w)
			log.Print(err)
			return
		}

		if len(username.Username) == 0 {
			api.RequestErrorHandler(w, errors.New(errUnauthorized))
			return
		}

		sessionTokenDb, err := server.AuthorizationService.GetSessionToken(username.Username)
		if err != nil {
			api.InternalErrorHandler(w)
			log.Print(err)
			return
		}

		if len(sessionToken.Value) == 0 || sessionToken.Value != sessionTokenDb {
			api.RequestErrorHandler(w, errors.New(errUnauthorized))
			return
		}
		next.ServeHTTP(w, r)
	})
}
