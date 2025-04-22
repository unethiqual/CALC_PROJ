package main

import (
	"github.com/unethiqual/CALC_PROJ/orchestrator/config"
	"github.com/unethiqual/CALC_PROJ/orchestrator/handlers"
	"github.com/unethiqual/CALC_PROJ/orchestrator/scheduler"
	"log"
	"net/http"
	"os"
)

func main() {
	if _, err := os.Stat("config.json"); err == nil {
		cfg, err := config.LoadConfig("config.json")
		if err != nil {
			log.Println("Ошибка загрузки конфигурации:", err)
		} else {
			scheduler.GlobalConfig = cfg
			log.Println("Конфигурация загружена из config.json")
		}
	} else {
		log.Println("Файл config.json не найден, используются переменные окружения")
	}

	
	http.HandleFunc("/", handlers.HandleFrontend)
	http.HandleFunc("/api/v1/calculate", handlers.HandleCalculate)
	http.HandleFunc("/api/v1/expressions", handlers.HandleGetExpressions)
	http.HandleFunc("/api/v1/expressions/", handlers.HandleGetExpression)
	http.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.HandleGetTask(w, r)
		case http.MethodPost:
			handlers.HandlePostTask(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Оркестратор запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
