package scheduler

import (
	"github.com/unethiqual/CALC_PROJ/orchestrator/config"
	"github.com/unethiqual/CALC_PROJ/orchestrator/models"
	"os"
	"strconv"
)

var GlobalConfig *config.Config

func ScheduleTasks(node *models.Node, exprID int64) {
	if node == nil {
		return
	}
	if node.Operator != "" {
		ScheduleTasks(node.Left, exprID)
		ScheduleTasks(node.Right, exprID)
		if node.Left != nil && node.Right != nil && node.Left.Computed && node.Right.Computed && !node.Computed {
			task := &models.Task{
				ID:            models.TaskIDCounter,
				ExpressionID:  exprID,
				Arg1:          node.Left.Value,
				Arg2:          node.Right.Value,
				Operation:     node.Operator,
				OperationTime: getOperationTime(node.Operator),
			}
			models.TaskIDCounter++
			models.TasksQueue = append(models.TasksQueue, task)
			node.TaskID = task.ID
		}
	}
}

func getOperationTime(op string) int {
	if GlobalConfig != nil {
		switch op {
		case "+":
			return GlobalConfig.TimeAdditionMs
		case "-":
			return GlobalConfig.TimeSubtractionMs
		case "*":
			return GlobalConfig.TimeMultiplicationsMs
		case "/":
			return GlobalConfig.TimeDivisionsMs
		default:
			return 1000
		}
	}

	var envVar string
	switch op {
	case "+":
		envVar = os.Getenv("TIME_ADDITION_MS")
		if v, err := strconv.Atoi(envVar); err == nil && v > 0 {
			return v
		}
		return 1000
	case "-":
		envVar = os.Getenv("TIME_SUBTRACTION_MS")
		if v, err := strconv.Atoi(envVar); err == nil && v > 0 {
			return v
		}
		return 1000
	case "*":
		envVar = os.Getenv("TIME_MULTIPLICATIONS_MS")
		if v, err := strconv.Atoi(envVar); err == nil && v > 0 {
			return v
		}
		return 2000
	case "/":
		envVar = os.Getenv("TIME_DIVISIONS_MS")
		if v, err := strconv.Atoi(envVar); err == nil && v > 0 {
			return v
		}
		return 2000
	default:
		return 1000
	}
}

func ScheduleReadyTasks(node *models.Node, exprID int64) {
	if node == nil {
		return
	}
	if node.Operator != "" && !node.Computed {
		if node.Left != nil && node.Right != nil && node.Left.Computed && node.Right.Computed {
			if node.TaskID == 0 {
				task := &models.Task{
					ID:            models.TaskIDCounter,
					ExpressionID:  exprID,
					Arg1:          node.Left.Value,
					Arg2:          node.Right.Value,
					Operation:     node.Operator,
					OperationTime: getOperationTime(node.Operator),
				}
				models.TaskIDCounter++
				models.TasksQueue = append(models.TasksQueue, task)
				node.TaskID = task.ID
			}
		}
	}
	ScheduleReadyTasks(node.Left, exprID)
	ScheduleReadyTasks(node.Right, exprID)
}

func UpdateASTWithTask(expr *models.Expression, taskID int64, result float64) bool {
	found := false
	var traverse func(node *models.Node)
	traverse = func(node *models.Node) {
		if node == nil || found {
			return
		}
		if node.TaskID == taskID && !node.Computed {
			node.Value = result
			node.Computed = true
			found = true
			node.TaskID = 0
			if node.Parent != nil {
				ScheduleReadyTasks(node.Parent, expr.ID)
			}
			return
		}
		traverse(node.Left)
		traverse(node.Right)
	}
	traverse(expr.Root)
	return found
}
