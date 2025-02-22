package service

func (tqs *TaskQueueService) taskExists(taskId uint64) (bool, error) {
	exists, err := tqs.taskQueueRepository.TaskExists(taskId)
	return exists, repositoryErrWrapper(err)
}

func (tqs *TaskQueueService) GetTaskQueue() ([]string, error) {
	taskQueue, err := tqs.taskQueueRepository.GetTaskQueue()
	return taskQueue, repositoryErrWrapper(err)
}

func (tqs *TaskQueueService) TaskQueuePush(task string) error {
	return repositoryErrWrapper(tqs.taskQueueRepository.TaskQueuePush(task))
}
