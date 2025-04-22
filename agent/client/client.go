package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

var BaseURL = "http://localhost:8080"

// Task – тип задачи
type Task struct {
	ID            int64   `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

// GetTask запрашивает задачу у оркестратора.
func GetTask() (*Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	var res struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res.Task, nil
}

// PostResult отправляет результат выполнения задачи.
func PostResult(taskID int64, result float64) error {
	data := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}
	body, _ := json.Marshal(data)
	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// Compute выполняет операцию с имитацией задержки.
func Compute(task *Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}