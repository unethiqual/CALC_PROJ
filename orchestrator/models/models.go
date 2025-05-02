package models

import (
    "github.com/unethiqual/CALC_PROJ/database"
)

type User struct {
    ID           int    `db:"id"`
    Login        string `db:"login"`
    PasswordHash string `db:"password_hash"`
}

type Expression struct {
    ID         int64   `db:"id"`
    UserID     int     `db:"user_id"`
    Expression string  `db:"expression"`
    Status     string  `db:"status"`
    Result     *float64 `db:"result"`
}

type Task struct {
    ID            int64   `db:"id"`
    ExpressionID  int64   `db:"expression_id"`
    Arg1          float64 `db:"arg1"`
    Arg2          float64 `db:"arg2"`
    Operation     string  `db:"operation"`
    OperationTime int     `db:"operation_time"`
    Status        string  `db:"status"`
}

func AddExpression(userID int, expression string) (int64, error) {
    result, err := database.DB.Exec(
        "INSERT INTO expressions (user_id, expression) VALUES ($1, $2) RETURNING id",
        userID, expression,
    )
    if err != nil {
        return 0, err
    }

    id, _ := result.LastInsertId()
    return id, nil
}

func GetExpressions(userID int) ([]Expression, error) {
    var expressions []Expression
    err := database.DB.Select(&expressions, "SELECT * FROM expressions WHERE user_id = $1", userID)
    return expressions, err
}

func GetExpressionByID(userID int, id int64) (*Expression, error) {
    var expression Expression
    err := database.DB.Get(&expression, "SELECT * FROM expressions WHERE id = $1 AND user_id = $2", id, userID)
    if err != nil {
        return nil, err
    }
    return &expression, nil
}