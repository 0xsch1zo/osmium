package service_test

import (
	"slices"
	"testing"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func TestAddTaskAndTaskExists(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agentId, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	taskId, err := testedServices.tasksService.AddTask(agentId.AgentId, "test task")
	if err != nil {
		t.Fatal(err)
	}

	err = testedServices.tasksService.TaskExists(agentId.AgentId, taskId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTasks(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	tasksGiven := []string{
		"test task 1",
		"test task 2",
		"test task 3",
	}

	var someTaskId uint64
	for _, task := range tasksGiven {
		someTaskId, err = testedServices.tasksService.AddTask(agent.AgentId, task)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = testedServices.tasksService.UpdateTaskStatus(agent.AgentId, someTaskId, teamserver.TaskFinished)
	if err != nil {
		t.Fatal(err)
	}

	tasksGot, err := testedServices.tasksService.GetTasks(agent.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	var tasksStrGot []string

	for _, task := range tasksGot {
		tasksStrGot = append(tasksStrGot, task.Task)
	}

	if slices.Compare(tasksGiven, tasksStrGot) != 0 {
		fatalErrUnexpectedData(
			t,
			"Tasks returned don't match with the original list.",
			tasksGiven,
			tasksGot,
		)
	}
}

func TestGetTasksWithStatuses(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	tasksGiven := []teamserver.Task{
		{TaskId: 0, Task: "test task 1"},
		{TaskId: 0, Task: "test task 2"},
	}

	for i := range tasksGiven {
		taskId, err := testedServices.tasksService.AddTask(agent.AgentId, tasksGiven[1].Task)
		if err != nil {
			t.Fatal(err)
		}

		tasksGiven[i].TaskId = taskId
	}

	err = testedServices.tasksService.UpdateTaskStatus(agent.AgentId, tasksGiven[0].TaskId, teamserver.TaskFinished)
	if err != nil {
		t.Fatal(err)
	}

	tasksGot, err := testedServices.tasksService.GetNewTasks(agent.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasksGot) != 1 ||
		tasksGot[0].TaskId != tasksGiven[1].TaskId ||
		tasksGot[0].Task != tasksGiven[1].Task {

		fatalErrUnexpectedData(
			t,
			"Tasks returned don't match with the original list.",
			tasksGiven,
			tasksGot,
		)
	}
}
