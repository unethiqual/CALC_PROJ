package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/unethiqual/CALC_PROJ/config"
    "github.com/unethiqual/CALC_PROJ/database"
    "github.com/unethiqual/CALC_PROJ/orchestrator/api"
    "github.com/unethiqual/CALC_PROJ/orchestrator/grpc"
)

func main() {
    cfg := config.LoadConfig()

    // Initialize database
    database.InitDB(cfg.DatabaseURL)

    router := mux.NewRouter()

    // Public routes
    router.HandleFunc("/api/v1/register", api.RegisterHandler).Methods("POST")
    router.HandleFunc("/api/v1/login", api.LoginHandler).Methods("POST")

    // Protected routes
    protected := router.PathPrefix("/api/v1").Subrouter()
    protected.Use(api.JWTMiddleware)
    protected.HandleFunc("/calculate", api.AddExpressionHandler).Methods("POST")
    protected.HandleFunc("/expressions", api.GetExpressionsHandler).Methods("GET")
    protected.HandleFunc("/expressions/{id}", api.GetExpressionByIDHandler).Methods("GET")

    // Start gRPC server in a separate goroutine
    go grpc.StartGRPCServer()

    log.Println("Orchestrator is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}