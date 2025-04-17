package http

import (
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
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

	agents, err := server.AgentService.ListAgents()
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	eventLog, err := server.EventLogService.GetEventLog()
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	homePage := templates.Dashboard(
		templates.AgentsView(agents),
		templates.EventLogView(eventLog),
	)
	err = homePage.Render(r.Context(), w)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}

func (server *Server) GetTaskResults(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	taskResults, err := server.TaskResultsService.GetTaskResults(agentId)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	err = templates.TaskResults(taskResults).Render(r.Context(), w)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}
