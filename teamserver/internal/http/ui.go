package http

import (
	"github.com/sentientbottleofwine/osmium/teamserver/internal/ui/templates"
	"net/http"
)

func (server *Server) RootHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		loginPage := templates.LoginPage()
		err := loginPage.Render(r.Context(), w)
		if err != nil {
			ApiErrorHandler(err, w)
			return
		}
		return
	} else if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	err = server.AuthorizationService.Authorize(token.Value)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	agentsView, err := server.AgentService.ListAgents()
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	homePage := templates.Dashboard(templates.AgentsView(agentsView))
	err = homePage.Render(r.Context(), w)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}
