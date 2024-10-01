package storage

import "errors"

type RaiStorage struct {
	results map[string]string
	status  map[string]string
}

func NewRaiStorage() *RaiStorage {
	return &RaiStorage{
		results: make(map[string]string),
		status:  make(map[string]string),
	}
}

func (rs *RaiStorage) GetTaskResult(taskID string) (*map[string]string, error) {
	value, exists := rs.results[taskID]
	if !exists {
		return nil, errors.New("Task not found")
	}

	result := map[string]string{
		"result": value,
	}

	return &result, nil
}

func (rs *RaiStorage) GetTaskStatus(taskID string) (*string, error) {
	value, exists := rs.status[taskID]
	if !exists {
		return nil, errors.New("Task not found")
	}
	return &value, nil
}

func (rs *RaiStorage) CreateTask(TaskID string) error {
	rs.status[TaskID] = "in_progress"
	return nil
}

func (rs *RaiStorage) DoTask(taskID string, result string) error {
	if _, exists := rs.status[taskID]; !exists {
		return errors.New("Task not found")
	}
	rs.status[taskID] = "ready"
	rs.results[taskID] = result
	return nil
}
