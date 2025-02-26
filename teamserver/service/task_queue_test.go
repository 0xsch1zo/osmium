package service_test

import (
	"slices"
	"strconv"
	"testing"
)

func TestTaskQueuePushAndTaskExists(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	taskId, err := testedServices.taskQueueService.TaskQueuePush("test task")
	if err != nil {
		t.Fatal(err)
	}

	exists, err := testedServices.taskQueueService.TaskExists(taskId)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("Pushed task doesn't exist, id: " + strconv.FormatUint(taskId, 10))
	}
}

func TestGetTaskQueue(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	tasksGiven := []string{
		"test task 1",
		"test task 2",
		"test task 3",
	}

	for _, task := range tasksGiven {
		_, err := testedServices.taskQueueService.TaskQueuePush(task)
		if err != nil {
			t.Fatal(err)
		}
	}

	tasksGot, err := testedServices.taskQueueService.GetTaskQueue()
	if err != nil {
		t.Fatal(err)
	}

	if slices.Compare(tasksGiven, tasksGot) != 0 {
		t.Error("Tasks returned don't match with the original list.")
		t.Error("Expected:")
		t.Error(tasksGiven)
		t.Error("Got:")
		t.Fatal(tasksGot)
	}
}
