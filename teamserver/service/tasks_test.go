package service_test

import (
	"slices"
	"testing"
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

func TestGetTaskQueue(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agentId, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	tasksGiven := []string{
		"test task 1",
		"test task 2",
		"test task 3",
	}

	for _, task := range tasksGiven {
		_, err := testedServices.tasksService.AddTask(agentId.AgentId, task)
		if err != nil {
			t.Fatal(err)
		}
	}

	tasksGot, err := testedServices.tasksService.GetTasks(agentId.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	var tasksStrGot []string

	for _, task := range tasksGot {
		tasksStrGot = append(tasksStrGot, task.Task)
	}

	if slices.Compare(tasksGiven, tasksStrGot) != 0 {
		t.Error("Tasks returned don't match with the original list.")
		t.Error("Expected:")
		t.Error(tasksGiven)
		t.Error("Got:")
		t.Fatal(tasksGot)
	}
}
