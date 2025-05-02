package scheduler

import (
    "errors"
    "github.com/unethiqual/CALC_PROJ/orchestrator/models"
)

func GetNextTask() (*models.Task, error) {
    models.Mu.Lock()
    defer models.Mu.Unlock()

    if len(models.TasksQueue) == 0 {
        return nil, errors.New("no tasks available")
    }

    task := models.TasksQueue[0]
    models.TasksQueue = models.TasksQueue[1:]
    return task, nil
}

func SubmitTaskResult(taskID int64, result float64) error {
    models.Mu.Lock()
    defer models.Mu.Unlock()

    for _, expr := range models.Expressions {
        if expr.ID == taskID {
            expr.Result = result
            expr.Status = "completed"
            return nil
        }
    }
    return errors.New("task not found")
}