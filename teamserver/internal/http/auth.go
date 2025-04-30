package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	username := r.Form["username"]
	password := r.Form["password"]
	if len(username) == 0 || len(password) == 0 {
		api.RequestErrorHandler(w, errors.New(errUnauthorized), http.StatusUnauthorized)
		return
	}

	targetInvalidCreds := &teamserver.ClientError{}
	token, err := server.AuthorizationService.Login(username[0], password[0])
	if errors.As(err, &targetInvalidCreds) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("content-type", "text/html")
		err = templates.LoginForm(true).Render(r.Context(), w)
		if err != nil {
			ApiErrorHandler(err, w)
			return
		}
		return
	} else if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	w.Header().Add("HX-Redirect", "/")
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  token.ExpiryTime,
		Value:    token.Token,
	})
}

func (server *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		api.RequestErrorHandler(w, errors.New(errUnauthorized), http.StatusUnauthorized)
		return
	}

	refreshedToken, err := server.AuthorizationService.RefreshToken(token.Value)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshedToken.ExpiryTime,
		Value:    refreshedToken.Token,
	})
}

func (server *Server) GetRefreshTime(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		api.RequestErrorHandler(w, errors.New(errUnauthorized), http.StatusUnauthorized)
		return
	}

	refTime, err := server.AuthorizationService.GetRefreshTime(token.Value)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(api.GetRefreshTimeResponse{RefTime: refTime})
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}
