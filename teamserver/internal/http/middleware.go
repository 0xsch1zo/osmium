package http

import (
	"errors"
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
		f, ok := w.(http.Flusher)
		if !ok {
			ApiErrorHandler(errors.New(errFailedFlushSse), w)
		}

		f.Flush()
		next.ServeHTTP(w, r)
	})
}

func (server *Server) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			api.RequestErrorHandler(w, errors.New(errUnauthorized))
		} else if err != nil {
			api.InternalErrorHandler(w)
		}

		err = server.AuthorizationService.Authorize(token.Value)
		if err != nil {
			ApiErrorHandler(err, w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
