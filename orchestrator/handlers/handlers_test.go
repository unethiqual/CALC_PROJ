package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"DistributedArithmeticExpressionCalculator/orchestrator/handlers"
	"DistributedArithmeticExpressionCalculator/orchestrator/models"
	"DistributedArithmeticExpressionCalculator/orchestrator/parser"
	"DistributedArithmeticExpressionCalculator/orchestrator/scheduler"
)

func TestHandleCalculate(t *testing.T) {
	models.Expressions = make(map[int64]*models.Expression)
	models.ExprIDCounter = 1
	models.TasksQueue = []*models.Task{}

	reqBody := []byte(`{"expression": "2+2*(2+5)*3"}`)
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleCalculate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp["id"] == nil {
		t.Error("response missing 'id' field")
	}
}

func TestHandleGetExpressions(t *testing.T) {
	models.Expressions = make(map[int64]*models.Expression)
	exprStr := "2+2*(2+5)*3"
	ast, err := parser.ParseExpression(exprStr)
	if err != nil {
		t.Fatal(err)
	}
	expr := &models.Expression{
		ID:     1,
		Expr:   exprStr,
		Status: models.StatusCompleted,
		Root:   ast,
		Result: 44,
	}
	models.Expressions[expr.ID] = expr

	req, err := http.NewRequest("GET", "/api/v1/expressions", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleGetExpressions)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if _, ok := resp["expressions"]; !ok {
		t.Error("expected 'expressions' array in response")
	}
}

func TestHandleGetExpression(t *testing.T) {
	models.Expressions = make(map[int64]*models.Expression)
	exprStr := "2+2*(2+5)*3"
	ast, err := parser.ParseExpression(exprStr)
	if err != nil {
		t.Fatal(err)
	}
	expr := &models.Expression{
		ID:     1,
		Expr:   exprStr,
		Status: models.StatusCompleted,
		Root:   ast,
		Result: 44,
	}
	models.Expressions[expr.ID] = expr

	req, err := http.NewRequest("GET", "/api/v1/expressions/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandleGetExpression)
	req.URL.Path = "/api/v1/expressions/1"
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	expression, ok := resp["expression"].(map[string]interface{})
	if !ok {
		t.Fatal("response format error: 'expression' field missing or invalid")
	}
	if expression["id"] == nil {
		t.Error("expected 'id' in expression object")
	}
}

func TestHandleTaskEndpoints(t *testing.T) {
	models.TasksQueue = []*models.Task{}
	models.Expressions = make(map[int64]*models.Expression)
	models.ExprIDCounter = 1
	models.TaskIDCounter = 1

	exprStr := "2+3"
	ast, err := parser.ParseExpression(exprStr)
	if err != nil {
		t.Fatal(err)
	}
	expr := &models.Expression{
		ID:     1,
		Expr:   exprStr,
		Status: models.StatusInProgress,
		Root:   ast,
	}
	models.Expressions[expr.ID] = expr
	scheduler.ScheduleTasks(expr.Root, expr.ID)
	if len(models.TasksQueue) == 0 {
		t.Fatal("expected at least one task in the queue")
	}
	task := models.TasksQueue[0]

	reqGet, err := http.NewRequest("GET", "/internal/task", nil)
	if err != nil {
		t.Fatal(err)
	}
	rrGet := httptest.NewRecorder()
	handlerGet := http.HandlerFunc(handlers.HandleGetTask)
	handlerGet.ServeHTTP(rrGet, reqGet)
	if rrGet.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rrGet.Code)
	}
	var getResp map[string]interface{}
	if err := json.Unmarshal(rrGet.Body.Bytes(), &getResp); err != nil {
		t.Fatal(err)
	}
	if getResp["task"] == nil {
		t.Error("expected task object in GET /internal/task response")
	}

	postBody := []byte(`{"id": ` + strconv.Itoa(int(task.ID)) + `, "result": 5}`)
	reqPost, err := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(postBody))
	if err != nil {
		t.Fatal(err)
	}
	reqPost.Header.Set("Content-Type", "application/json")
	rrPost := httptest.NewRecorder()
	handlerPost := http.HandlerFunc(handlers.HandlePostTask)
	handlerPost.ServeHTTP(rrPost, reqPost)
	if rrPost.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rrPost.Code)
	}
}
