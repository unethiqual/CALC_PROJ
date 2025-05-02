package api

import (
    "encoding/json"
    "net/http"

    "golang.org/x/crypto/bcrypt"
    "github.com/unethiqual/CALC_PROJ/database"
    "github.com/unethiqual/CALC_PROJ/orchestrator/models"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }

    _, err = database.DB.Exec(
        "INSERT INTO users (login, password_hash) VALUES ($1, $2)",
        req.Login, string(passwordHash),
    )
    if err != nil {
        http.Error(w, "Failed to register user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var user models.User
    err := database.DB.Get(&user, "SELECT * FROM users WHERE login = $1", req.Login)
    if err != nil {
        http.Error(w, "Invalid login or password", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        http.Error(w, "Invalid login or password", http.StatusUnauthorized)
        return
    }

    token, err := generateJWT(user.ID)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Expression string `json:"expression"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusUnprocessableEntity)
        return
    }

    exprID, err := models.AddExpression(req.Expression)
    if err != nil {
        http.Error(w, "Failed to add expression", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int64{"id": exprID})
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
    expressions, err := models.GetExpressions()
    if err != nil {
        http.Error(w, "Failed to fetch expressions", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    expression, err := models.GetExpressionByID(id)
    if err != nil {
        http.Error(w, "Expression not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{"expression": expression})
}