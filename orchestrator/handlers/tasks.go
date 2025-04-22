package handlers

import (
	"github.com/unethiqual/CALC_PROJ/orchestrator/models"
	"github.com/unethiqual/CALC_PROJ/orchestrator/scheduler"
	"encoding/json"
	"net/http"
)

// HandleGetTask – обработчик GET /internal/task
func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	models.Mu.Lock()
	defer models.Mu.Unlock()
	if len(models.TasksQueue) == 0 {
		http.Error(w, "Нет задач", http.StatusNotFound)
		return
	}
	task := models.TasksQueue[0]
	models.TasksQueue = models.TasksQueue[1:]
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": task,
	})
}

// HandlePostTask – обработчик POST /internal/task
func HandlePostTask(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		ID     int64   `json:"id"`
		Result float64 `json:"result"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидные данные", http.StatusUnprocessableEntity)
		return
	}
	models.Mu.Lock()
	defer models.Mu.Unlock()
	found := false
	for _, expr := range models.Expressions {
		if scheduler.UpdateASTWithTask(expr, req.ID, req.Result) {
			found = true
			if expr.Root.Computed {
				expr.Status = models.StatusCompleted
				expr.Result = expr.Root.Value
			}
		}
	}
	if !found {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
