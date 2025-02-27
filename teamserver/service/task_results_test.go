package service_test

import (
	"errors"
	"testing"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

type saveTaskResutlsTestCase struct {
	agent         uint64
	taskResultsIn teamserver.TaskResultIn
}

func isClientError(err error) bool {
	target := &teamserver.ClientError{}
	return errors.As(err, &target)
}

func TestSaveAndGetTaskResults(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	validAgent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	const testTask = "test task"
	validTaskId, err := testedServices.taskQueueService.TaskQueuePush(testTask)
	if err != nil {
		t.Fatal(err)
	}

	given := []saveTaskResutlsTestCase{
		// Success
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: validTaskId, Output: "some output"}},

		// Failure
		{2137, teamserver.TaskResultIn{TaskId: validTaskId, Output: "some output"}},
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: 2137, Output: "some output"}},
		{2137, teamserver.TaskResultIn{TaskId: 2137, Output: "some output"}},
	}

	testCaseIndex := 0
	err = testedServices.taskResultsService.SaveTaskResults(given[testCaseIndex].agent, []teamserver.TaskResultIn{given[testCaseIndex].taskResultsIn})
	if err != nil {
		t.Fatal(err)
	}

	taskResult, err := testedServices.taskResultsService.GetTaskResult(validAgent.AgentId, given[testCaseIndex].taskResultsIn.TaskId)
	if err != nil {
		t.Fatal(err)
	}

	if taskResult.Output != given[testCaseIndex].taskResultsIn.Output ||
		taskResult.TaskId != validTaskId ||
		taskResult.Task != testTask {
		t.Error("Task result data doesn't match")
		t.Error("Expected:")
		t.Error(given[testCaseIndex].taskResultsIn.Output)
		t.Error(validTaskId)
		t.Error(testTask)
		t.Error("Got: ")
		t.Error(taskResult)
	}

	taskResults, err := testedServices.taskResultsService.GetTaskResults(validAgent.AgentId) // Same thing but different function, need to check that too
	if err != nil {
		t.Fatal(err)
	}

	if taskResults[0] != *taskResult {
		t.Error("Invalid task result returned by GetTaskdResults:")
		t.Error(taskResults[0])
	}

	testCaseIndex++
	for ; testCaseIndex < len(given); testCaseIndex++ {
		err = testedServices.taskResultsService.SaveTaskResults(given[testCaseIndex].agent, []teamserver.TaskResultIn{given[testCaseIndex].taskResultsIn})
		if err == nil {
			t.Fatal("No error detected with incorrect arguments")
		} else if !isClientError(err) {
			t.Fatal("Non client error returned with incorrect args")
		}
	}
}
