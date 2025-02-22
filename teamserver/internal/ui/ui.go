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
	taskQueueService   *service.TaskQueueService
	taskResultsService *service.TaskResultsService
}

func NewUi(as *service.AgentService, tqs *service.TaskQueueService, trr *service.TaskResultsService) *Ui {
	return &Ui{
		agentService:       as,
		taskQueueService:   tqs,
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

	taskQueueView, err := ui.taskQueueContainer()
	if err != nil {
		UiErrorHandler(w, r, err.Error())
	}

	homePage := templates.Index(agentsView, taskQueueView)
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

func (ui *Ui) taskQueueContainer() (templ.Component, error) {
	taskQueue, err := ui.taskQueueService.GetTaskQueue()
	if err != nil {
		return nil, err
	}

	return templates.TaskQueueView(taskQueue), nil
}
