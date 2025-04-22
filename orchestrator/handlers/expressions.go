package handlers

import (
	"DistributedArithmeticExpressionCalculator/orchestrator/models"
	"DistributedArithmeticExpressionCalculator/orchestrator/parser"
	"DistributedArithmeticExpressionCalculator/orchestrator/scheduler"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// HandleCalculate – обработчик POST /api/v1/calculate
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Expression string `json:"expression"`
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "Невалидные данные", http.StatusUnprocessableEntity)
		return
	}
	root, err := parser.ParseExpression(req.Expression)
	if err != nil {
		http.Error(w, "Ошибка разбора: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	models.Mu.Lock()
	expr := &models.Expression{
		ID:     models.ExprIDCounter,
		Expr:   req.Expression,
		Status: models.StatusPending,
		Root:   root,
	}
	models.ExprIDCounter++
	models.Expressions[expr.ID] = expr
	scheduler.ScheduleTasks(expr.Root, expr.ID)
	if len(models.TasksQueue) > 0 {
		expr.Status = models.StatusInProgress
	}
	models.Mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": expr.ID,
	})
}

// HandleGetExpressions – обработчик GET /api/v1/expressions
func HandleGetExpressions(w http.ResponseWriter, r *http.Request) {
	models.Mu.Lock()
	defer models.Mu.Unlock()
	list := make([]*models.Expression, 0, len(models.Expressions))
	for _, expr := range models.Expressions {
		list = append(list, expr)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": list,
	})
}

// HandleGetExpression – обработчик GET /api/v1/expressions/{id}
// Извлекаем ID из URL вручную.
func HandleGetExpression(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	models.Mu.Lock()
	defer models.Mu.Unlock()
	expr, ok := models.Expressions[id]
	if !ok {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": expr,
	})
}
