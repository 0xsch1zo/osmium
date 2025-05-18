package service_test

import (
	"errors"
	"strconv"
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

	validAgents := make([]teamserver.Agent, 2)
	for i := range len(validAgents) {
		agent, err := testedServices.agentService.AddAgent(teamserver.AgentRegisterInfo{})
		if err != nil {
			t.Fatal(err)
		}
		validAgents[i] = *agent
	}

	const testTask = "test task"
	validTaskIds := make([]uint64, 2)
	for i := range len(validTaskIds) {
		validTaskIds[i], err = testedServices.tasksService.AddTask(validAgents[0].AgentId, testTask)
		if err != nil {
			t.Fatal(err)
		}
	}

	given := []saveTaskResutlsTestCase{
		// Success
		{validAgents[0].AgentId, teamserver.TaskResultIn{TaskId: validTaskIds[0], Output: "some output"}},
		{validAgents[0].AgentId, teamserver.TaskResultIn{TaskId: validTaskIds[1], Output: "some output"}},

		// Failure
		{2137, teamserver.TaskResultIn{TaskId: validTaskIds[0], Output: "some output"}},
		{validAgents[0].AgentId, teamserver.TaskResultIn{TaskId: 2137, Output: "some output"}},
		{2137, teamserver.TaskResultIn{TaskId: 2137, Output: "some output"}},
	}

	testCaseIndex := 0

	for testCaseIndex < 2 {
		err = testedServices.taskResultsService.SaveTaskResult(given[testCaseIndex].agent, &given[testCaseIndex].taskResultsIn)
		if err != nil {
			t.Fatal(err)
		}

		taskResult, err := testedServices.taskResultsService.GetTaskResult(validAgents[0].AgentId, given[testCaseIndex].taskResultsIn.TaskId)
		if err != nil {
			t.Fatal(err)
		}

		if taskResult.Output != given[testCaseIndex].taskResultsIn.Output ||
			taskResult.TaskId != validTaskIds[testCaseIndex] ||
			taskResult.Task != testTask {
			fatalErrUnexpectedData(t, "Task result data doesn't match", struct {
				taskId uint64
				task   string
				output string
			}{validTaskIds[testCaseIndex], testTask, given[testCaseIndex].taskResultsIn.Output}, taskResult)
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

func TestGetTaskResultAndResultsMultiAgent(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agents := make([]teamserver.Agent, 2)
	tasks := make([]teamserver.Task, len(agents))
	const testTaskPref = "test task"
	for i := range len(agents) {
		agent, err := testedServices.agentService.AddAgent(teamserver.AgentRegisterInfo{})
		if err != nil {
			t.Fatal(err)
		}
		agents[i] = *agent

		tasks[i].Task = testTaskPref + strconv.FormatUint(agent.AgentId, 10)
		tasks[i].TaskId, err = testedServices.tasksService.AddTask(agent.AgentId, tasks[i].Task)
		if err != nil {
			t.Fatal(err)
		}
	}

	taskResults := make([]teamserver.TaskResultIn, len(tasks))
	const testTaskResultPref = "test output"
	for i := range len(taskResults) {
		taskResults[i] = teamserver.TaskResultIn{
			TaskId: tasks[i].TaskId,
			Output: testTaskResultPref + strconv.FormatUint(agents[i].AgentId, 10),
		}

		err = testedServices.taskResultsService.SaveTaskResult(agents[i].AgentId, &taskResults[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := range len(agents) {
		taskResultGot, err := testedServices.taskResultsService.GetTaskResult(agents[i].AgentId, tasks[i].TaskId)
		if err != nil {
			t.Fatal(err)
		}

		if taskResultGot.Output != taskResults[i].Output ||
			taskResultGot.TaskId != taskResults[i].TaskId ||
			taskResultGot.Task != tasks[i].Task {
			fatalErrUnexpectedData(t, "Task result data doesn't match", struct {
				taskId uint64
				task   string
				output string
			}{tasks[i].TaskId, tasks[i].Task, taskResults[i].Output}, *taskResultGot)
		}

		taskResultsGot, err := testedServices.taskResultsService.GetTaskResults(agents[i].AgentId)
		if err != nil {
			t.Fatal(err)
		}

		if len(taskResultsGot) != 1 {
			t.Error("Unexpected task results recieved length")
			fatalErrUnexpectedData(t, "Task result data doesn't match", struct {
				taskId uint64
				task   string
				output string
			}{tasks[i].TaskId, tasks[i].Task, taskResults[i].Output}, taskResultsGot)
		}

		if taskResultsGot[0].Output != taskResults[i].Output ||
			taskResultsGot[0].TaskId != taskResults[i].TaskId ||
			taskResultsGot[0].Task != tasks[i].Task {
			fatalErrUnexpectedData(t, "Task result data doesn't match", struct {
				taskId uint64
				task   string
				output string
			}{tasks[i].TaskId, tasks[i].Task, taskResults[i].Output}, *taskResultGot)
		}
	}
}

func TestGetTaskResults(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	validAgent, err := testedServices.agentService.AddAgent(teamserver.AgentRegisterInfo{})
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
