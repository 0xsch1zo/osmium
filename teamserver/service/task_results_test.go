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

func TestSaveAndGetTaskResult(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	validAgent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	const testTask = "test task"
	validTaskIds := make([]uint64, 2)
	validTaskIds[0], err = testedServices.tasksService.AddTask(validAgent.AgentId, testTask)
	if err != nil {
		t.Fatal(err)
	}

	validTaskIds[1], err = testedServices.tasksService.AddTask(validAgent.AgentId, testTask)
	if err != nil {
		t.Fatal(err)
	}

	given := []saveTaskResutlsTestCase{
		// Success
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: validTaskIds[0], Output: "some output"}},
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: validTaskIds[1], Output: "some output"}},

		// Failure
		{2137, teamserver.TaskResultIn{TaskId: validTaskIds[0], Output: "some output"}},
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: 2137, Output: "some output"}},
		{2137, teamserver.TaskResultIn{TaskId: 2137, Output: "some output"}},
	}

	testCaseIndex := 0

	for testCaseIndex < 2 {
		err = testedServices.taskResultsService.SaveTaskResult(given[testCaseIndex].agent, &given[testCaseIndex].taskResultsIn)
		if err != nil {
			t.Fatal(err)
		}

		taskResult, err := testedServices.taskResultsService.GetTaskResult(validAgent.AgentId, given[testCaseIndex].taskResultsIn.TaskId)
		if err != nil {
			t.Fatal(err)
		}

		if taskResult.Output != given[testCaseIndex].taskResultsIn.Output ||
			taskResult.TaskId != validTaskIds[testCaseIndex] ||
			taskResult.Task != testTask {
			fatalErrUnexpectedData(t, "Task result data doesn't match", struct {
				output string
				taskId uint64
				task   string
			}{given[testCaseIndex].taskResultsIn.Output, validTaskIds[testCaseIndex], testTask}, taskResult)
		}

		testCaseIndex++
	}

	for ; testCaseIndex < len(given); testCaseIndex++ {
		err = testedServices.taskResultsService.SaveTaskResult(given[testCaseIndex].agent, &given[testCaseIndex].taskResultsIn)
		if err == nil {
			t.Fatal("No error detected with incorrect arguments")
		} else if !isClientError(err) {
			t.Fatal("Non client error returned with incorrect args")
		}
	}
}

func TestGetTaskResults(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	validAgent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	const testTask = "test task"
	validTaskIds := make([]uint64, 2)
	validTaskIds[0], err = testedServices.tasksService.AddTask(validAgent.AgentId, testTask)
	if err != nil {
		t.Fatal(err)
	}

	validTaskIds[1], err = testedServices.tasksService.AddTask(validAgent.AgentId, testTask)
	if err != nil {
		t.Fatal(err)
	}

	given := []saveTaskResutlsTestCase{
		// Success
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: validTaskIds[0], Output: "some output"}},
		{validAgent.AgentId, teamserver.TaskResultIn{TaskId: validTaskIds[1], Output: "some output"}},
	}

	for testCaseIndex := 0; testCaseIndex < 2; testCaseIndex++ {
		err = testedServices.taskResultsService.SaveTaskResult(given[testCaseIndex].agent, &given[testCaseIndex].taskResultsIn)
		if err != nil {
			t.Fatal(err)
		}

		taskResults, err := testedServices.taskResultsService.GetTaskResults(validAgent.AgentId)
		if err != nil {
			t.Fatal(err)
		}

		if taskResults[testCaseIndex].Output != given[testCaseIndex].taskResultsIn.Output ||
			taskResults[testCaseIndex].TaskId != validTaskIds[testCaseIndex] ||
			taskResults[testCaseIndex].Task != testTask {
			fatalErrUnexpectedData(t, "Task result data doesn't match", struct {
				output string
				taskId uint64
				task   string
			}{given[testCaseIndex].taskResultsIn.Output, validTaskIds[testCaseIndex], testTask}, taskResults[testCaseIndex])
		}
	}
}
