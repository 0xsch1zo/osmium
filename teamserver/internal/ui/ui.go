package ui

import (
	"github.com/a-h/templ"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/ui/templates"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
	"net/http"
)

type Ui struct {
	agentService       *service.AgentService
	tasksService       *service.TasksService
	taskResultsService *service.TaskResultsService
}

func NewUi(as *service.AgentService, ts *service.TasksService, trr *service.TaskResultsService) *Ui {
	return &Ui{
		agentService:       as,
		tasksService:       ts,
		taskResultsService: trr,
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
	agentsView, err := ui.agentsContainer()
	if err != nil {
		UiErrorHandler(w, r, err.Error())
	}

	homePage := templates.Index(agentsView)
	err = homePage.Render(r.Context(), w)
	if err != nil {
		UiErrorHandler(w, r, err.Error())
	}
}

func (ui *Ui) agentsContainer() (templ.Component, error) {
	agents, err := ui.agentService.ListAgents()
	if err != nil {
		return nil, err
	}

	return templates.AgentsView(agents), nil
}
