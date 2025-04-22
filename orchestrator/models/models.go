package models

import "sync"

// ExpressionStatus Статусы вычисления
type ExpressionStatus string

const (
	StatusPending    ExpressionStatus = "pending"
	StatusInProgress ExpressionStatus = "in-progress"
	StatusCompleted  ExpressionStatus = "completed"
)

type Expression struct {
	ID     int64            `json:"id"`
	Expr   string           `json:"expression"`
	Status ExpressionStatus `json:"status"`
	Result float64          `json:"result"`
	Root   *Node            `json:"-"`
}

type Task struct {
	ID            int64   `json:"id"`
	ExpressionID  int64   `json:"-"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Node struct {
	Operator string  // например, "+", "-", "*" или "/"
	Value    float64 // значение, если вычислено
	Left     *Node
	Right    *Node
	Computed bool  // вычислено ли значение
	TaskID   int64 // ID задачи, соответствующей этому узлу (если не вычислено)
	Parent   *Node // для отслеживания зависимости
}

var (
	ExprIDCounter int64 = 1
	TaskIDCounter int64 = 1

	Expressions = make(map[int64]*Expression)
	TasksQueue  = make([]*Task, 0)
	Mu          sync.Mutex
)
