package scheduler_test

import (
	"testing"

	"github.com/unethiqual/CALC_PROJ/orchestrator/models"
	"github.com/unethiqual/CALC_PROJ/orchestrator/parser"
	"github.com/unethiqual/CALC_PROJ/orchestrator/scheduler"
)

func TestScheduleAndUpdateAST(t *testing.T) {
	exprStr := "2+3"
	ast, err := parser.ParseExpression(exprStr)
	if err != nil {
		t.Fatalf("error parsing expression: %v", err)
	}
	expr := &models.Expression{
		ID:     1,
		Expr:   exprStr,
		Status: models.StatusInProgress,
		Root:   ast,
	}
	models.TasksQueue = []*models.Task{}
	scheduler.ScheduleTasks(expr.Root, expr.ID)
	if len(models.TasksQueue) == 0 {
		t.Fatal("expected at least one scheduled task")
	}
	task := models.TasksQueue[0]
	result := task.Arg1 + task.Arg2
	updated := scheduler.UpdateASTWithTask(expr, task.ID, result)
	if !updated {
		t.Fatal("UpdateASTWithTask failed to update AST")
	}
	if !expr.Root.Computed {
		t.Error("expected root node to be computed after update")
	}
	if expr.Root.Value != result {
		t.Errorf("expected %v, got %v", result, expr.Root.Value)
	}
}
