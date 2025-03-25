package ui

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/ui/templates"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Ui struct {
	agentService         *service.AgentService
	tasksService         *service.TasksService
	taskResultsService   *service.TaskResultsService
	authorizationService *service.AuthorizationService
}

func NewUi(as *service.AgentService, ts *service.TasksService, trr *service.TaskResultsService, auths *service.AuthorizationService) *Ui {
	return &Ui{
		agentService:         as,
		tasksService:         ts,
		taskResultsService:   trr,
		authorizationService: auths,
	}
}

func UiErrorHandler(w http.ResponseWriter, r *http.Request, error string) {
	errorPage := templates.ErrorPage(error)
	w.WriteHeader(500)
	err := errorPage.Render(r.Context(), w)
	if err != nil {
		api.InternalErrorHandler(w)
	}
}

func (ui *Ui) RootHandler(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		loginPage := templates.LoginPage()
		err := loginPage.Render(r.Context(), w)
		if err != nil {
			UiErrorHandler(w, r, err.Error())
			return
		}
		return
	} else if err != nil {
		UiErrorHandler(w, r, err.Error())
		return
	}

	err = ui.authorizationService.Authorize(token.Value)
	if err != nil {
		UiErrorHandler(w, r, err.Error())
		return
	}

	agentsView, err := ui.agentsContainer()
	if err != nil {
		UiErrorHandler(w, r, err.Error())
		return
	}

	homePage := templates.Dashboard(agentsView)
	err = homePage.Render(r.Context(), w)
	if err != nil {
		UiErrorHandler(w, r, err.Error())
		return
	}
}

func (ui *Ui) agentsContainer() (templ.Component, error) {
	agents, err := ui.agentService.ListAgents()
	if err != nil {
		return nil, err
	}

	return templates.AgentsView(agents), nil
}
